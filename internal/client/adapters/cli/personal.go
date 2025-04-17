package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gokeeper/pkg/domain"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type PrivateService interface {
	Save(ctx context.Context, pd domain.Data, inputUser domain.InUserRequest, saveLocalOnError bool) error
	GetAll(ctx context.Context, gpr domain.GetAllRequest, inputUser domain.InUserRequest) ([]domain.Data, error)
	Get(ctx context.Context, id string, inputUser domain.InUserRequest) (*domain.Data, error)
	Delete(ctx context.Context, pd domain.DeleteRequest) error
	Upload(ctx context.Context) error
}

type PrivateCLI struct {
	privateService PrivateService
}

func NewPrivateCLI(privateService PrivateService) *PrivateCLI {
	return &PrivateCLI{
		privateService: privateService,
	}
}

func (pc *PrivateCLI) GetCommands() []*cobra.Command {
	return []*cobra.Command{
		pc.createSaveCommand(),
		pc.createGetCommand(),
		pc.createGetAllCommand(),
		pc.createDeleteCommand(),
		pc.createUploadCommand(),
	}
}

func (pc *PrivateCLI) createSaveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save data. Available data types: auth, file, card, text",
		Run:   pc.save,
	}

	addCommonAuthFlags(cmd)
	cmd.Flags().String("type", "", "Data type")
	cmd.Flags().String("id", "", "Data key")
	cmd.Flags().String("meta", "", "Meta information (not encrypted)")
	cmd.Flags().Bool("save-local-on-error", false, "Save locally if server unavailable")

	// Data-specific flags
	cmd.Flags().String("data-login", "", "Login for external resource")
	cmd.Flags().String("data-password", "", "Password for external resource")
	cmd.Flags().String("data-number", "", "Card number")
	cmd.Flags().String("data-name", "", "Card name")
	cmd.Flags().String("data-date", "", "Card expiration date")
	cmd.Flags().Uint16("data-secure", 0, "Card security code")
	cmd.Flags().String("text", "", "Text for saving")
	cmd.Flags().String("file", "", "File with data for saving")

	return cmd
}

func (pc *PrivateCLI) createGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get data",
		Run:   pc.get,
	}

	addCommonAuthFlags(cmd)
	cmd.Flags().String("id", "", "Data key")
	cmd.Flags().String("output", "", "Output file")

	return cmd
}

func (pc *PrivateCLI) createGetAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_all",
		Short: "Get all private data",
		Run:   pc.getAll,
	}

	addCommonAuthFlags(cmd)
	cmd.Flags().Uint64("limit", 10, "Number of elements")
	cmd.Flags().Uint64("offset", 0, "Page number")
	cmd.Flags().String("output", "", "Output file")

	return cmd
}

func (pc *PrivateCLI) createDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete private data",
		Run:   pc.delete,
	}

	cmd.Flags().String("id", "", "Data key")

	return cmd
}

func (pc *PrivateCLI) createUploadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload private data",
		Run:   pc.upload,
	}
}

func addCommonAuthFlags(cmd *cobra.Command) {
	cmd.Flags().String("login", "", "Authentication login")
	cmd.Flags().String("password", "", "Authentication password")
}

type dataHandler func(*cobra.Command) ([]byte, error)

var dataHandlers = map[domain.Type]dataHandler{
	domain.LOGIN_PASSWORD: func(cmd *cobra.Command) ([]byte, error) {
		return handleAuthData(cmd)
	},
	domain.CARD: func(cmd *cobra.Command) ([]byte, error) {
		return handleCardData(cmd)
	},
	domain.TEXT: func(cmd *cobra.Command) ([]byte, error) {
		return handleTextData(cmd)
	},
	domain.BYTES: func(cmd *cobra.Command) ([]byte, error) {
		return handleFileData(cmd, true)
	},
}

func (pc *PrivateCLI) save(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	u := pc.authenticate(cmd)

	id := getInputString(cmd, "id", "Enter id: ")
	dataTypeStr := getInputString(cmd, "type", "Enter data type: ")
	metaDataStr := getInputString(cmd, "meta", "Enter meta data: ")
	saveLocalOnError, _ := cmd.Flags().GetBool("save-local-on-error")

	dataType := parseType(dataTypeStr)
	if dataType == domain.UNKNOWN {
		pc.handleError(fmt.Errorf("invalid data type"))
		return
	}

	handler, exists := dataHandlers[dataType]
	if !exists {
		pc.handleError(fmt.Errorf("unsupported data type: %s", dataType))
		return
	}

	data, err := handler(cmd)
	if err != nil {
		pc.handleError(err)
		return
	}

	pd := domain.Data{
		ID:       id,
		DataType: dataType,
		MetaData: []byte(metaDataStr),
		Data:     data,
		SavedAt:  time.Now(),
	}

	if err := pc.privateService.Save(ctx, pd, *u, saveLocalOnError); err != nil {
		if errors.Is(err, domain.WarnServerUnavailable) {
			fmt.Println("Your data was saved locally, try command \"upload\" for uploading your data to the server")
			return
		}
		pc.handleError(err)
		return
	}

	fmt.Println("Your data was successfully saved")
}

func (pc *PrivateCLI) get(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	u := pc.authenticate(cmd)
	id := getInputString(cmd, "id", "Enter id: ")

	data, err := pc.privateService.Get(ctx, id, *u)
	if err != nil {
		pc.handleGetError(err, id)
		return
	}

	if err := pc.handleOutput(cmd, data); err != nil {
		pc.handleError(err)
	}
}

func (pc *PrivateCLI) getAll(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	u := pc.authenticate(cmd)

	limit, _ := cmd.Flags().GetUint64("limit")
	offset, _ := cmd.Flags().GetUint64("offset")

	data, err := pc.privateService.GetAll(ctx, domain.GetAllRequest{Limit: limit, Offset: offset}, *u)
	if err != nil {
		pc.handleError(err)
		return
	}

	resBytes, err := json.Marshal(data)
	if err != nil {
		pc.handleError(fmt.Errorf("marshaling error: %w", err))
		return
	}

	if err := pc.saveToOutput(cmd, resBytes); err != nil {
		pc.handleError(err)
	}
}

func (pc *PrivateCLI) delete(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()
	id := getInputString(cmd, "id", "Enter id: ")

	if err := pc.privateService.Delete(ctx, domain.DeleteRequest{ID: id, DeletedAt: time.Now()}); err != nil {
		pc.handleError(err)
		return
	}

	fmt.Printf("Your data with id %s was successfully deleted\n", id)
}

func (pc *PrivateCLI) upload(cmd *cobra.Command, _ []string) {
	if err := pc.privateService.Upload(cmd.Context()); err != nil {
		pc.handleError(err)
		return
	}
	fmt.Println("Your data was successfully uploaded")
}

func parseType(dataType string) domain.Type {
	switch dataType {
	case "auth":
		return domain.LOGIN_PASSWORD
	case "card":
		return domain.CARD
	case "text":
		return domain.TEXT
	case "file":
		return domain.BYTES
	default:
		return domain.UNKNOWN
	}
}

func getInputString(cmd *cobra.Command, flagName, prompt string) string {
	val, err := cmd.Flags().GetString(flagName)
	if err != nil || val == "" {
		fmt.Print(prompt)
		fmt.Scanf("%s", &val)
	}
	return val
}

func (pc *PrivateCLI) authenticate(cmd *cobra.Command) *domain.InUserRequest {
	return &domain.InUserRequest{
		Login:    getInputString(cmd, "login", "Enter your login: "),
		Password: getInputString(cmd, "password", "Enter your password: "),
	}
}

func (pc *PrivateCLI) handleError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func (pc *PrivateCLI) handleGetError(err error, id string) {
	if errors.Is(err, domain.ErrPrivateDataNotFound) {
		fmt.Printf("Data with id %s was not found\n", id)
		return
	}
	pc.handleError(err)
}

func (pc *PrivateCLI) handleOutput(cmd *cobra.Command, data *domain.Data) error {
	switch data.DataType {
	case domain.LOGIN_PASSWORD, domain.CARD:
		fmt.Printf("%s\n\n%s\n", data.MetaData, data.Data)
		return nil
	case domain.TEXT, domain.BYTES:
		return pc.saveToOutput(cmd, data.Data)
	default:
		return errors.New("unsupported data type")
	}
}

func (pc *PrivateCLI) saveToOutput(cmd *cobra.Command, data []byte) error {
	filePath, _ := cmd.Flags().GetString("output")
	if filePath == "" {
		fmt.Print("Enter output file: ")
		fmt.Scanf("%s", &filePath)
	}
	return pc.saveDataToFile(data, filePath)
}

func (pc *PrivateCLI) saveDataToFile(data []byte, filePath string) error {
	return os.WriteFile(filePath, data, 0666)
}

func handleAuthData(cmd *cobra.Command) ([]byte, error) {
	dataLogin := getInputString(cmd, "data-login", "Enter login for saving: ")
	dataPassword := getInputString(cmd, "data-password", "Enter password for saving: ")

	return json.Marshal(domain.LoginPasswordData{
		Login:    dataLogin,
		Password: dataPassword,
	})
}

func handleCardData(cmd *cobra.Command) ([]byte, error) {
	dataNumber := getInputString(cmd, "data-number", "Enter card number: ")
	dataName := getInputString(cmd, "data-name", "Enter card name: ")
	dataDate := getInputString(cmd, "data-date", "Enter expiration date: ")
	dataSecure := getInputUint16(cmd, "data-secure", "Enter CVV: ")

	return json.Marshal(domain.CardData{
		Number: dataNumber,
		Name:   dataName,
		Date:   dataDate,
		Secure: dataSecure,
	})
}

func handleTextData(cmd *cobra.Command) ([]byte, error) {
	if text, _ := cmd.Flags().GetString("text"); text != "" {
		return []byte(text), nil
	}
	return handleFileData(cmd, false)
}

func handleFileData(cmd *cobra.Command, required bool) ([]byte, error) {
	filePath, _ := cmd.Flags().GetString("file")
	if filePath == "" && required {
		fmt.Print("Enter file path: ")
		fmt.Scanf("%s", &filePath)
	}
	return os.ReadFile(filePath)
}

func getInputUint16(cmd *cobra.Command, flagName, prompt string) uint16 {
	val, err := cmd.Flags().GetUint16(flagName)
	if err != nil || val == 0 {
		fmt.Print(prompt)
		fmt.Scanf("%d", &val)
	}
	return val
}

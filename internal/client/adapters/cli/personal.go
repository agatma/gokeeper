package cli

import (
	"bytes"
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
	GetPage(ctx context.Context, gpr domain.GetAllRequest, inputUser domain.InUserRequest) ([]domain.Data, error)
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
	cmdSave := &cobra.Command{
		Use:   "save",
		Short: "save data. Available data types: auth, file, card, text",
		Run:   pc.save,
	}
	// for authentication and encryption
	cmdSave.Flags().String("login", "", "authentication on server")
	cmdSave.Flags().String("password", "", "authentication on server")

	// data type
	cmdSave.Flags().String("type", "", "data type")

	// data key
	cmdSave.Flags().String("id", "", "data key")

	// meta information
	cmdSave.Flags().String("meta", "", "meta information for saving. Will not encrypted")

	// save-local-on-error
	cmdSave.Flags().Bool("save-local-on-error", false, "data will be saved locally if server is unavailable. You should run upload command when you will have access to the server")

	// for login and password
	cmdSave.Flags().String("data-login", "", "login on external resourse")
	cmdSave.Flags().String("data-password", "", "password on external resource")

	// for card data
	cmdSave.Flags().String("data-number", "", "card number")
	cmdSave.Flags().String("data-name", "", "card name")
	cmdSave.Flags().String("data-date", "", "card expiration date")
	cmdSave.Flags().Uint16("data-secure", 0, "card security code")

	// for text and bytes
	cmdSave.Flags().String("text", "", "text for saving")
	cmdSave.Flags().String("file", "", "file with data for saving")

	cmdGet := &cobra.Command{
		Use:   "get",
		Short: "get data",
		Run:   pc.get,
	}
	// for authentication and encryption
	cmdGet.Flags().String("login", "", "authentication on server")
	cmdGet.Flags().String("password", "", "authentication on server")

	cmdGet.Flags().String("id", "", "data key")

	cmdGet.Flags().String("output", "", "file for output data")

	cmdGetPage := &cobra.Command{
		Use:   "getpage",
		Short: "getpage data",
		Run:   pc.getPage,
	}
	// for authentication and encryption
	cmdGetPage.Flags().String("login", "", "authentication on server")
	cmdGetPage.Flags().String("password", "", "authentication on server")

	cmdGetPage.Flags().Uint64("limit", 10, "num of elements")
	cmdGetPage.Flags().Uint64("offset", 0, "page number")

	cmdGetPage.Flags().String("output", "", "file for output data")

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: "delete data",
		Run:   pc.delete,
	}
	// for authentication and encryption

	cmdDelete.Flags().String("id", "", "data key")

	cmdUpload := &cobra.Command{
		Use:   "upload",
		Short: "upload data",
		Run:   pc.upload,
	}

	return []*cobra.Command{cmdSave, cmdGet, cmdGetPage, cmdDelete, cmdUpload}
}

func (pc *PrivateCLI) getCardData(cmd *cobra.Command) ([]byte, error) {
	dataNumber, err := cmd.Flags().GetString("data-number")
	if err != nil || dataNumber == "" {
		fmt.Print("Enter your card number (without spaces): ")
		fmt.Scanf("%s", &dataNumber)
	}
	dataName, err := cmd.Flags().GetString("data-name")
	if err != nil || dataName == "" {
		fmt.Print("Enter your name: ")
		fmt.Scanf("%s", &dataName)
	}
	dataDate, err := cmd.Flags().GetString("data-date")
	if err != nil || dataDate == "" {
		fmt.Print("Enter card expiration date: ")
		fmt.Scanf("%s", &dataDate)
	}
	dataSecure, err := cmd.Flags().GetUint16("data-secure")
	if err != nil || dataSecure == 0 {
		fmt.Print("Enter your cvv: ")
		fmt.Scanf("%d", &dataSecure)
	}

	cardData := domain.CardData{
		Number: dataNumber,
		Name:   dataName,
		Secure: dataSecure,
		Date:   dataDate,
	}
	data, err := json.Marshal(cardData)
	if err != nil {
		return nil, errors.New("internal error")
	}
	return data, nil
}

func (pc *PrivateCLI) getTextData(cmd *cobra.Command) []byte {
	dataText, err := cmd.Flags().GetString("text")
	if err != nil || dataText == "" {
		data, err := pc.getFileData(cmd, false)
		if err != nil {
			fmt.Print("Enter your text: ")
			fmt.Scanf("%s", &dataText)
			return []byte(dataText)
		}
		return data
	}
	return []byte(dataText)
}

func (pc *PrivateCLI) saveDataToFile(data []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func (pc *PrivateCLI) getFileData(cmd *cobra.Command, askFile bool) ([]byte, error) {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil || filePath == "" {
		if !askFile {
			return nil, err
		}
		fmt.Print("Enter file path: ")
		fmt.Scanf("%s", &filePath)
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (pc *PrivateCLI) getBytesData(cmd *cobra.Command) ([]byte, error) {
	return pc.getFileData(cmd, true)
}

func (pc *PrivateCLI) getAuthData(cmd *cobra.Command) ([]byte, error) {
	dataLogin, err := cmd.Flags().GetString("data-login")
	if err != nil || dataLogin == "" {
		fmt.Print("Enter your login for saving: ")
		fmt.Scanf("%s", &dataLogin)
	}
	dataPassword, err := cmd.Flags().GetString("data-password")
	if err != nil || dataPassword == "" {
		fmt.Print("Enter your password for saving: ")
		fmt.Scanf("%s", &dataPassword)
	}

	loginPasswordData := domain.LoginPasswordData{
		Login:    dataLogin,
		Password: dataPassword,
	}
	data, err := json.Marshal(loginPasswordData)
	if err != nil {
		return nil, errors.New("internal error")
	}
	return data, nil
}

func (pc *PrivateCLI) authenticate(cmd *cobra.Command) *domain.InUserRequest {
	login, err := cmd.Flags().GetString("login")
	if err != nil || login == "" {
		fmt.Print("Enter your login: ")
		fmt.Scanf("%s", &login)
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil || password == "" {
		fmt.Print("Enter your password: ")
		fmt.Scanf("%s", &password)
	}

	return &domain.InUserRequest{
		Login:    login,
		Password: password,
	}
}

func (pc *PrivateCLI) save(cmd *cobra.Command, _ []string) {
	saveLocalOnError, _ := cmd.Flags().GetBool("save-local-on-error")

	u := pc.authenticate(cmd)

	id, err := cmd.Flags().GetString("id")
	if err != nil || id == "" {
		fmt.Print("Enter id: ")
		fmt.Scanf("%s", &id)
	}

	dataTypeStr, err := cmd.Flags().GetString("type")
	if err != nil || dataTypeStr == "" {
		fmt.Print("Enter data type: ")
		fmt.Scanf("%s", &dataTypeStr)
	}

	metaDataStr, err := cmd.Flags().GetString("meta")
	if err != nil || metaDataStr == "" {
		fmt.Print("Enter meta data: ")
		fmt.Scanf("%s", &metaDataStr)
	}

	var dataType domain.Type
	var data []byte
	switch dataTypeStr {
	case "auth":
		dataType = domain.LOGIN_PASSWORD
		data, err = pc.getAuthData(cmd)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	case "card":
		dataType = domain.CARD
		data, err = pc.getCardData(cmd)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	case "text":
		dataType = domain.TEXT
		data = pc.getTextData(cmd)
	case "file":
		dataType = domain.BYTES
		data, err = pc.getBytesData(cmd)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	default:
		fmt.Printf("invalid data type provided")
		return
	}

	pd := domain.Data{
		ID:       id,
		DataType: dataType,
		MetaData: []byte(metaDataStr),
		Data:     data,
		SavedAt:  time.Now(),
	}

	err = pc.privateService.Save(cmd.Context(), pd, *u, saveLocalOnError)
	if err != nil {
		if errors.Is(err, domain.WarnServerUnavailable) {
			fmt.Printf("Your data was saved locally, try command \"upload\" for uploading your data to the server")
			return
		}
		fmt.Printf("Exception occured: %v", err)
		return
	}

	fmt.Println("Your data was successfully saved")
}

func (pc *PrivateCLI) get(cmd *cobra.Command, _ []string) {
	u := pc.authenticate(cmd)

	id, err := cmd.Flags().GetString("id")
	if err != nil || id == "" {
		fmt.Print("Enter id: ")
		fmt.Scanf("%s", &id)
	}

	data, err := pc.privateService.Get(cmd.Context(), id, *u)
	if err != nil {
		if errors.Is(err, domain.ErrPrivateDataNotFound) {
			fmt.Printf("Data with id %s was not found", id)
			return
		}
		fmt.Printf("Exception occured: %v", err)
		return
	}

	switch data.DataType {
	case domain.LOGIN_PASSWORD:
		fmt.Printf("%s\n\n", data.MetaData)
		fmt.Printf("%s\n", string(data.Data))
	case domain.CARD:
		fmt.Printf("%s\n\n", data.MetaData)
		fmt.Printf("%s\n", string(data.Data))
	case domain.TEXT:
		filePath, err := cmd.Flags().GetString("output")
		if err != nil || filePath == "" {
			fmt.Printf("%s\n\n", data.MetaData)
			fmt.Printf("%s\n", string(data.Data))
		} else {
			pc.saveDataToFile(data.Data, filePath)
		}
	case domain.BYTES:
		filePath, err := cmd.Flags().GetString("output")
		if err != nil || filePath == "" {
			fmt.Print("Enter output file: ")
			fmt.Scanf("%s", &filePath)
		}
		pc.saveDataToFile(data.Data, filePath)
	default:
		fmt.Printf("invalid data type :(")
		return
	}
}

func (pc *PrivateCLI) getPage(cmd *cobra.Command, _ []string) {
	u := pc.authenticate(cmd)

	limit, err := cmd.Flags().GetUint64("limit")
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	offset, err := cmd.Flags().GetUint64("offset")
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	getPageRequest := domain.GetAllRequest{
		Limit:  limit,
		Offset: offset,
	}

	data, err := pc.privateService.GetPage(cmd.Context(), getPageRequest, *u)
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	filePath, err := cmd.Flags().GetString("output")
	if err != nil || filePath == "" {
		fmt.Print("Enter output file: ")
		fmt.Scanf("%s", &filePath)
	}

	resBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	err = pc.saveDataToFile(resBytes, filePath)
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}
}

func (pc *PrivateCLI) delete(cmd *cobra.Command, _ []string) {
	id, err := cmd.Flags().GetString("id")
	if err != nil || id == "" {
		fmt.Print("Enter id: ")
		fmt.Scanf("%s", &id)
	}

	deleteRequest := domain.DeleteRequest{
		ID:        id,
		DeletedAt: time.Now(),
	}

	err = pc.privateService.Delete(cmd.Context(), deleteRequest)
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	fmt.Printf("Your data with id %s was successfully deleted\n", id)
}

func (pc *PrivateCLI) upload(cmd *cobra.Command, _ []string) {
	err := pc.privateService.Upload(cmd.Context())
	if err != nil {
		fmt.Printf("Exception occured: %v", err)
		return
	}

	fmt.Printf("Your data was successfully uploaded")
}

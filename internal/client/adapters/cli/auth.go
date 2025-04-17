package cli

import (
	"context"
	"errors"
	"fmt"
	"gokeeper/pkg/domain"
	"log"

	"github.com/spf13/cobra"
)

type AuthService interface {
	Register(ctx context.Context, user domain.InUserRequest, saveJWT bool) error
	Login(ctx context.Context, user domain.InUserRequest, saveJWT bool) (string, error)
	GetJwt(ctx context.Context) (string, error)
}

type AuthCLI struct {
	authService AuthService
}

func NewAuthCLI(authService AuthService) *AuthCLI {
	return &AuthCLI{
		authService: authService,
	}
}

func (ac *AuthCLI) GetCommands() []*cobra.Command {
	cmdLogin := &cobra.Command{
		Use:   "login",
		Short: "Command for login user",
		Run:   ac.login,
	}
	cmdLogin.Flags().String("login", "", "authentication on server")
	cmdLogin.Flags().String("password", "", "authentication on server")

	cmdRegister := &cobra.Command{
		Use:   "register",
		Short: "Command for registration",
		Run:   ac.register,
	}
	cmdRegister.Flags().String("login", "", "register on server")
	cmdRegister.Flags().String("password", "", "register on server")
	return []*cobra.Command{cmdLogin, cmdRegister}
}

func (ac *AuthCLI) login(cmd *cobra.Command, _ []string) {
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

	u := domain.InUserRequest{
		Login:    login,
		Password: password,
	}

	_, err = ac.authService.Login(cmd.Context(), u, true)
	if err != nil {
		if errors.Is(err, domain.ErrUserAuthentication) {
			log.Printf("authentication failed")
			return
		}
		log.Printf("server unavailable, try later")
	}

	fmt.Print("Successfully logged in\n")
}

func (ac *AuthCLI) register(cmd *cobra.Command, _ []string) {
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

	u := domain.InUserRequest{
		Login:    login,
		Password: password,
	}

	err = ac.authService.Register(cmd.Context(), u, true)
	if err != nil {
		if errors.Is(err, domain.ErrUserConflict) {
			log.Printf("user with same login already exists")
			return
		}
		log.Printf("server unavailable, try later")
	}

	fmt.Print("Successfully registered\n")
}

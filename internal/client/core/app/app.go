package app

import (
	"context"
	"fmt"
	"gokeeper/internal/client/adapters/cli"
	clients "gokeeper/internal/client/adapters/handlers"
	"gokeeper/internal/client/core/config"
	"gokeeper/internal/client/core/service"
	"gokeeper/internal/client/core/service/workers"
	"gokeeper/pkg/encrypter"

	"github.com/spf13/cobra"
)

type Client struct {
	CLI *cli.CLI
}

func NewClient(cfg *config.Config) *Client {
	c := clients.NewClients(cfg)
	w := workers.NewWorkers(cfg, c.PrivateClient)
	services := service.NewServices(
		w.FileWorker.JWTWorker,
		c.AuthClient,
		c.PrivateClient,
		encrypter.NewEncrypter(),
		w.FileWorker.PrivateFileWorker,
		w.Sender,
	)
	return &Client{
		CLI: cli.NewCLI(services.PrivateService, services.AuthService),
	}
}

func (a *Client) Run(ctx context.Context) error {
	var rootCmd = &cobra.Command{
		Use:   "gophkeeper",
		Short: "Client for private data",
		Long:  "Client for private data",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to gophkeeper")
		},
	}

	for _, cmd := range a.CLI.AuthCLI.GetCommands() {
		rootCmd.AddCommand(cmd)
	}
	for _, cmd := range a.CLI.PrivateCLI.GetCommands() {
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}

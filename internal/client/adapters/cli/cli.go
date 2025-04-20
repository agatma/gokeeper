package cli

type CLI struct {
	PrivateCLI *PrivateCLI
	AuthCLI    *AuthCLI
}

func NewCLI(privateService PrivateService, authService AuthService) *CLI {
	return &CLI{
		PrivateCLI: NewPrivateCLI(privateService),
		AuthCLI:    NewAuthCLI(authService),
	}
}

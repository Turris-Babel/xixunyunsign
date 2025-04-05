// cmd/login.go
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"xixunyunsign/internal/wire" // Import wire package
	// "xixunyunsign/service" // No longer needed here
	// "xixunyunsign/service/impl" // No longer needed here
	// "xixunyunsign/utils"      // No longer needed here
)

// LoginCmdRunner and NewLoginCmdRunner are now defined in internal/wire/wire.go

// RunLogin executes the login logic using the injected AuthService.
// This function will be called by the runner created by wire.
// The runner type is now wire.LoginCmdRunner.
func RunLogin(runner *wire.LoginCmdRunner, account, password, schoolID string) {
	token, err := runner.AuthService.Login(account, password, schoolID)
	if err != nil {
		log.Printf("登录失败: %v\n", err) // Log error
		return
	}
	log.Printf("登录成功！Token: %s\n", token) // Log success
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录到系统",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the runner using the wire injector
		runner, err := wire.InitializeLoginCmdRunner() // Call the injector
		if err != nil {
			log.Fatalf("无法初始化 Login 命令: %v", err)
		}

		// Run the command logic using the initialized runner
		// Call the RunLogin function defined in this package, passing the runner
		RunLogin(runner, u.account, u.password, u.schoolID)
	},
}

func init() {
	LoginCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	LoginCmd.Flags().StringVarP(&u.password, "password", "p", "", "密码")
	LoginCmd.Flags().StringVarP(&u.schoolID, "school_id", "i", "7", "学校id")
	LoginCmd.MarkFlagRequired("account")
	LoginCmd.MarkFlagRequired("password")
}

// Removed old Login function as logic is now in AuthService and LoginCmdRunner.Run
// Removed old getStringFromResult function

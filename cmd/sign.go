package cmd

import (
	"log"
	// "fmt" // Not needed if PushMsgToWechat is elsewhere
	// "io/ioutil" // Not needed if PushMsgToWechat is elsewhere
	// "net/http" // Not needed if PushMsgToWechat is elsewhere
	// "net/url" // Not needed if PushMsgToWechat is elsewhere

	"github.com/spf13/cobra"
	"xixunyunsign/internal/wire" // Import wire
	// "xixunyunsign/service" // No longer needed here
	// "xixunyunsign/utils" // No longer needed here
)

// signconfig is defined in cmd/struct.go.
// Define the Config variable of type signconfig (from struct.go)
var Config signconfig

// Assume secret_key is defined elsewhere in cmd package (e.g., topic.go or a shared var file)
// Assume PushMsgToWechat is defined elsewhere in cmd package (e.g., topic.go)

// SignCmdRunner holds dependencies for the sign command.
// Will be defined in internal/wire/wire.go

// RunSignIn executes the sign-in logic using the injected SignService.
func RunSignIn(runner *wire.SignCmdRunner, account, address, addressName, latitude, longitude string) {
	// Note: The service layer now handles getting user info, encryption, etc.
	message, err := runner.SignService.SignIn(account, address, addressName, latitude, longitude)
	if err != nil {
		log.Printf("签到失败: %v\n", err)
		// Handle Wechat push on error if needed
		// Use secret_key defined elsewhere
		if secret_key != "" {
			// Call PushMsgToWechat defined elsewhere
			PushMsgToWechat("签到失败", err.Error(), "9", secret_key)
		}
		return
	}
	log.Println(message) // Service returns "签到成功" on success
	// Handle Wechat push on success if needed
	// Use secret_key defined elsewhere
	if secret_key != "" {
		// Maybe get coordinates again if needed for the message?
		// Requires importing utils and handling error
		// lat, lon, _ := utils.GetCoordinates(account)
		// Call PushMsgToWechat defined elsewhere
		PushMsgToWechat("签到成功", message /* + formatted coordinates */, "9", secret_key)
	}
}

// SignCmd 定义签到命令
var SignCmd = &cobra.Command{
	Use:   "sign",
	Short: "执行签到",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize runner via wire (injector to be added later)
		runner, err := wire.InitializeSignCmdRunner() // Placeholder name
		if err != nil {
			log.Fatalf("无法初始化 Sign 命令: %v", err)
		}
		// Call the local function, passing the runner and flags
		RunSignIn(runner, u.account, Config.address, Config.address_name, Config.latitude, Config.longitude)
	},
}

func init() {
	SignCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号") // Keep one definition
	// SignCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号") // Remove duplicate
	SignCmd.Flags().StringVarP(&Config.address, "address", "", "", "地址(具体名称_小字部分)") // Use Config from cmd/struct.go
	SignCmd.Flags().StringVarP(&Config.address_name, "address_name", "", "", "地址名称")
	SignCmd.Flags().StringVarP(&Config.latitude, "latitude", "", "", "纬度")
	SignCmd.Flags().StringVarP(&Config.longitude, "longitude", "", "", "经度")
	SignCmd.Flags().StringVarP(&Config.remark, "remark", "", "0", "备注")
	SignCmd.Flags().StringVarP(&Config.comment, "comment", "", "", "评论")
	SignCmd.Flags().StringVarP(&Config.province, "province", "p", "", "省份")
	SignCmd.Flags().StringVarP(&Config.city, "city", "c", "", "城市")
	SignCmd.Flags().BoolVarP(&Config.debug, "debug", "d", false, "启用调试模式")
	SignCmd.Flags().StringVarP(&secret_key, "secret_key", "k", "", "server酱密钥") // Use secret_key from elsewhere

	// 标记必需的标志
	SignCmd.MarkFlagRequired("account")
	SignCmd.MarkFlagRequired("address")
}

// Removed old SignIn, rsaEncrypt, extractProvinceAndCity functions
// Logic is now in SignService

// Assume PushMsgToWechat is defined elsewhere (e.g., cmd/topic.go)
/*
func PushMsgToWechat(title, desp, channel, secret string) {
	...
}
*/

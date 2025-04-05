// cmd/query.go
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"xixunyunsign/internal/wire" // Import wire
	// "xixunyunsign/service" // No longer needed here
)

// QueryCmdRunner and NewQueryCmdRunner are now defined in internal/wire/wire.go

// RunQuerySignInfo executes the query sign info logic.
// Accepts the runner type defined in the wire package.
func RunQuerySignInfo(runner *wire.QueryCmdRunner, account string) {
	// Use the passed 'runner' variable
	data, err := runner.QueryService.QuerySignInfo(account)
	if err != nil {
		log.Printf("查询签到信息失败: %v\n", err)
		return
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

// RunSearchSchool is now defined in internal/wire/wire.go as a method on QueryCmdRunner

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询签到信息",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize runner via wire (injector to be added later)
		runner, err := wire.InitializeQueryCmdRunner() // Call the injector
		if err != nil {
			log.Fatalf("无法初始化 Query 命令: %v", err)
		}
		// Call the local function, passing the runner
		RunQuerySignInfo(runner, u.account)
	},
}

func init() {
	QueryCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	QueryCmd.MarkFlagRequired("account")
}

// Removed old Query function

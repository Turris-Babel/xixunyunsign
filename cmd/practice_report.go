package cmd

import (
	"fmt" // Keep for printing generated content
	"log"
	"xixunyunsign/internal/wire" // Import wire

	"github.com/spf13/cobra"
	// Imports no longer needed here: bytes, encoding/json, io, mime/multipart, net/http, net/url, os, time, utils
)

var (
	filePath     string
	role         string
	month        int8
	businessType string
	startDate    string
	endDate      string
	// attachment is now handled within the runner methods
	apiKey string
)

// var u UserInfo // u is defined in cmd/struct.go

// RunPracticeReport is called by the cobra command after initializing the runner.
// It orchestrates the report generation process using the runner's methods.
func RunPracticeReport(runner *wire.PracticeReportCmdRunner, account, filePath, role, apiKey, businessType, startDate, endDate string, month int8) {
	// Note: The runner methods now take 'account' as an argument where needed.
	attachment := runner.UploadImages(account, filePath) // Pass account
	if attachment == "" {
		log.Println("上传图片失败，中止操作。")
		return
	}

	// Pass month to GenerateContent
	content, err := runner.GenerateContent(role, apiKey, month)
	if err != nil {
		log.Printf("生成内容失败: %v\n", err)
		return
	}
	fmt.Println("生成的内容:") // Keep debug output for now
	fmt.Println(content)

	// Pass account to ReportsMonth
	runner.ReportsMonth(account, businessType, startDate, endDate, content, attachment)
}

var ExperimentalCmd = &cobra.Command{
	Use:   "experimental",
	Short: "实验性命令(自动月报)",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize runner via wire
		runner, err := wire.InitializePracticeReportCmdRunner()
		if err != nil {
			log.Fatalf("无法初始化 PracticeReport 命令: %v", err)
		}
		// Call the local function, passing the runner and flags
		// Pass u.account and month flag value
		RunPracticeReport(runner, u.account, filePath, role, apiKey, businessType, startDate, endDate, month)
	},
}

func init() {
	ExperimentalCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	ExperimentalCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "文件地址")
	ExperimentalCmd.Flags().StringVarP(&role, "role", "r", "", "工作角色")
	ExperimentalCmd.Flags().Int8VarP(&month, "month", "M", 1, "第几月（默认为1）")
	ExperimentalCmd.Flags().StringVarP(&businessType, "businessType", "b", "month", "报告类型(默认month)")
	ExperimentalCmd.Flags().StringVarP(&startDate, "startDate", "s", "", "开始日期(格式为20xx/xx/xx)")
	ExperimentalCmd.Flags().StringVarP(&endDate, "endDate", "e", "", "结束日期(格式为20xx/xx/xx)")
	ExperimentalCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "apikey(gemini-1.5-flash:generateContent)")
	ExperimentalCmd.MarkFlagRequired("filePath")
	ExperimentalCmd.MarkFlagRequired("role")
	ExperimentalCmd.MarkFlagRequired("account")
	ExperimentalCmd.MarkFlagRequired("startDate")
	ExperimentalCmd.MarkFlagRequired("endDate") // Corrected from endData
	ExperimentalCmd.MarkFlagRequired("apiKey")
}

// Removed old functions: MonthReportUploadSelectFile, UploadImages, GenerateContent, ReportsMonth
// Their logic is now part of PracticeReportCmdRunner in internal/wire/wire.go

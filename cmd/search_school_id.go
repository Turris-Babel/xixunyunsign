package cmd

import (
	"log"
	"xixunyunsign/internal/wire" // Import wire
	"xixunyunsign/utils"         // Re-add utils import

	"github.com/spf13/cobra"
)

var (
	schoolName string
)

var SchoolSearchIDCmd = &cobra.Command{
	Use:   "search", // Keep original Use value
	Short: "通过学校名称查询学校ID",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Consider moving this check/fetch logic to a service or initialization step
		// 检查数据库中是否有学校数据
		isEmpty, err := utils.IsSchoolInfoTableEmpty()
		if err != nil {
			log.Printf("检查数据库时发生错误: %v", err)
			return
		}

		// 如果学校数据表为空，则获取并保存学校数据
		if isEmpty {
			log.Println("学校信息表为空，正在从 API 获取...")
			err := utils.FetchAndSaveSchoolData()
			if err != nil {
				log.Printf("获取并保存学校数据失败: %v", err)
				// Continue anyway, maybe the search still works or gives a better error
			} else {
				log.Println("学校信息获取并保存成功。")
			}
		}

		// Initialize the QueryCmdRunner via wire
		runner, err := wire.InitializeQueryCmdRunner()
		if err != nil {
			log.Fatalf("无法初始化查询命令 runner: %v", err)
		}

		// Call the runner's search method
		// Note: The runner's RunSearchSchool method already handles printing
		runner.RunSearchSchool(schoolName)
	},
}

func init() {
	// 定义参数
	SchoolSearchIDCmd.Flags().StringVarP(&schoolName, "school_name", "s", "", "学校名称")
	SchoolSearchIDCmd.MarkFlagRequired("school_name")
}

// Removed old SearchSchoolID function as logic is now in QueryService and called via QueryCmdRunner

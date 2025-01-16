package cmd

import (
	"fmt"
	"log"
	"xixunyunsign/utils"

	"github.com/spf13/cobra"
)

var (
	schoolName string
)

var SchoolSearchIDCmd = &cobra.Command{
	Use:   "search",
	Short: "通过学校名称查询学校ID",
	Run: func(cmd *cobra.Command, args []string) {
		// 检查数据库中是否有学校数据
		isEmpty, err := utils.IsSchoolInfoTableEmpty()
		if err != nil {
			log.Printf("检查数据库时发生错误: %v", err)
			return
		}

		// 如果学校数据表为空，则获取并保存学校数据
		if isEmpty {
			err := utils.FetchAndSaveSchoolData()
			if err != nil {
				log.Printf("Error fetching and saving school data: %v", err)
				return
			}
		}

		// 调用查询函数
		_, err = SearchSchoolID(schoolName)
		if err != nil {
			return
		}
	},
}

func init() {
	// 定义参数
	SchoolSearchIDCmd.Flags().StringVarP(&schoolName, "school_name", "s", "", "学校名称")
	SchoolSearchIDCmd.MarkFlagRequired("school_name")
}

func SearchSchoolID(schoolName string) ([]utils.SchoolInfo, error) {
	// 调用utils中的查询函数，模糊匹配学校名称
	schools, err := utils.SearchSchoolID(schoolName)
	if err != nil {
		fmt.Println("查询时发生错误:", err)
		return nil, err
	}

	if len(schools) == 0 {
		fmt.Println("没有找到匹配的学校")
		return nil, nil
	}

	var schoolInfos []utils.SchoolInfo
	// 收集所有匹配学校的详细信息
	for _, school := range schools {
		fmt.Printf("学校名称: %s 对应的学校ID是: %s\n", school.SchoolName, school.SchoolID)
		schoolInfos = append(schoolInfos, utils.SchoolInfo{
			SchoolID:   school.SchoolID,
			SchoolName: school.SchoolName,
		})
	}

	return schoolInfos, nil
}

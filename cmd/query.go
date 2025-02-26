// cmd/query.go
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"xixunyunsign/utils"
)

func init() {
	QueryCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	err := QueryCmd.MarkFlagRequired("account")
	if err != nil {
		return
	}
}

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询签到信息",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := Query(u.account)
		if err != nil {
			return
		}
	},
}

// Query retrieves sign-in information for the given account.
func Query(account string) (map[string]interface{}, error) {
	token, _, _, err := utils.GetUser(account)
	if err != nil || token == "" {
		return nil, fmt.Errorf("未找到账号 %s 的 token", account)
	}

	userData, err := utils.GetAdditionalUserData(account)
	if err != nil {
		return nil, fmt.Errorf("获取用户额外信息失败: %v", err)
	}

	apiURL := "https://api.xixunyun.com/signin40/homepage"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	query := req.URL.Query()
	query.Add("month_date", "2024-12")
	query.Add("token", token)
	query.Add("from", "app")
	query.Add("version", "5.1.3")
	query.Add("platform", "android")
	query.Add("entrance_year", userData["entrance_year"])
	query.Add("graduate_year", userData["graduation_year"])
	query.Add("school_id", userData["school_id"])
	req.URL.RawQuery = query.Encode()

	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	code, ok := result["code"].(float64)
	if !ok || code != 20000 {
		message, _ := result["message"].(string)
		return nil, fmt.Errorf("查询失败: %s", message)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("解析数据失败：无效的响应结构")
	}

	signResourcesInfo, ok := data["sign_resources_info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("解析签到资源信息失败：无效的响应结构")
	}

	midLatitude := fmt.Sprintf("%v", signResourcesInfo["mid_sign_latitude"])
	midLongitude := fmt.Sprintf("%v", signResourcesInfo["mid_sign_longitude"])

	// 更新数据库中的经纬度信息
	err = utils.UpdateCoordinates(account, midLatitude, midLongitude)
	if err != nil {
		return nil, fmt.Errorf("保存经纬度信息失败: %v", err)
	}

	return data, nil
}

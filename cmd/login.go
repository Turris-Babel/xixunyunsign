// cmd/login.go
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"xixunyunsign/utils"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录到系统",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := Login(u.account, u.password, u.schoolID)
		if err != nil {
			return
		}
	},
}

func init() {
	LoginCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	LoginCmd.Flags().StringVarP(&u.password, "password", "p", "", "密码")
	LoginCmd.Flags().StringVarP(&u.schoolID, "school_id", "i", "7", "学校id")
	LoginCmd.MarkFlagRequired("account")
	LoginCmd.MarkFlagRequired("password")
}

// Login performs the login operation and returns the token or an error.
func Login(account, password, schoolID string) (string, error) {
	apiURL := "https://api.xixunyun.com/login/api"

	data := url.Values{}
	data.Set("app_version", "5.1.3")
	data.Set("registration_id", "")
	data.Set("uuid", "fd9dc13a49cc850c")
	data.Set("request_source", "3")
	data.Set("platform", "2")
	data.Set("mac", "7C:F3:1B:BB:F1:C4")
	data.Set("password", password)
	data.Set("system", "10")
	data.Set("school_id", schoolID)
	data.Set("model", "LM-G820")
	data.Set("app_id", "cn.vanber.xixunyun.saas")
	data.Set("account", account)
	data.Set("key", "")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	query := req.URL.Query()
	query.Add("from", "app")
	query.Add("version", "5.1.3")
	query.Add("platform", "android")
	query.Add("entrance_year", "0")
	query.Add("graduate_year", "0")
	query.Add("school_id", schoolID)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if code, ok := result["code"].(float64); !ok || code != 20000 {
		message, _ := result["message"].(string)
		return "", fmt.Errorf("登录失败: %s", message)
	}

	dataMap := result["data"].(map[string]interface{})
	token, ok := dataMap["token"].(string)
	if !ok {
		return "", fmt.Errorf("获取 token 失败")
	}

	// 保存到数据库
	err = utils.SaveUser(
		account,
		password,
		token,
		"", "", // 经纬度信息留空
		getStringFromResult(dataMap, "bind_phone"),
		getStringFromResult(dataMap, "user_number"),
		getStringFromResult(dataMap, "user_name"),
		dataMap["school_id"].(float64),
		getStringFromResult(dataMap, "sex"),
		getStringFromResult(dataMap, "class_name"),
		getStringFromResult(dataMap, "entrance_year"),
		getStringFromResult(dataMap, "graduation_year"),
	)
	if err != nil {
		return "", fmt.Errorf("保存用户信息失败: %v", err)
	}

	return token, nil
}

func getStringFromResult(dataMap map[string]interface{}, key string) string {
	if value, ok := dataMap[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return "" // 如果字段不存在或类型不匹配，返回空字符串
}

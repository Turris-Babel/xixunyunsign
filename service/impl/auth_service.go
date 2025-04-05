// service/impl/auth_service.go
package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"xixunyunsign/service"
)

// AuthServiceImpl implements the AuthService interface.
type AuthServiceImpl struct {
	userRepo   service.UserRepository
	httpClient *http.Client // Keep http client internal for now
}

// NewAuthService creates a new AuthService implementation.
func NewAuthService(repo service.UserRepository) service.AuthService {
	// Return the exported type
	return &AuthServiceImpl{
		userRepo:   repo,
		httpClient: &http.Client{}, // Initialize internal http client
	}
}

// Login performs the user login operation.
func (a *AuthServiceImpl) Login(account, password, schoolID string) (string, error) {
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
	q := req.URL.Query()
	q.Add("from", "app")
	q.Add("version", "5.1.3")
	q.Add("platform", "android")
	q.Add("entrance_year", "0")
	q.Add("graduate_year", "0")
	q.Add("school_id", schoolID)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
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
		return "", errors.New("登录失败: " + message)
	}

	dataMap, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("响应数据格式错误")
	}
	token, ok := dataMap["token"].(string)
	if !ok {
		return "", errors.New("获取 token 失败")
	}

	// 保存用户数据到数据库
	err = a.userRepo.SaveUser(
		account, password, token,
		"", "", // 经纬度为空
		getString(dataMap, "bind_phone"),
		getString(dataMap, "user_number"),
		getString(dataMap, "user_name"),
		dataMap["school_id"].(float64),
		getString(dataMap, "sex"),
		getString(dataMap, "class_name"),
		getString(dataMap, "entrance_year"),
		getString(dataMap, "graduation_year"),
	)
	if err != nil {
		return "", fmt.Errorf("保存用户信息失败: %v", err)
	}

	return token, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

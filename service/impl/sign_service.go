// service/impl/sign_service.go
package impl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"xixunyunsign/service"
	// "xixunyunsign/utils" // No longer needed if extractProvinceAndCity is local or moved
)

// SignServiceImpl implements the SignService interface.
type SignServiceImpl struct {
	userRepo service.UserRepository // Add UserRepository dependency
	// httpClient *http.Client      // Consider injecting http client
}

// NewSignService creates a new SignService implementation.
func NewSignService(repo service.UserRepository) service.SignService {
	// Return the exported type
	return &SignServiceImpl{
		userRepo: repo,
		// httpClient: &http.Client{},
	}
}

// SignIn performs the sign-in operation for a user.
func (s *SignServiceImpl) SignIn(account, address, addressName, latitude, longitude string) (string, error) {
	// Use injected userRepo
	token, dbLat, dbLong, err := s.userRepo.GetUser(account)
	if err != nil || token == "" {
		return "", fmt.Errorf("获取用户 token 失败或未找到账号 %s: %w", account, err)
	}
	if latitude == "" {
		latitude = dbLat
	}
	if longitude == "" {
		longitude = dbLong
	}
	if latitude == "" || longitude == "" {
		return "", errors.New("未提供经纬度信息，且数据库中不存在")
	}

	encryptedLat, err := rsaEncrypt([]byte(latitude))
	if err != nil {
		return "", fmt.Errorf("加密纬度失败: %v", err)
	}
	encryptedLong, err := rsaEncrypt([]byte(longitude))
	if err != nil {
		return "", fmt.Errorf("加密经度失败: %v", err)
	}

	// 从 address 提取省份和城市
	province, city, err := extractProvinceAndCity(address)
	if err != nil {
		return "", fmt.Errorf("地址格式不正确，无法提取省份和城市: %v", err)
	}

	apiURL := "https://api.xixunyun.com/signin_rsa"
	data := url.Values{}
	data.Set("address", address)
	data.Set("province", province)
	data.Set("city", city)
	data.Set("latitude", encryptedLat)
	data.Set("longitude", encryptedLong)
	data.Set("remark", "0")
	data.Set("comment", "")
	data.Set("address_name", addressName)
	data.Set("change_sign_resource", "0")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	q := req.URL.Query()
	q.Add("token", token)
	q.Add("from", "app")
	q.Add("version", "5.1.3")
	q.Add("platform", "android")
	q.Add("entrance_year", "0")
	q.Add("graduate_year", "0")
	q.Add("school_id", account) // 假设 schoolID 存在于 account 字段（或可修改为其他来源）
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应数据失败: %v", err)
	}
	if code, ok := result["code"].(float64); !ok || code != 20000 {
		return "", fmt.Errorf("签到失败: %v", result["message"])
	}
	return "签到成功", nil
}

// rsaEncrypt 使用公钥加密数据
func rsaEncrypt(origData []byte) (string, error) {
	publicKey := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDlYsiV3DsG+t8OFMLyhdmG2P2J
4GJwmwb1rKKcDZmTxEphPiYTeFIg4IFEiqDCATAPHs8UHypphZTK6LlzANyTzl9L
jQS6BYVQk81LhQ29dxyrXgwkRw9RdWaMPtcXRD4h6ovx6FQjwQlBM5vaHaJOHhEo
rHOSyd/deTvcS+hRSQIDAQAB
-----END PUBLIC KEY-----`
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("公钥解码失败")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("解析公钥类型失败")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func extractProvinceAndCity(address string) (string, string, error) {
	re := regexp.MustCompile(`(?P<province>[^省]+省)?(?P<city>[^市]+市)?`)
	matches := re.FindStringSubmatch(address)
	if len(matches) >= 3 {
		return matches[1], matches[2], nil
	}
	return "", "", fmt.Errorf("无法提取省份和城市")
}

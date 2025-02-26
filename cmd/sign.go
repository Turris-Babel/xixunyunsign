package cmd

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
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"xixunyunsign/utils"
)

var Config signconfig

// SignCmd 定义签到命令
var SignCmd = &cobra.Command{
	Use:   "sign",
	Short: "执行签到",
	Run: func(cmd *cobra.Command, args []string) {
		SignIn(u.account, Config.address, Config.address_name, Config.latitude, Config.longitude)
	},
}

func init() {
	SignCmd.Flags().StringVarP(&u.account, "account", "a", "", "账号")
	SignCmd.Flags().StringVarP(&Config.address, "address", "", "", "地址(具体名称_小字部分)")
	SignCmd.Flags().StringVarP(&Config.address_name, "address_name", "", "", "地址名称")
	SignCmd.Flags().StringVarP(&Config.latitude, "latitude", "", "", "纬度")
	SignCmd.Flags().StringVarP(&Config.longitude, "longitude", "", "", "经度")
	SignCmd.Flags().StringVarP(&Config.remark, "remark", "", "0", "备注")
	SignCmd.Flags().StringVarP(&Config.comment, "comment", "", "", "评论")
	SignCmd.Flags().StringVarP(&Config.province, "province", "p", "", "省份")
	SignCmd.Flags().StringVarP(&Config.city, "city", "c", "", "城市")
	SignCmd.Flags().BoolVarP(&Config.debug, "debug", "d", false, "启用调试模式") // 添加 debug 标志
	SignCmd.Flags().StringVarP(&secret_key, "secret_key", "k", "", "server酱密钥")

	// 标记必需的标志
	SignCmd.MarkFlagRequired("account")
	SignCmd.MarkFlagRequired("address")
}

// signIn 执行签到逻辑
func SignIn(account, address, address_name, latitude, longitude string) string {
	// 获取用户信息
	token, dbLatitude, dbLongitude, err := utils.GetUser(account)
	if err != nil || token == "" {
		if Config.debug {
			fmt.Printf("获取用户信息失败: %v\n", err)
		}
		fmt.Println("未找到该账号的 token，请先登录。")
		return "未找到该账号的 token，请先登录。"
	}

	// 如果命令行未提供 latitude 和 longitude，则使用数据库中的值
	if latitude == "" {
		Config.latitude = dbLatitude
		latitude = Config.latitude
	}
	if longitude == "" {
		Config.longitude = dbLongitude
		longitude = Config.longitude
	}

	if Config.latitude == "" || Config.longitude == "" {
		fmt.Println("未提供经纬度信息，且数据库中不存在，请先查询签到信息或手动提供经纬度。")
		return "未提供经纬度信息，且数据库中不存在，请先查询签到信息或手动提供经纬度。"
	}

	// 使用公钥加密 latitude 和 longitude
	encryptedLatitude, err := rsaEncrypt([]byte(Config.latitude))
	if err != nil {
		if Config.debug {
			fmt.Printf("加密纬度失败: %v\n", err)
		}
		fmt.Println("加密纬度失败:", err)
		return "加密纬度失败"
	}

	encryptedLongitude, err := rsaEncrypt([]byte(Config.longitude))
	if err != nil {
		if Config.debug {
			fmt.Printf("加密经度失败: %v\n", err)
		}
		fmt.Println("加密经度失败:", err)
		return "加密经度失败"
	}
	Config.address = address
	// 从 address 提取 province 和 city
	if address != "" {
		extractedProvince, extractedCity, err := extractProvinceAndCity(Config.address)
		if err != nil {
			if Config.debug {
				fmt.Printf("提取省份和城市失败: %v\n", err)
			}
			fmt.Println("地址格式不正确，无法提取省份和城市:", err)
			return "地址格式不正确，无法提取省份和城市"
		}
		Config.province = extractedProvince
		Config.city = extractedCity
	}

	apiURL := "https://api.xixunyun.com/signin_rsa"

	data := url.Values{}
	data.Set("address", Config.address)
	data.Set("province", Config.province)
	data.Set("city", Config.city)
	data.Set("latitude", encryptedLatitude)
	data.Set("longitude", encryptedLongitude)
	data.Set("remark", Config.remark)
	data.Set("comment", Config.comment)
	data.Set("address_name", address_name)
	data.Set("change_sign_resource", "0")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		if Config.debug {
			fmt.Printf("创建请求失败: %v\n", err)
		}
		fmt.Println("创建请求失败:", err)
		return "创建请求失败"
	}

	query := req.URL.Query()
	query.Add("token", token)
	query.Add("from", "app")
	query.Add("version", "5.1.3")
	query.Add("platform", "android")
	query.Add("entrance_year", "0")
	query.Add("graduate_year", "0")
	query.Add("school_id", u.schoolID) // 确保 school_id 已正确获取
	req.URL.RawQuery = query.Encode()

	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", "PHPSESSID=qkc555lu6050h43e204crialf0") // 注释掉不必要的 Cookie

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if Config.debug {
			fmt.Printf("发送 HTTP 请求失败: %v\n", err)
		}
		fmt.Println("请求失败:", err)
		return "请求失败"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if Config.debug {
			fmt.Printf("读取响应体失败: %v\n", err)
		}
		fmt.Println("读取响应体失败:", err)
		return "读取响应体失败"
	}

	if Config.debug {
		fmt.Printf("响应状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		if Config.debug {
			fmt.Printf("解析 JSON 失败: %v\n", err)
		}
		fmt.Println("解析响应数据失败:", err)
		return "解析响应数据失败"
	}

	// 检查响应码是否为 20000
	if code, ok := result["code"].(float64); !ok || code != 20000 {
		if Config.debug {
			fmt.Printf("签到失败，响应内容: %v\n", result)
		}
		fmt.Println("签到失败:", result["message"])
		if secret_key != "" {
			resultStr, err := json.Marshal(result)
			if err != nil {
				println("json解析错误")
			}
			PushMsgToWechat("签到失败", result["message"].(string)+"\r\n"+"错误详细信息："+string(resultStr), "9", secret_key)
		}
		return fmt.Sprintf("签到失败+%s", result["message"].(string))
	}

	fmt.Println("签到成功！")

	if secret_key != "" {
		Config.latitude, Config.longitude, err = utils.GetCoordinates(u.account)
		if err != nil {
			log.Fatalln(err)
		}
		PushMsgToWechat("签到成功", result["message"].(string)+"\n签到地址（使用高德地图查询）为"+Config.latitude+"\n"+Config.longitude, "9", secret_key)
	}
	return "签到成功"
}

// rsaEncrypt 使用提供的公钥对数据进行加密
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

	// 使用 base64 编码
	encryptedString := base64.StdEncoding.EncodeToString(encryptedData)
	return encryptedString, nil
}

// extractProvinceAndCity 从地址中提取省份和城市
func extractProvinceAndCity(address string) (string, string, error) {
	// 定义正则表达式，匹配省份和城市
	re := regexp.MustCompile(`(?P<province>[^省]+省)?(?P<city>[^市]+市)?`)
	matches := re.FindStringSubmatch(address)

	if len(matches) >= 3 {
		// 提取省份和城市
		province := strings.TrimSpace(matches[1]) // 去掉多余的空格
		city := strings.TrimSpace(matches[2])
		return province, city, nil
	}
	return "", "", fmt.Errorf("地址格式不正确，无法提取省份和城市")
}

// service/impl/query_service.go
package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"xixunyunsign/service"
	"xixunyunsign/utils" // Keep for GetAdditionalUserData for now, ideally refactor that too
)

// QueryServiceImpl implements the QueryService interface.
type QueryServiceImpl struct {
	userRepo service.UserRepository // Add UserRepository dependency
	// httpClient *http.Client      // Consider injecting http client if needed for external calls
}

// NewQueryService creates a new QueryService implementation.
func NewQueryService(repo service.UserRepository) service.QueryService {
	// Return the exported type
	return &QueryServiceImpl{
		userRepo: repo,
		// httpClient: &http.Client{}, // Initialize http client if needed
	}
}

// QuerySignInfo retrieves sign-in information for a user.
func (q *QueryServiceImpl) QuerySignInfo(account string) (map[string]interface{}, error) {
	// Use injected userRepo
	token, _, _, err := q.userRepo.GetUser(account)
	if err != nil || token == "" {
		return nil, fmt.Errorf("获取用户 token 失败或未找到账号 %s: %w", account, err)
	}

	// TODO: Refactor GetAdditionalUserData to be part of UserRepository or another service
	userData, err := utils.GetAdditionalUserData(account)
	if err != nil {
		return nil, fmt.Errorf("获取用户额外信息失败: %w", err)
	}

	apiURL := "https://api.xixunyun.com/signin40/homepage"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	qry := req.URL.Query()
	qry.Add("month_date", "2024-12")
	qry.Add("token", token)
	qry.Add("from", "app")
	qry.Add("version", "5.1.3")
	qry.Add("platform", "android")
	qry.Add("entrance_year", userData["entrance_year"])
	qry.Add("graduate_year", userData["graduation_year"])
	qry.Add("school_id", userData["school_id"])
	req.URL.RawQuery = qry.Encode()

	req.Header.Set("User-Agent", "okhttp/3.8.0")
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

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

	// 更新数据库中的经纬度信息
	signResourcesInfo, ok := data["sign_resources_info"].(map[string]interface{})
	if ok {
		midLatitude := fmt.Sprintf("%v", signResourcesInfo["mid_sign_latitude"])
		midLongitude := fmt.Sprintf("%v", signResourcesInfo["mid_sign_longitude"])
		// Use injected userRepo
		_ = q.userRepo.UpdateCoordinates(account, midLatitude, midLongitude) // Error handling might be needed
	}

	return data, nil
}

// SearchSchool searches for schools by name.
func (q *QueryServiceImpl) SearchSchool(schoolName string) ([]service.SchoolInfo, error) {
	// TODO: Ideally, SearchSchoolID in utils should return []service.SchoolInfo
	// or this logic should be moved entirely into the repository.
	// For now, we call the utils function and convert the type.
	utilsSchools, err := utils.SearchSchoolID(schoolName)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败: %w", err)
	}
	if len(utilsSchools) == 0 {
		return nil, nil // Or return an empty slice: []service.SchoolInfo{}
	}

	// Convert []utils.SchoolInfo to []service.SchoolInfo
	serviceSchools := make([]service.SchoolInfo, len(utilsSchools))
	for i, uSchool := range utilsSchools {
		serviceSchools[i] = service.SchoolInfo{ // Assuming fields match
			SchoolID:   uSchool.SchoolID,
			SchoolName: uSchool.SchoolName,
		}
	}

	return serviceSchools, nil
}

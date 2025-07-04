// service/query.go
package service

// QueryService 定义查询相关功能接口
type QueryService interface {
	// QuerySignInfo 查询签到信息并更新经纬度信息
	QuerySignInfo(account string) (map[string]interface{}, error)
	// SearchSchool 查询学校信息（模糊匹配）
	SearchSchool(schoolName string) ([]SchoolInfo, error) // Use local service.SchoolInfo
}

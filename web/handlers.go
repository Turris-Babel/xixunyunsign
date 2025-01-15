// web/handlers.go
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"xixunyunsign/cmd"
	"xixunyunsign/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the expected login payload
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	SchoolID string `json:"school_id" binding:"required"`
	//Token    string `json:"token" binding:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// handleLogin processes the login request
func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 调用cmd包中的登录逻辑
	token, err := cmd.Login(req.Account, req.Password, req.SchoolID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Message: "登录失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message: "登录成功",
		Token:   token,
	})
}

// QueryResponse represents the query response payload
type QueryResponse struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// handleQuery processes the query request
func handleQuery(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusBadRequest, QueryResponse{
			Message: "缺少 account 参数",
		})
		return
	}

	// 调用cmd包中的查询逻辑
	data, err := cmd.Query(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, QueryResponse{
			Message: "查询失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, QueryResponse{
		Message: "查询成功",
		Data:    data,
	})
}

// SignRequest represents the expected sign-in payload
type SignRequest struct {
	Account     string `json:"account" binding:"required"`
	Address     string `json:"address" binding:"required"`
	AddressName string `json:"address_name"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Remark      string `json:"remark"`
	Comment     string `json:"comment"`
	Province    string `json:"province"`
	City        string `json:"city"`
	SecretKey   string `json:"secret_key"`
}

// SignResponse represents the sign-in response payload
type SignResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// handleSign processes the sign-in request
func handleSign(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SignResponse{
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 调用cmd包中的签到逻辑
	success := cmd.SignIn(req.Account)
	//if success == false {
	//	c.JSON(http.StatusInternalServerError, SignResponse{
	//		Message: "签到失败",
	//		Error:   fmt.Sprintf("%v", success),
	//	})
	//	return
	//}
	if success {
		c.JSON(http.StatusOK, SignResponse{
			Message: "签到成功",
		})
	} else {
		c.JSON(http.StatusInternalServerError, SignResponse{
			Message: "签到失败",
			Error:   fmt.Sprintf("%v", success),
		})
	}
}

// SearchSchoolResponse represents the search school ID response payload
type SearchSchoolResponse struct {
	Message string             `json:"message"`
	Schools []utils.SchoolInfo `json:"schools,omitempty"`
	Error   string             `json:"error,omitempty"`
}

// handleSearchSchoolID processes the search school ID request
func handleSearchSchoolID(c *gin.Context) {
	schoolName := c.Query("school_name")
	if schoolName == "" {
		c.JSON(http.StatusBadRequest, SearchSchoolResponse{
			Message: "缺少 school_name 参数",
		})
		return
	}

	// 调用 cmd 包中的查询学校 ID 逻辑
	rawData, err := cmd.SearchSchoolID(schoolName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SearchSchoolResponse{
			Message: "查询失败",
			Error:   err.Error(),
		})
		return
	}

	// 解析 JSON 字符串为 []utils.SchoolInfo
	var schools []utils.SchoolInfo
	if err := json.Unmarshal([]byte(rawData), &schools); err != nil {
		c.JSON(http.StatusInternalServerError, SearchSchoolResponse{
			Message: "数据解析失败",
			Error:   err.Error(),
		})
		return
	}

	// 如果没有找到匹配的学校
	if len(schools) == 0 {
		c.JSON(http.StatusNotFound, SearchSchoolResponse{
			Message: "没有找到匹配的学校",
		})
		return
	}

	// 返回查询成功结果
	c.JSON(http.StatusOK, SearchSchoolResponse{
		Message: "查询成功",
		Schools: schools,
	})
}

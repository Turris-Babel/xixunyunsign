// web/handlers.go
package web

import (
	"net/http"
	// "xixunyunsign/cmd" // No longer calling cmd directly
	"xixunyunsign/service" // Depend on service interfaces
	// "xixunyunsign/utils" // No longer needed directly if SchoolInfo is in service

	"github.com/gin-gonic/gin"
)

// Handlers holds dependencies for web handlers, like services.
type Handlers struct {
	AuthService  service.AuthService
	QueryService service.QueryService
	SignService  service.SignService
	// TODO: Add SchoolSearch service if needed, or keep logic in QueryService
}

// NewHandlers creates a new Handlers struct with injected services.
func NewHandlers(authService service.AuthService, queryService service.QueryService, signService service.SignService) *Handlers {
	return &Handlers{
		AuthService:  authService,
		QueryService: queryService,
		SignService:  signService,
	}
}

// LoginRequest represents the expected login payload
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	SchoolID string `json:"school_id" binding:"required"`
	Token    string `json:"token"   // binding:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// handleLogin processes the login request - now a method on Handlers
func (h *Handlers) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 调用注入的 AuthService
	token, err := h.AuthService.Login(req.Account, req.Password, req.SchoolID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Message: "登录失败", // Consider more specific messages if possible
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

// handleQuery processes the query request - now a method on Handlers
func (h *Handlers) handleQuery(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusBadRequest, QueryResponse{
			Message: "缺少 account 参数",
		})
		return
	}

	// 调用注入的 QueryService
	// Assuming QueryService.QuerySignInfo is the correct method for this handler
	data, err := h.QueryService.QuerySignInfo(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, QueryResponse{
			Message: "查询失败", // Consider more specific messages
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

// handleSign processes the sign-in request - now a method on Handlers
func (h *Handlers) handleSign(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SignResponse{
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 调用注入的 SignService
	// Note: The original cmd.SignIn returned a string "签到成功" or an error string.
	// The service.SignIn interface returns (string, error). We adapt the logic.
	message, err := h.SignService.SignIn(req.Account, req.Address, req.AddressName, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SignResponse{
			Message: "签到失败",
			Error:   err.Error(), // Use the error returned by the service
		})
		return
	}

	// Assuming the service returns a success message string when err is nil
	c.JSON(http.StatusOK, SignResponse{
		Message: message, // Use the message from the service
	})
}

// SearchSchoolResponse represents the search school ID response payload
// Note: utils.SchoolInfo is used here. Consider moving SchoolInfo definition
// to the service package if it's primarily a service-level concept,
// or keep it in utils if it's a general utility struct.
type SearchSchoolResponse struct {
	Message string               `json:"message"`
	Schools []service.SchoolInfo `json:"schools,omitempty"` // Use service.SchoolInfo
	Error   string               `json:"error,omitempty"`
}

// handleSearchSchoolID processes the search school ID request - now a method on Handlers
func (h *Handlers) handleSearchSchoolID(c *gin.Context) {
	// 获取查询参数
	schoolName := c.Query("school_name")
	if schoolName == "" {
		c.JSON(http.StatusBadRequest, SearchSchoolResponse{
			Message: "缺少 school_name 参数",
		})
		return
	}

	// 调用注入的 QueryService's SearchSchool method
	// Assuming QueryService handles the school search logic now.
	// The original cmd.SearchSchoolID returned []utils.SchoolInfo, error.
	// The service.SearchSchool interface returns ([]service.SchoolInfo, error).
	// We need to ensure the types match or adapt. Let's assume service.SchoolInfo exists
	// and is compatible or identical to utils.SchoolInfo for now.
	// If not, we'd need to adjust the service interface or the handler response.
	schools, err := h.QueryService.SearchSchool(schoolName) // Using QueryService
	if err != nil {
		c.JSON(http.StatusInternalServerError, SearchSchoolResponse{
			Message: "查询失败", // Consider more specific messages
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
	// Adapt the response to use the type returned by the service.
	// Assuming service.SchoolInfo is compatible with utils.SchoolInfo for JSON marshalling.
	// If service.SchoolInfo is different, we might need a conversion step here.
	c.JSON(http.StatusOK, SearchSchoolResponse{
		Message: "查询成功",
		Schools: schools, // Use the result from QueryService.SearchSchool
	})
}

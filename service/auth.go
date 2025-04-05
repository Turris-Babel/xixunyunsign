// service/auth.go
package service

// AuthService 定义登录认证接口
type AuthService interface {
	Login(account, password, schoolID string) (string, error)
}

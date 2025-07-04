// service/sign.go
package service

// SignService 定义执行签到操作的接口
type SignService interface {
	SignIn(account, address, addressName, latitude, longitude string) (string, error)
}

package service

// UserRepository 定义了用户数据的存储库接口
type UserRepository interface {
	SaveUser(account, password, token, latitude, longitude, bindPhone, userNumber, userName string, schoolID float64, sex, className, entranceYear, graduationYear string) error
	GetUser(account string) (token, latitude, longitude string, err error)
	UpdateCoordinates(account, latitude, longitude string) error
}

// utils/user_repository.go
package utils

import (
	"database/sql"
	"xixunyunsign/service"
)

// UserRepositoryImpl implements the service.UserRepository interface.
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository 创建一个新的 UserRepository 实现
// 它接受一个数据库连接作为依赖项
func NewUserRepository(db *sql.DB) service.UserRepository {
	// Return the exported type
	return &UserRepositoryImpl{db: db}
}

// SaveUser saves or updates user data.
func (r *UserRepositoryImpl) SaveUser(account, password, token, latitude, longitude, bindPhone, userNumber, userName string, schoolID float64, sex, className, entranceYear, graduationYear string) error {
	// 使用注入的 db 连接
	insertSQL := `
    INSERT INTO users (account, password, token, latitude, longitude, bind_phone, user_number, user_name, school_id, sex, class_name, entrance_year, graduation_year)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(account) DO UPDATE SET 
        password = excluded.password,
        token = excluded.token,
        latitude = excluded.latitude,
        longitude = excluded.longitude,
        bind_phone = excluded.bind_phone,
        user_number = excluded.user_number,
        user_name = excluded.user_name,
        school_id = excluded.school_id,
        sex = excluded.sex,
        class_name = excluded.class_name,
        entrance_year = excluded.entrance_year,
        graduation_year = excluded.graduation_year;
    `
	_, err := r.db.Exec(insertSQL, account, password, token, latitude, longitude, bindPhone, userNumber, userName, schoolID, sex, className, entranceYear, graduationYear)
	return err
}

// GetUser retrieves user token and coordinates.
func (r *UserRepositoryImpl) GetUser(account string) (token, latitude, longitude string, err error) {
	// 使用注入的 db 连接
	querySQL := `SELECT token, latitude, longitude FROM users WHERE account = ?;`
	row := r.db.QueryRow(querySQL, account)
	err = row.Scan(&token, &latitude, &longitude)
	return
}

// UpdateCoordinates updates user coordinates.
func (r *UserRepositoryImpl) UpdateCoordinates(account, latitude, longitude string) error {
	// 使用注入的 db 连接
	updateSQL := `UPDATE users SET latitude = ?, longitude = ? WHERE account = ?;`
	_, err := r.db.Exec(updateSQL, latitude, longitude, account)
	return err
}

package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"sync"
	"xixunyunsign/service" // Import service package
)

var (
	db    *sql.DB
	once  sync.Once
	dbErr error
)

// ProvideDB initializes and returns a singleton database connection.
// It ensures that InitDB is called only once.
func ProvideDB() (*sql.DB, error) {
	once.Do(func() {
		// Set the database file path to "config.db"
		dbPath := "config.db"
		// Open the database connection
		db, dbErr = sql.Open("sqlite3", dbPath)
		if dbErr != nil {
			return // dbErr will be returned outside once.Do
		}

		// Run migrations/table creations
		dbErr = createTables(db)
		if dbErr != nil {
			db.Close() // Close connection if table creation fails
			db = nil   // Set db to nil to indicate failure
		}
	})
	return db, dbErr
}

// createTables creates the necessary database tables.
func createTables(db *sql.DB) error {
	// Create the users table
	createUserTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        account TEXT PRIMARY KEY,
        password TEXT,
        token TEXT,
        latitude TEXT,
        longitude TEXT,
        bind_phone TEXT,
        user_number TEXT,
        user_name TEXT,
        school_id REAL, -- Changed from INT to REAL to match SaveUser usage
        sex TEXT,
        class_name TEXT,
        entrance_year TEXT,
        graduation_year TEXT
    );
    `
	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		return fmt.Errorf("创建 users 表失败: %w", err)
	}

	// Create the school_info table (updated to include city_name and city_id)
	createSchoolTableSQL := `
    CREATE TABLE IF NOT EXISTS school_info (
        school_id TEXT PRIMARY KEY,
        school_name TEXT,
        city_name TEXT,  -- Added city_name
        city_id TEXT     -- Added city_id
    );
    `
	_, err = db.Exec(createSchoolTableSQL)
	if err != nil {
		return fmt.Errorf("创建 school_info 表失败: %w", err)
	}

	// Create schedules table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schedules (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        account TEXT,
        address TEXT,
        latitude TEXT,
        longitude TEXT,
        province TEXT,
        city TEXT,
        remark TEXT,
        comment TEXT,
        cron_expr TEXT,
        enabled INTEGER DEFAULT 1
    )`)
	if err != nil {
		return fmt.Errorf("创建 schedules 表失败: %w", err)
	}
	return nil
}

// SchoolInfo is now defined in the service package.
// type SchoolInfo struct {
// 	SchoolID   string `json:"school_id"`
// 	SchoolName string `json:"school_name"`
// }

// --- Functions below still use the global db or InitDB logic ---
// --- These should ideally be refactored to accept *sql.DB or moved to repositories ---
// --- For now, we focus on providing the DB via ProvideDB for wire ---

// SaveSchoolInfo saves or updates the school information in the database.
// TODO: Refactor to accept *sql.DB or move to a repository
func SaveSchoolInfo(cityName, cityID, schoolID, schoolName string) error {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	insertSQL := `
    INSERT INTO school_info (city_name, city_id, school_id, school_name)
    VALUES (?, ?, ?, ?)
    ON CONFLICT(school_id) DO UPDATE SET
        school_name = excluded.school_name,
        city_name = excluded.city_name,
        city_id = excluded.city_id;
    `

	_, err = localDB.Exec(insertSQL, cityName, cityID, schoolID, schoolName)
	return err
}

// FetchAndSaveSchoolData fetches the school data from the given API and saves it to the database.
// TODO: Refactor to accept *sql.DB or move to a repository/service
func FetchAndSaveSchoolData() error {
	// Make the GET request to the API
	resp, err := http.Get("https://api.xixunyun.com/login/schoolmap")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			CityName string               `json:"name"`
			CityId   string               `json:"id"`
			Schools  []service.SchoolInfo `json:"list"` // Use service.SchoolInfo
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Loop through the data and save each school info
	for _, group := range result.Data {
		for _, school := range group.Schools {
			if err := SaveSchoolInfo(group.CityName, group.CityId, school.SchoolID, school.SchoolName); err != nil {
				log.Printf("Error saving school %s: %v", school.SchoolName, err)
			}
		}
	}

	return nil
}

// SaveUser saves or updates a user in the database.
// NOTE: This function is now redundant as the logic is in UserRepository.
// Keeping it here temporarily might break things if it's called directly elsewhere.
// Ideally, remove this and update callers to use UserRepository.
// TODO: Remove this function or refactor callers.
func SaveUser(account, password, token, latitude, longitude, bindPhone, userNumber, userName string, schoolID float64, sex, className, entranceYear, graduationYear string) error {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	insertSQL := `
    INSERT INTO users (
        account, password, token, latitude, longitude, bind_phone, 
        user_number, user_name, school_id, sex, class_name, entrance_year, graduation_year
    )
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

	_, err = localDB.Exec(insertSQL, account, password, token, latitude, longitude, bindPhone, userNumber, userName, schoolID, sex, className, entranceYear, graduationYear)
	return err
}

// GetUser retrieves user information from the database by account.
// NOTE: This function is now redundant as the logic is in UserRepository.
// TODO: Remove this function or refactor callers.
func GetUser(account string) (token, latitude, longitude string, err error) {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return "", "", "", fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return "", "", "", fmt.Errorf("数据库连接未初始化")
	}
	querySQL := `SELECT token, latitude, longitude FROM users WHERE account = ?;`
	row := localDB.QueryRow(querySQL, account)
	err = row.Scan(&token, &latitude, &longitude)
	return
}

// CloseDB closes the database connection.
// Note: Closing the singleton DB might affect other parts of the application.
// Consider managing the DB lifecycle elsewhere (e.g., in main).
func CloseDB() error {
	if db != nil {
		err := db.Close()
		db = nil // Reset db variable after closing
		// Reset sync.Once if you need to re-initialize later (complex scenario)
		// once = sync.Once{}
		return err
	}
	return nil
}

// UpdateCoordinates updates the latitude and longitude for a given account.
// NOTE: This function is now redundant as the logic is in UserRepository.
// TODO: Remove this function or refactor callers.
func UpdateCoordinates(account, latitude, longitude string) error {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	updateSQL := `
        UPDATE users
        SET latitude = ?, longitude = ?
        WHERE account = ?;
    `
	_, err = localDB.Exec(updateSQL, latitude, longitude, account)
	return err
}

// TODO: Refactor to accept *sql.DB or move to a repository/service
func GetCoordinates(account string) (string, string, error) {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return "", "", fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return "", "", fmt.Errorf("数据库连接未初始化")
	}
	querySQL := `
        SELECT latitude, longitude
        FROM users
        WHERE account = ?;
    `
	var latitude, longitude string
	err = localDB.QueryRow(querySQL, account).Scan(&latitude, &longitude)
	if err != nil {
		return "", "", err
	}
	return latitude, longitude, nil
}

// GetAdditionalUserData retrieves additional user data for constructing the query parameters.
// TODO: Refactor to accept *sql.DB or move to a repository/service
func GetAdditionalUserData(account string) (map[string]string, error) {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	querySQL := `
        SELECT entrance_year, graduation_year, school_id
        FROM users
        WHERE account = ?;
    `
	row := localDB.QueryRow(querySQL, account)

	var entranceYear, graduateYear, schoolID string // Assuming schoolID is stored as TEXT or can be scanned into string
	err = row.Scan(&entranceYear, &graduateYear, &schoolID)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"entrance_year":   entranceYear,
		"graduation_year": graduateYear,
		"school_id":       schoolID,
	}, nil
}

// SearchSchoolID searches for all school IDs by school name using fuzzy matching.
// TODO: Refactor to accept *sql.DB or move to a repository/service
func SearchSchoolID(schoolName string) ([]service.SchoolInfo, error) { // Return service.SchoolInfo
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// Use SQL LIKE for fuzzy matching
	querySQL := `SELECT school_id, school_name FROM school_info WHERE school_name LIKE ?;`
	likeName := "%" + schoolName + "%"

	rows, err := localDB.Query(querySQL, likeName)
	if err != nil {
		return nil, fmt.Errorf("查询学校ID时发生错误: %w", err)
	}
	defer rows.Close()

	var schools []service.SchoolInfo // Use service.SchoolInfo
	for rows.Next() {
		var school service.SchoolInfo // Use service.SchoolInfo
		if err := rows.Scan(&school.SchoolID, &school.SchoolName); err != nil {
			return nil, fmt.Errorf("读取查询结果时发生错误: %w", err)
		}
		schools = append(schools, school)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历查询结果时发生错误: %w", err)
	}

	return schools, nil
}

// TODO: Refactor to accept *sql.DB or move to a repository/service
func IsSchoolInfoTableEmpty() (bool, error) {
	localDB, err := ProvideDB() // Use ProvideDB to ensure initialization
	if err != nil {
		return false, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if localDB == nil {
		return false, fmt.Errorf("数据库连接未初始化")
	}
	var count int
	querySQL := `SELECT COUNT(*) FROM school_info;`
	err = localDB.QueryRow(querySQL).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("查询 school_info 表时发生错误: %w", err)
	}
	return count == 0, nil
}

package logger

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseLogWriter 数据库日志写入器
type DatabaseLogWriter struct {
	db        *sql.DB
	dbType    string
	insertSQL string
}

// NewDatabaseLogWriter 创建数据库日志写入器
func NewDatabaseLogWriter(dbType, connStr string) (*DatabaseLogWriter, error) {
	var db *sql.DB
	var err error
	var createTableSQL string
	var insertSQL string

	switch dbType {
	case "postgres":
		db, err = sql.Open("postgres", connStr)
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS request_logs (
				id SERIAL PRIMARY KEY,
				request_id VARCHAR(36),
				log_type VARCHAR(10),
				log_time TIMESTAMP,
				log_data JSONB
			)`
		insertSQL = `
			INSERT INTO request_logs (request_id, log_type, log_time, log_data) 
			VALUES ($1, $2, $3, $4)`

	case "mysql":
		db, err = sql.Open("mysql", connStr)
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS request_logs (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				request_id VARCHAR(36),
				log_type VARCHAR(10),
				log_time TIMESTAMP,
				log_data JSON
			)`
		insertSQL = `
			INSERT INTO request_logs (request_id, log_type, log_time, log_data) 
			VALUES (?, ?, ?, ?)`

	case "sqlite":
		db, err = sql.Open("sqlite3", connStr)
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS request_logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				request_id TEXT,
				log_type TEXT,
				log_time TIMESTAMP,
				log_data TEXT
			)`
		insertSQL = `
			INSERT INTO request_logs (request_id, log_type, log_time, log_data) 
			VALUES (?, ?, ?, ?)`

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 创建日志表
	if _, err = db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &DatabaseLogWriter{
		db:        db,
		dbType:    dbType,
		insertSQL: insertSQL,
	}, nil
}

func (w *DatabaseLogWriter) Write(log map[string]interface{}) error {
	logJSON, _ := json.Marshal(log)
	var logJSONStr string

	// SQLite 需要将 JSON 转换为字符串
	if w.dbType == "sqlite" {
		logJSONStr = string(logJSON)
	} else {
		logJSONStr = string(logJSON)
	}

	_, err := w.db.Exec(
		w.insertSQL,
		log["request_id"],
		log["type"],
		log["time"],
		logJSONStr,
	)
	return err
}

func (w *DatabaseLogWriter) Close() error {
	return w.db.Close()
}

package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/olivermgi/golang-andromeda-api/config"
)

var DB *sql.DB

func init() {
	dbConfig := config.GetDatabaseConfig()
	driver := config.GetDatabaseDriver()

	conn := ""
	if driver == "mysql" {
		conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig["username"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["database"])
	} else if driver == "postgres" {
		conn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbConfig["username"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["database"])
	}

	var err error
	DB, err = sql.Open(driver, conn)
	if err != nil {
		log.Fatalln("資料庫設定錯誤，錯誤資訊：", err)
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	ctx := context.Background()
	err = DB.PingContext(ctx)
	if err != nil {
		log.Fatalln("資料庫連線失敗，錯誤資訊：", err)
	}

	log.Println("資料庫成功連線")
}

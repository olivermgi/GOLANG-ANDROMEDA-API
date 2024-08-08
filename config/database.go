package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("讀取不到 .env 檔")
	}
}

func GetMysqlConfig() map[string]string {
	return map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"database": os.Getenv("DB_DATABASE"),
		"username": os.Getenv("DB_USERNAME"),
		"password": os.Getenv("DB_PASSWORD"),
	}
}

//             'host' => env('DB_HOST', '127.0.0.1'),
//             'port' => env('DB_PORT', '3306'),
//             'database' => env('DB_DATABASE', 'forge'),
//             'username' => env('DB_USERNAME', 'forge'),

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

func GetGcpConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"common": {
			"project_id": os.Getenv("GOOGLE_GCP_PROJECT_ID"),
		},
		"storage": {
			"bucketName": os.Getenv("GOOGLE_GCP_STORAGE_BUCKET_NAME"),
		},
		"transcoder": {
			"location":   os.Getenv("GOOGLE_GCP_TRANSCODER_LOCATION"),
			"input_uri":  "gs://" + os.Getenv("GOOGLE_GCP_STORAGE_BUCKET_NAME") + os.Getenv("GOOGLE_GCP_TRANSCODER_INPUT_PATH"),
			"output_uri": "gs://" + os.Getenv("GOOGLE_GCP_STORAGE_BUCKET_NAME") + os.Getenv("GOOGLE_GCP_TRANSCODER_OUTPUT_PATH"),
		},
	}
}

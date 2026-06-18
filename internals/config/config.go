package config

import "os"

var (
	Port       = os.Getenv("PORT")
	DbFilePath = os.Getenv("DB_FILE_PATH")
)

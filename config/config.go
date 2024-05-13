package config

import (
	"crowdfunding-minpro-alterra/database"
	"os"
)

func InitConfigDB() database.Config {
	return database.Config{
		DBName: os.Getenv("DBName"),
		DBUser: os.Getenv("DBUser"),
		DBPass: os.Getenv("DBPass"),
		DBHost: os.Getenv("DBHost"),
		DBPort: os.Getenv("DBPort"),
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var loggingToDB = false
var FILEBUS_URL = "https://filebus.nova.moe/"
var DB_HOST = "localhost"
var DB_USERNAME = "root"
var DB_PASSWORD = "password"
var DB_DBNAME = "filebus"
var db *gorm.DB

func main() {

	logToDB := os.Getenv("ENABLE_LOG")
	FILEBUS_URL = os.Getenv("FILEBUS_URL")

	if logToDB == "TRUE" {
		loggingToDB = true
		DB_HOST = os.Getenv("DB_HOST")
		DB_USERNAME = os.Getenv("DB_USERNAME")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_DBNAME = os.Getenv("DB_DBNAME")

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_DBNAME)
		db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024 * 1024, // Limit to 1TB
	})

	app.Get("/delete/*", handle_delete)
	app.Post("/upload", handle_upload)
	app.Static("/", "/data")

	app.Listen(":3000")
}

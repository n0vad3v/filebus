package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DeleteLog struct {
	Filename  string
	Filepath  string
	DeleterIp string
	DeletedAt time.Time
}

func logDeleteRequestToDB(db *gorm.DB, fileName string, filePath string, deleterIp string) {

	delete := DeleteLog{
		Filename:  fileName,
		Filepath:  filePath,
		DeleterIp: deleterIp,
		DeletedAt: time.Now(),
	}

	result := db.Create(&delete)
	fmt.Println(result)
}

func handle_delete(c *fiber.Ctx) error {
	fileNameWithPath := "/data/" + c.Params("*")
	err := os.Remove(fileNameWithPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{err.Error(): err.Error()})
	}
	if loggingToDB == true {
		logDeleteRequestToDB(db, c.Params("*"), fileNameWithPath, c.IP())
	}
	return c.Status(200).JSON(fiber.Map{"message": "File deleted successfully"})
}

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type UploadLog struct {
	Filename   string
	Filepath   string
	Filesize   int64
	Filehash   string
	UploaderIp string
	UploadedAt time.Time
}

func logUploadRequestToDB(db *gorm.DB, fileName string, filePath string, fileSize int64, uploaderIp string, fileHash string) {

	upload := UploadLog{
		Filename:   fileName,
		Filepath:   filePath,
		Filesize:   fileSize,
		Filehash:   fileHash,
		UploaderIp: uploaderIp,
		UploadedAt: time.Now(),
	}

	result := db.Create(&upload)
	fmt.Println(result)
}

func getFileHash(filePath string) string {
	fileContent, err := os.Open(filePath)
	defer fileContent.Close()
	hash := md5.New()
	_, err = io.Copy(hash, fileContent)
	hashContent := fmt.Sprintf("%x", hash.Sum(nil))
	if err != nil {
		return ""
	}
	return hashContent
}

func handle_upload(c *fiber.Ctx) error {
	// filePathWithFilename -> path/to/main.py
	form, err := c.MultipartForm()

	filePathWithFilename := ""

	if err != nil {
		fmt.Println(err)
	} else {
		var result map[string]interface{}
		mapstructure.Decode(form, &result)
		for filename, _ := range result["File"].(map[string][]*multipart.FileHeader) {
			filePathWithFilename = filename
		}
	}

	fmt.Println("filePathWithFilename:", filePathWithFilename)

	if filePathWithFilename != "file" {
		fmt.Println("Old upload is recorded, filePathWithFilename is: ", filePathWithFilename)

		file, _ := c.FormFile(filePathWithFilename)
		// Create path for filename
		// filepath.Dir(filePathWithFilename) -> path/to
		os.MkdirAll(filepath.Dir(filePathWithFilename), os.ModePerm)
		c.SaveFile(file, fmt.Sprintf("/data/%s", filePathWithFilename))

		if loggingToDB == true {
			logUploadRequestToDB(db, file.Filename, filePathWithFilename, file.Size, c.IP(), getFileHash(filePathWithFilename))
		}

		return c.JSON(fiber.Map{
			"url": FILEBUS_URL + filePathWithFilename,
			"md5": getFileHash(filePathWithFilename),
		})

	} else {
		filePathWithFilename := c.FormValue("filepath")
		file, _ := c.FormFile("file")
		fmt.Println("New upload is recorded, filePathWithFilename is: ", filePathWithFilename)

		// Create path for filename
		// filepath.Dir(filePathWithFilename) -> path/to
		os.MkdirAll(filepath.Dir(filePathWithFilename), os.ModePerm)
		c.SaveFile(file, fmt.Sprintf("/data/%s", filePathWithFilename))

		if loggingToDB == true {
			logUploadRequestToDB(db, file.Filename, filePathWithFilename, file.Size, c.IP(), getFileHash(filePathWithFilename))
		}

		return c.JSON(fiber.Map{
			"url": FILEBUS_URL + filePathWithFilename,
			"md5": getFileHash(filePathWithFilename),
		})
	}

}

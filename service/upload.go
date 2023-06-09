package service

import (
	"fmt"
	"os"
	"project_absensi/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadFormData(imagePath string, c *gin.Context) (string, error) {
	_ = os.MkdirAll(imagePath, os.ModePerm)
	File, err := c.FormFile("image")
	if err != nil {
		return "", err
	}
	extension := strings.Split(File.Filename, ".")
	fileName := helper.RandomString(20)
	filePath := fmt.Sprintf("%s%s.%s", imagePath, fileName, extension[1])
	for helper.FileExists(filePath) {
		fileName = helper.RandomString(20)
		filePath = fmt.Sprintf("%s%s%s", imagePath, fileName, extension[1])
	}

	err = c.SaveUploadedFile(File, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

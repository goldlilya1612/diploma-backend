package media

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	imagesDir = "./resources/images/"
)

func ImageUpload(ctx *gin.Context) (imageName string, imagePath string, url string, err error) {

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		return "", "", "", err
	}

	fileExt := filepath.Ext(header.Filename)
	originalImageName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	imageName = strings.ReplaceAll(strings.ToLower(originalImageName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	url = "http://localhost:8080/api/courses/images/" + imageName //  TODO: подставлять урл из настроек сервера

	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			return "", "", "", err
		}
	}

	imagePath = imagesDir + imageName
	out, err := os.Create(imagePath)
	if err != nil {
		return "", "", "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", "", "", err
	}

	return imageName, imagePath, url, nil
}

func RemoteUpload(ctx *gin.Context) {

}

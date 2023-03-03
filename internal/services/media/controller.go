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

func FileUpload(ctx *gin.Context) (string, error) {

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	url := "http://localhost:8080/api/courses/images/" + filename

	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	out, err := os.Create(imagesDir + filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return url, nil
}

func RemoteUpload(ctx *gin.Context) {

}

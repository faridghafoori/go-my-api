package controllers

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/services"
	"gin-mongo-api/utils"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		extention := strings.Split(header.Header["Content-Type"][0], "/")[1]
		filename := utils.GetMD5Hash(header.Filename) + "." + extention
		sise := header.Size

		folderPrefix := filename[:3]

		_, err = os.Stat("public/" + folderPrefix)
		if os.IsNotExist(err) {
			err := os.Mkdir("public/"+folderPrefix, 0755)
			utils.GenerateErrorOutput(http.StatusBadRequest, err, c)
		}

		out, err := os.Create("public/" + folderPrefix + "/" + filename)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		defer out.Close()
		_, err = io.Copy(out, file)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		filepath := configs.ENV_RUNABLE_PROJECT_URI() + "/file/" + folderPrefix + "/" + filename

		utils.GenerateSuccessOutput(
			map[string]interface{}{
				"filepath": filepath,
				"size":     sise,
			},
			c,
		)
	}
}

func UploadToMinio() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)
		files := form.File["file"]

		var bucketName string = "images"
		if c.Request.URL.Query().Get("bucket") != "" {
			bucketName = c.Request.URL.Query().Get("bucket")
		}

		var extention string = "/default"
		if c.Request.URL.Query().Get("extention") != "" {
			extention = c.Request.URL.Query().Get("extention")
		}

		var links []string
		for _, file := range files {
			link := services.UploadFile(c, bucketName, extention, file)
			links = append(links, link)
		}

		utils.GenerateSuccessOutput(
			map[string]interface{}{
				"links": links,
			},
			c,
		)
	}
}

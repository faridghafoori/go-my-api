package controllers

import (
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

		filepath := "http://localhost:6000/file/" + folderPrefix + "/" + filename

		utils.GenerateSuccessOutput(
			map[string]interface{}{
				"filepath": filepath,
				"size":     sise,
			},
			c,
		)
	}
}

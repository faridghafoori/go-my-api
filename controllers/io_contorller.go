package controllers

import (
	"fmt"
	"gin-mongo-api/responses"
	"gin-mongo-api/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
			return
		}
		extention := strings.Split(header.Header["Content-Type"][0], "/")[1]
		filename := utils.GetMD5Hash(header.Filename) + "." + extention
		sise := header.Size

		folderPrefix := filename[:3]

		folderInfo, err := os.Stat("public/" + folderPrefix)
		if os.IsNotExist(err) {
			log.Fatal("Folder does not exist.")
			// err := os.Mkdir(folderPrefix, 0755)
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
		log.Println(folderInfo)

		out, err := os.Create("public/" + folderPrefix + "/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		filepath := "http://localhost:6000/file/" + folderPrefix + "/" + filename
		c.JSON(
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data: map[string]interface{}{
					"filepath": filepath,
					"size":     sise,
				},
			},
		)
	}
}

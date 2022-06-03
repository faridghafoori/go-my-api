package utils

import (
	"crypto/md5"
	"encoding/hex"
	"gin-mongo-api/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ValidateStruct(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validationErr := validate.Struct(&model); validationErr != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.GeneralResponse{
					Status:  http.StatusBadRequest,
					Message: ErrorMessage,
					Data:    validationErr.Error(),
				},
			)
			return
		}
	}

}

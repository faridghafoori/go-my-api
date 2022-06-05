package utils

import (
	"crypto/md5"
	"encoding/hex"
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
		validationErr := validate.Struct(&model)
		GenerateErrorOutput(http.StatusBadRequest, validationErr, c)
	}
}

package utils

import (
	"gin-mongo-api/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateErrorOutput(statusCode int, e error, c *gin.Context, optionalParams ...map[string]interface{}) {
	if e != nil {
		var showableMessage string = http.StatusText(statusCode)
		var showableData interface{} = e.Error()
		if len(optionalParams) > 0 {
			if optionalParams[0]["data"] != nil {
				showableData = optionalParams[0]["data"]
			}
			if optionalParams[0]["message"] != nil {
				showableMessage = optionalParams[0]["message"].(string)
			}
		}
		c.JSON(
			statusCode,
			responses.GeneralResponse{
				Status:  statusCode,
				Message: showableMessage,
				Data:    showableData,
			},
		)
		c.Abort()
		panic(e)
	}
}

func GenerateSuccessOutput(result interface{}, c *gin.Context, statusCode ...int) {
	var showableStatucCode int = http.StatusOK
	if len(statusCode) > 0 {
		showableStatucCode = statusCode[0]
	}
	c.JSON(
		showableStatucCode,
		responses.GeneralResponse{
			Status:  showableStatucCode,
			Message: http.StatusText(showableStatucCode),
			Data:    result,
		},
	)
	c.Abort()
}

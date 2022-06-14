package middleware

import (
	"gin-mongo-api/services"
	"gin-mongo-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(tokenType ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestTokenType string
		if tokenType != nil {
			requestTokenType = tokenType[0]
		} else {
			requestTokenType = ""
		}
		err := services.TokenValid(c.Request, requestTokenType)
		utils.GenerateErrorOutput(
			http.StatusUnauthorized,
			err,
			c,
		)
		if err != nil {
			return
		}
		c.Next()
	}
}

package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/services"
	"gin-mongo-api/utils"
	"image/png"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inputUser models.Authentication
		bindError := c.BindJSON(&inputUser)
		utils.GenerateErrorOutput(http.StatusUnprocessableEntity, bindError, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		var user models.User
		findError := userCollection.FindOne(ctx, bson.M{"username": inputUser.Username}).Decode(&user)
		utils.GenerateErrorOutput(http.StatusBadRequest, findError, c)

		inputUser.Password = utils.GetSHA256Hash(inputUser.Password)
		if user.Username != inputUser.Username || user.Password != inputUser.Password {
			utils.GenerateErrorOutput(
				http.StatusUnauthorized,
				errors.New(""),
				c,
				map[string]interface{}{
					"message": "Invalid Username or Password",
					"data":    "Please provide valid login details",
				},
			)
		}

		if !user.TotpActive {
			ts, err := services.CreateToken(user.Id.Hex())
			utils.GenerateErrorOutput(http.StatusUnprocessableEntity, err, c)

			saveErr := services.CreateAuth(user.Id.Hex(), ts)
			utils.GenerateErrorOutput(http.StatusUnprocessableEntity, saveErr, c)

			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}

			utils.GenerateSuccessOutput(tokens, c)
			return
		}

		ts, err := services.CreateTOTPToken(user.Id.Hex())
		utils.GenerateErrorOutput(http.StatusUnprocessableEntity, err, c)

		tokens := map[string]string{
			"token": ts.AccessToken,
		}

		if user.TotpKey == "" {
			utils.GenerateSuccessOutput(
				map[string]interface{}{
					"message": "Your TOTP key is not set, must be set",
					"data":    tokens,
				},
				c,
				http.StatusCreated,
			)
		}

		utils.GenerateSuccessOutput(tokens, c, http.StatusAccepted)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := services.ExtractTokenMetadata(c.Request)
		utils.GenerateErrorOutput(http.StatusUnprocessableEntity, err, c)

		deleted, delErr := services.DeleteAuth(au.AccessUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			utils.GenerateErrorOutput(http.StatusUnprocessableEntity, delErr, c)
		}

		utils.GenerateSuccessOutput("Successfully logged out", c)
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		bindErr := c.ShouldBindJSON(&mapToken)
		utils.GenerateErrorOutput(http.StatusUnprocessableEntity, bindErr, c)

		refreshToken := mapToken["refresh_token"]

		//verify the token
		token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(configs.ENV_JWT_REFRESH_SECRET()), nil
		})
		//if there is an error, the token must have expired
		utils.GenerateErrorOutput(
			http.StatusUnauthorized,
			err,
			c,
			map[string]interface{}{
				"data": "Refresh token expired",
			},
		)

		//is token valid?
		validError := token.Claims.Valid()
		utils.GenerateErrorOutput(http.StatusUnauthorized, validError, c)

		//Since token is valid, get the uuid:
		claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
		if ok && token.Valid {
			refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
			if !ok {
				utils.GenerateErrorOutput(
					http.StatusUnprocessableEntity,
					errors.New("token claims maybe damaged"),
					c,
				)
			}
			userId := claims["sub"]
			//Delete the previous Refresh Token
			deleted, delErr := services.DeleteAuth(refreshUuid)
			// if any goes wrong
			if delErr != nil || deleted == 0 {
				utils.GenerateErrorOutput(
					http.StatusUnauthorized,
					errors.New("Refresh token not valid"),
					c,
				)
			}
			//Create new pairs of refresh and access tokens
			ts, createErr := services.CreateToken(userId.(string))
			utils.GenerateErrorOutput(
				http.StatusForbidden,
				createErr,
				c,
			)
			//save the tokens metadata to redis
			saveErr := services.CreateAuth(userId.(string), ts)
			utils.GenerateErrorOutput(
				http.StatusForbidden,
				saveErr,
				c,
			)

			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}

			utils.GenerateSuccessOutput(tokens, c)
		} else {
			utils.GenerateErrorOutput(
				http.StatusUnauthorized,
				errors.New(""),
				c,
				map[string]interface{}{
					"data": "Refresh expired",
				},
			)
		}
	}
}

func TOTPGenerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		AccessDetails, err := services.ExtractTOTPTokenMetadata(c.Request)
		utils.GenerateErrorOutput(http.StatusUnauthorized, err, c)

		user, err := services.FetchUser(AccessDetails.UserId)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "mygoapi.com",
			AccountName: user.Username,
		})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)
		user.TotpKey = key.Secret()

		_, err = userCollection.UpdateOne(ctx, bson.M{"id": user.Id}, bson.M{"$set": user})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		// Convert TOTP key into a PNG
		var buf bytes.Buffer
		img, err := key.Image(200, 200)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)
		png.Encode(&buf, img)

		// display the QR code to the user.
		qrCodeLink := services.Display(key, buf.Bytes())
		utils.GenerateSuccessOutput(qrCodeLink, c)
	}
}

func VerifyTOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		AccessDetails, err := services.ExtractTOTPTokenMetadata(c.Request)
		utils.GenerateErrorOutput(http.StatusUnauthorized, err, c)

		user, err := services.FetchUser(AccessDetails.UserId)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		// Now Validate that the user's successfully added the passcode.
		var verifyTotp models.VerifyTOTP
		err = c.BindJSON(&verifyTotp)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&verifyTotp)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		valid := totp.Validate(verifyTotp.Passcode, user.TotpKey)
		if valid {
			ts, err := services.CreateToken(user.Id.Hex())
			utils.GenerateErrorOutput(http.StatusUnprocessableEntity, err, c)

			saveErr := services.CreateAuth(user.Id.Hex(), ts)
			utils.GenerateErrorOutput(http.StatusUnprocessableEntity, saveErr, c)

			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}

			utils.GenerateSuccessOutput(tokens, c)
		} else {
			utils.GenerateErrorOutput(http.StatusNotAcceptable, errors.New("TOTP code was wrong"), c)
		}

	}
}

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

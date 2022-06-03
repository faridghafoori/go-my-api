package controllers

import (
	"context"
	"fmt"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"gin-mongo-api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inputUser models.Authentication
		if err := c.BindJSON(&inputUser); err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}

		//use the validator library to validate required fields
		utils.ValidateStruct(&inputUser)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"username": inputUser.Username}).Decode(&user)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.GeneralResponse{
					Status:  http.StatusInternalServerError,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}

		inputUser.Password = utils.GetMD5Hash(inputUser.Password)
		if user.Username != inputUser.Username || user.Password != inputUser.Password {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: "Invalid Username or Password",
					Data:    "Please provide valid login details",
				},
			)
			return
		}

		ts, err := CreateToken(user.Id.Hex())
		if err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
		}

		saveErr := CreateAuth(user.Id.Hex(), ts)
		if saveErr != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		c.JSON(
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    tokens,
			},
		)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}
		deleted, delErr := DeleteAuth(au.AccessUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    "Successfully logged out",
			},
		)
	}
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		if err := c.ShouldBindJSON(&mapToken); err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.GeneralResponse{
					Status:  http.StatusUnprocessableEntity,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}
		refreshToken := mapToken["refresh_token"]

		//verify the token
		token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(configs.EnvJWTRefreshSecret()), nil
		})
		//if there is an error, the token must have expired
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: utils.ErrorMessage,
					Data:    "Refresh token expired",
				},
			)
			return
		}
		//is token valid?
		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}
		//Since token is valid, get the uuid:
		claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
		if ok && token.Valid {
			refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
			if !ok {
				c.JSON(
					http.StatusUnprocessableEntity,
					responses.GeneralResponse{
						Status:  http.StatusUnprocessableEntity,
						Message: utils.ErrorMessage,
						Data:    err.Error(),
					},
				)
				return
			}
			userId := claims["sub"]
			//Delete the previous Refresh Token
			deleted, delErr := DeleteAuth(refreshUuid)
			if delErr != nil || deleted == 0 { //if any goes wrong
				c.JSON(
					http.StatusUnauthorized,
					responses.GeneralResponse{
						Status:  http.StatusUnauthorized,
						Message: utils.UnauthorizedMessage,
						Data:    err,
					},
				)
				return
			}
			//Create new pairs of refresh and access tokens
			ts, createErr := CreateToken(userId.(string))
			if createErr != nil {
				c.JSON(
					http.StatusForbidden,
					responses.GeneralResponse{
						Status:  http.StatusForbidden,
						Message: utils.ForbidenMessage,
						Data:    createErr.Error(),
					},
				)
				return
			}
			//save the tokens metadata to redis
			saveErr := CreateAuth(userId.(string), ts)
			if saveErr != nil {
				c.JSON(
					http.StatusForbidden,
					responses.GeneralResponse{
						Status:  http.StatusForbidden,
						Message: utils.ForbidenMessage,
						Data:    createErr.Error(),
					},
				)
				return
			}
			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}
			c.JSON(
				http.StatusCreated,
				responses.GeneralResponse{
					Status:  http.StatusCreated,
					Message: utils.SuccessMessage,
					Data:    tokens,
				},
			)
		} else {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: utils.UnauthorizedMessage,
					Data:    "refresh expired",
				},
			)
		}
	}
}

func CreateToken(userId string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["sub"] = userId
	atClaims["exp"] = td.AtExpires
	atClaims["iss"] = "auth"
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(configs.EnvJWTAcessSecret()))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["sub"] = userId
	rtClaims["exp"] = td.RtExpires
	atClaims["iss"] = "auth"
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(configs.EnvJWTRefreshSecret()))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userid string, td *models.TokenDetails) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := configs.RDB.Set(ctx, td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := configs.RDB.Set(ctx, td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.EnvJWTAcessSecret()), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     claims["sub"].(string),
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *models.AccessDetails) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userid, err := configs.RDB.Get(ctx, authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deleted, err := configs.RDB.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: utils.UnauthorizedMessage,
					Data:    err.Error(),
				},
			)
			c.Abort()
			return
		}
		c.Next()
	}
}

package services

import (
	"encoding/base32"
	"gin-mongo-api/utils"
	"io/ioutil"
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Display(key *otp.Key, data []byte) string {
	fileName := key.AccountName() + "-" + "totp-qr-code"
	filePath := "public/" + fileName + ".png"
	ioutil.WriteFile(filePath, data, 0644)
	fileInLocal, _ := os.Open(filePath)
	multiPartFile, _ := utils.GetFileHeader(fileInLocal)
	link := UploadFile("totp-images", "/default", multiPartFile, fileName)
	return link
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom paramters
func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA256,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

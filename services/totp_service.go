package services

import (
	"encoding/base32"
	"gin-mongo-api/configs"
	"io/ioutil"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Display(key *otp.Key, data []byte) string {
	// fmt.Printf("Issuer:       %s\n", key.Issuer())
	// fmt.Printf("Account Name: %s\n", key.AccountName())
	// fmt.Printf("Secret:       %s\n", key.Secret())
	// fmt.Println("Writing PNG to qr-code.png....")
	file := "public/totp_codes/" + key.AccountName() + "-qr-code.png"
	ioutil.WriteFile(file, data, 0644)
	link := configs.ENV_RUNABLE_PROJECT_URI() + "/file/totp_codes/" + key.AccountName() + "-qr-code.png"
	// fmt.Println("")
	// fmt.Println("Please add your TOTP to your OTP Application now!")
	// fmt.Println("")
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

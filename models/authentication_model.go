package models

type Authentication struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type VerifyTOTP struct {
	Passcode string `validate:"required"`
}

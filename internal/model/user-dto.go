package model

import (
	"fmt"
	"regexp"
	"unicode"
)

type SignUpUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	validUsername = `^[\w]{6,50}$`
	validEmail    = `^[\w-\.]{6,54}@([\w-]{1,40}\.)[\w-]{2,4}$`
)

func (u *SignUpUserDTO) Validate() error {
	isMatched, err := regexp.MatchString(validUsername, u.Username)
	if err != nil {
		return err
	}
	if !isMatched {
		return fmt.Errorf("username must consist of letters or numbers, also it must contain from 6 to 50 characters")
	}
	isMatched, err = regexp.MatchString(validEmail, u.Email)
	if err != nil {
		return err
	}
	if !isMatched {
		return fmt.Errorf("email must consist of letters and numbers, also it mustn't exceed 100 characters")
	}
	return u.validatePassword()
}

func (u *SignUpUserDTO) validatePassword() error {
	var (
		hasDigit      bool
		hasLowerAlpha bool
		hasUpperAlpha bool
	)
	for _, sym := range u.Password {
		switch {
		case unicode.IsDigit(sym):
			hasDigit = true
		case unicode.IsLower(sym):
			hasLowerAlpha = true
		case unicode.IsUpper(sym):
			hasUpperAlpha = true
		}
	}
	length := len(u.Password)
	if hasDigit && hasLowerAlpha && hasUpperAlpha && length >= 8 && length <= 30 {
		return nil
	}
	return fmt.Errorf("password must contain from 8 to 30 characters, be at least one uppercase letter, one lowercase letter and one number")
}

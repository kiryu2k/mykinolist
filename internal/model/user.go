package model

import (
	"fmt"
	"regexp"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	validUsername = `^[\w]{6,50}$`
	validEmail    = `^[\w-\.]{6,30}@([\w-]{1,10}\.)[\w-]{2,4}$`
	// validPassword = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,30}$` TODO PASSWORD VALID REGEX
)

func (u *User) Validate() error {
	fmt.Println(*u)
	isMatched, err := regexp.MatchString(validUsername, u.Username)
	fmt.Println(isMatched, err)
	if err != nil {
		return err
	}
	if !isMatched {
		return fmt.Errorf("username must consist of letters and numbers, also it must contain from 6 to 50 characters")
	}
	isMatched, err = regexp.MatchString(validEmail, u.Email)
	if err != nil {
		return err
	}
	if !isMatched {
		return fmt.Errorf("email must consist of letters and numbers, also it mustn't exceed 50 characters")
	}
	// isMatched, err = regexp.MatchString(validPassword, u.Password)
	// if err != nil {
	// 	return err
	// }
	// if !isMatched {
	// 	return fmt.Errorf("password must contain from 8 to 30 characters, be at least one uppercase letter, one lowercase letter, one number and one special character")
	// }
	return nil
}

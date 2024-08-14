package userService

import (
	"FORUM/app/models"
	"FORUM/config/database"
	"unicode"
)

func CheckUsrExistUsername(username string) error {
	result := database.DB.Where("username= ?", username).First(&models.User{})
	return result.Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username= ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func ComparePwd(pwd1 string, pwd2 string) bool {
	return pwd1 == pwd2
}

func IsUsernameAllDigits(username string) bool {
	for _, char := range username {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func CheckPasswordLength(password string, minLength, maxLength int) bool {
	if len(password) >= minLength && len(password) <= maxLength {
		return true
	}
	return false
}

func Register(user models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

var validUserTypes = []int{1, 2}

func CheckUserType(userType int) bool {
	for _, validType := range validUserTypes {
		if userType == validType {
			return true
		}
	}
	return false
}

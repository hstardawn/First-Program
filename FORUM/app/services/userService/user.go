package userService

import (
	"FORUM/app/models"
	"FORUM/config/database"
	"unicode"

	"github.com/go-playground/validator/v10"
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

type UserPassword struct{
	Password string ` validate:"min=8,max=16"`
}
func CheckPasswordLength(password string) bool {
	validate := validator.New()

	var userpassword UserPassword
	userpassword.Password = password
	err := validate.Struct(userpassword)
	if err == nil {
		return true
	} else {
		return false
	}
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

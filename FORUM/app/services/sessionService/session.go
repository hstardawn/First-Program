package sessionService

import (
	"FORUM/app/models"
	"FORUM/app/services/userService"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetUserSession(c *gin.Context, user *models.User) error {
	webSession := sessions.Default(c)
	webSession.Options(sessions.Options{MaxAge: 3600 * 24 * 7, Path: "/api"})
	webSession.Set("userid", user.UserId)
	return webSession.Save()
}

func GetUserSession(c *gin.Context) (*models.User, error) {
    webSession := sessions.Default(c)
    userid := webSession.Get("userid")
    if userid == nil {
        return nil, errors.New("")
    }
    user, _ := userService.GetUserByUserid(userid.(uint))
    if user == nil {
        ClearUserSession(c)
        return nil, errors.New("")
    }
    return user, nil
}

func CheckUserSession(c *gin.Context) bool {
	webSession := sessions.Default(c)
	userid := webSession.Get("userid")
	return userid != nil
}

func ClearUserSession(c *gin.Context) {
	webSession := sessions.Default(c)
	webSession.Delete("userid")
	webSession.Save()
}
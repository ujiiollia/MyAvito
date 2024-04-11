package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type RoleGetter interface {
	GetRole(string, *gin.Context) (string, error)
}

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		str := ctx.Request.Header.Get("Authorization")
		if !strings.Contains(str, "Bearer") { //todo: написать в README
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("token miossing Bearer"))
			return
		}
		token := strings.Split(str, "Bearer ")
		t, _, err := new(jwt.Parser).ParseUnverified(token[1], jwt.MapClaims{})
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("token parse error"))
			return
		}
		var userId string
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			uu := claims["userId"]
			userId = uu.(string)
		}
		if userId == "" {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("token userId is null"))
			return
		}
		ctx.Set("userId", userId)
	}
}

const (
	AdminRole = "admin"
	UserRole  = "user"
	//GuestRole = "guest"
)

func CheckRole(isAdmin bool, rg RoleGetter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := ctx.Get("userId")
		if !ok {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("role not found"))
			return
		}

		role, err := rg.GetRole(userId.(string), ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusForbidden, errors.New("no user in data base"))
			return
		}
		if isAdmin {
			if role != AdminRole {
				ctx.AbortWithError(http.StatusForbidden, errors.New("token userId is null"))
				return
			}
		}
	}
}

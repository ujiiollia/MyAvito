package auth

import (
	"app/internal/storage/postgresql"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleGetter interface {
	GetRole(string, *gin.Context) (string, error)
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		str := r.Header.Get("Authorization")
		if !strings.Contains(str, "Bearer") {
			http.Error(w, "Token missing Bearer", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(str, "Bearer ")[1]
		t, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			http.Error(w, "Token parse error", http.StatusUnauthorized)
			return
		}

		var userId string
		if claims, ok := t.Claims.(jwt.MapClaims); ok && claims["userId"] != nil {
			userId, ok = claims["userId"].(string)
			if !ok {
				http.Error(w, "Token userId is not a string", http.StatusUnauthorized)
				return
			}
		}

		if userId == "" {
			http.Error(w, "Token userId is null", http.StatusUnauthorized)
			return
		}

		// Использование контекста для прокидывания userId в дальнейшие обработчики
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const (
	AdminRole = "admin"
	UserRole  = "user"
	//GuestRole = "guest"
	AdminRightsRequired = true
	UserRightsRequired  = false
)

func CheckRole(isAdmin bool, pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userId, ok := r.Context().Value("userId").(string)
			if !ok {
				http.Error(w, "role not found", http.StatusUnauthorized)
				return
			}
			role, err := postgresql.GetRole(userId, pool)
			if err != nil {
				http.Error(w, "no user in data base", http.StatusForbidden)
				return
			}
			if isAdmin {
				if role != AdminRole {
					http.Error(w, "token userId is null", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

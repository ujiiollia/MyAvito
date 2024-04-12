package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	takeNewestData         = "true"
	takeFiveMinutesOldData = "false"
)

func GetUserBanner(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		useLastRevision := ctx.Query("use_last_revision")

		switch useLastRevision {
		case takeNewestData:
			rows, err := pool.Query(context.Background(), "SELECT content FROM user_banner WHERE is_active = TRUE")
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer rows.Close()

			var banner []string
			for rows.Next() {
				var content string
				err := rows.Scan(&content)
				if err != nil {
					ctx.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				banner = append(banner, content)
			}

			if len(banner) == 0 {
				ctx.AbortWithError(http.StatusNotFound, errors.New("no active user_banner found"))
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"user_banner": banner})
			return
		case takeFiveMinutesOldData:
			//todo read cache
			return
		default:
			ctx.AbortWithError(http.StatusNotFound, errors.New("use_last_revision valuse is out of range"))
			return
		}
	}
}

package middleware

import (
	"app/internal/middleware/auth"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	takeNewestData         = "true"
	takeFiveMinutesOldData = "false"
	adminRightsRequired    = true
)

type UserBanner struct {
	ID        int64           `json:"id"`
	TagID     int64           `json:"tag_id"`
	FeatureID int64           `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

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

func GetBanner(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//check admin rights
		auth.Authorization()
		auth.CheckRole(adminRightsRequired, pool)

		tagID, ok := ctx.GetQuery("tag_id")
		if !ok {
			ctx.AbortWithError(http.StatusNotFound, errors.New("tagID not found in Query"))
			return
		}
		featureID, ok := ctx.GetQuery("feature_id")
		if !ok {
			ctx.AbortWithError(http.StatusNotFound, errors.New("featureID not found in Query"))
			return
		}
		sql := `SELECT id, tag_id, feature_id, content, is_active, created_at, updated_at
				FROM user_banner
				WHERE tag_id=$1 AND feature_id=$2 AND is_active=true`

		rows, err := pool.Query(context.Background(), sql, tagID, featureID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var banners []UserBanner
		for rows.Next() {
			var b UserBanner
			err = rows.Scan(&b.ID, &b.TagID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			banners = append(banners, b)
		}

		ctx.JSON(http.StatusOK, banners)
	}
}

package middleware

import (
	checkdigits "app/internal/lib/checkDigitsInStr"
	mapcache "app/internal/storage/cache"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	takeNewestData         = "true"
	takeFiveMinutesOldData = "false"
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

func GetUserBanner(pool *pgxpool.Pool, cace *mapcache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		useLastRevision := r.URL.Query().Get("use_last_revision")
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		switch useLastRevision {

		case takeFiveMinutesOldData:
			// read cache
			if banner, isExist := cace.Get(tagID + ":" + featureID); isExist {
				json.NewEncoder(w).Encode(banner)
				json.NewEncoder(w).Encode(banner)
				return
			}
			fallthrough

		case takeNewestData:

			query := fmt.Sprintf("SELECT * FROM user_banner WHERE tag_id = %s AND feature_id = %s AND is_active = %t", tagID, featureID, true)

			rows, err := pool.Query(r.Context(), query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var banners []UserBanner
			for rows.Next() {
				var banner UserBanner
				err := rows.Scan(&banner.ID, &banner.TagID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				banners = append(banners, banner)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(banners)
			return

		}
	}
}

func GetAllBannerByFeatureAndTag(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		featureID := r.URL.Query().Get("feature_id")
		tagID := r.URL.Query().Get("tag_id")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		// проверка на то что в limit только цифры ("123" - true, "1dva3" - false)
		if ok, err := checkdigits.CheckInSring(limit); !ok || err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// проверка на то что в offset только цифры ("123" - true, "1dva3" - false)
		if ok, err := checkdigits.CheckInSring(offset); !ok || err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var banners []UserBanner
		rows, err := pool.Query(context.Background(), "SELECT * FROM user_banner WHERE feature_id=$1 AND tag_id=$2 LIMIT $3 OFFSET $4",
			featureID, tagID, limit, offset)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var banner UserBanner
			err = rows.Scan(&banner.ID, &banner.TagID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			banners = append(banners, banner)
		}

		jsonBanners, err := json.Marshal(banners)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBanners)
	}

}

func CreateBanner(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userBanner UserBanner
		if err := json.NewDecoder(r.Body).Decode(&userBanner); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createTableSQL := `
		            CREATE TABLE IF NOT EXISTS user_banner (
		                id SERIAL PRIMARY KEY,
		                tag_id INTEGER NOT NULL,
		                feature_id INTEGER NOT NULL,
		                content JSON NOT NULL,
		                is_active BOOLEAN DEFAULT TRUE,
		                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		            );`

		_, err := pool.Exec(context.Background(), createTableSQL)
		if err != nil {
			http.Error(w, "Failed to insert user banner", http.StatusInternalServerError)
			return
		}

		_, err = pool.Exec(context.Background(),
			`INSERT INTO user_banner (tag_id, feature_id, content, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
			userBanner.TagID, userBanner.FeatureID, userBanner.Content, userBanner.IsActive, time.Now().String(), time.Now().String())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user_banner created successfully"))
	}
}

func PatchBanner(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		content := r.URL.Query().Get("content")
		isActive := r.URL.Query().Get("is_active")

		updateSQL := `UPDATE user_banner SET tag_id=$1, feature_id=$2, content=$3, is_active=$4, updated_at=$5 WHERE id=$6`
		_, err := pool.Exec(context.Background(), updateSQL, tagID, featureID, content, isActive, time.Now(), id)
		if err != nil {
			http.Error(w, "Failed to update user_banner", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user_banner updated successfully"))
	}
}

func DeleteBanner(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем id из URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Вызываем функцию для удаления UserBanner из базы данных
		// err = removeUserBannerByID(r.Context(), pool, id)
		_, err = pool.Exec(r.Context(), `DELETE FROM user_banners WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Failed to delete UserBanner", http.StatusInternalServerError)
			return
		}

		// Отправляем подтверждение об успешном удалении
		w.WriteHeader(http.StatusNoContent)
	}
}

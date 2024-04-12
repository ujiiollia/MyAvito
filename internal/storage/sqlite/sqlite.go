package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Banner представляет структуру баннера
type Banner struct {
	ID        int
	TagIDs    []int
	FeatureID int
	Content   string
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}

type User struct {
	ID        int
	Name      string
	RoleToken AdminToken
	Email     string
	TagIDs    []int
	FeatureID int
}

type AdminToken struct {
	JWTAdmin string
}

func NewBanner(storagePath string) (*Storage, error) {
	const el = `sqlite.NewBanner`
	db, err := sql.Open(`sqlite3`, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)

	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS banners (
		id INTEGER PRIMARY KEY ,
		tag_ids INTEGER NOT NULL,
		feature_id INTEGER NOT NULL,
		content TEXT,
		active BOOLEAN,
		last_updated_time DATETIME
	);`) //todo: CREATE INDEX IF NOT EXISTS ...

	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	return &Storage{db: db}, nil
}

func NewUser(storagePath string) (*Storage, error) {
	const el = `sqlite.NewUser`
	db, err := sql.Open(`sqlite3`, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)

	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS admin_token (
		id INTEGER PRIMARY KEY
		);`)
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT VARCHAR (100) NOT NULL,gti
		role_token INTEGER,
		email TEXT VARCHAR (100) NOT NULL,
		tag_ids TEXT,
		feature_id INTEGER,
		FOREIGN KEY (role_token) REFERENCES admin_token(id)
	   );`) //todo: CREATE INDEX IF NOT EXISTS ...

	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)

	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	return &Storage{db: db}, nil
}

// Поиск баннера по фиче и тегу
func (s *Storage) GetBannerByFeatureAndTag(featureID int, tagID int, lastRevisionFlag string) (*Banner, error) {
	const el = "sqlite.GetBannerByFeatureAndTag"
	var banner Banner
	var err error

	if lastRevisionFlag == "use_last_revision" {
		err = s.db.QueryRow(
			`SELECT id, feature_id, tag_ids, json_data, active, last_updated_time 
			FROM banners WHERE feature_id=? AND ? IN (SELECT tag_ids FROM banner_tags WHERE banner_id=banners.id)`,
			featureID, tagID).Scan(&banner.ID, &banner.FeatureID, &banner.TagIDs, &banner.Content, &banner.IsActive, &banner.UpdatedAt)

	} else {
		err = s.db.QueryRow(
			`SELECT id, feature_id, tag_ids, json_data, active, last_updated_time 
		FROM banners WHERE feature_id=? AND ? IN (SELECT tag_ids FROM banner_tags WHERE banner_id=banners.id)`,
			featureID, tagID).Scan(&banner.ID, &banner.FeatureID, &banner.TagIDs, &banner.Content, &banner.IsActive, &banner.UpdatedAt)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	return &banner, nil
}

// // Получение баннера по идентификатору
// func (s *Storage) GetBannerByID(bannerID int) (*Banner, error) {
// 	const el = "sqlite.getBannerByID"
// 	var banner Banner
// 	err := s.db.QueryRow("SELECT * FROM banners WHERE banner_id = ?",
// 		bannerID).Scan(&banner.ID, &banner.FeatureID, &banner.TagIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", el, err)
// 	}
// 	return &banner, nil
// }

func (s *Storage) AddBanner(banner Banner) error {
	_, err := s.db.Exec(`INSERT INTO banners (feature_id, tag_ids, json_data, active, last_updated_time) VALUES (?, ?, ?, ?, ?)`,
		banner.FeatureID, pq.Array(banner.TagIDs), banner.Content, banner.IsActive, banner.UpdatedAt)
	return err
}

func (s *Storage) UpdateBanner(banner Banner) error {
	_, err := s.db.Exec(`UPDATE banners SET feature_id=?, tag_ids=?, json_data=?, active=?, last_updated_time=? WHERE id=?`,
		banner.FeatureID, pq.Array(banner.TagIDs), banner.Content, banner.IsActive, banner.UpdatedAt, banner.ID)
	return err
}

func (s *Storage) DeleteBanner(id int) error {
	_, err := s.db.Exec(`DELETE FROM banners WHERE id=?`, id)
	return err
}

func (s *Storage) GetAllBanners() ([]Banner, error) {
	rows, err := s.db.Query("SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at FROM banners")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []Banner
	for rows.Next() {
		var banner Banner
		err := rows.Scan(&banner.ID, &banner.TagIDs, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

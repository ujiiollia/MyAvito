package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Banner представляет структуру баннера
type Banner struct {
	ID              int
	FeatureID       int
	TagIDs          []int
	JSONData        string
	Active          int
	LastUpdatedTime string
}

// струтктура тега
type Tag struct {
	ID int
}

// структура фичи
type Feature struct {
	ID int
}

type User struct {
	ID        int
	Name      string
	RoleToken string
	Email     string
	TagIDs    []int
	FeatureID int
}

func New(storagePath string) (*Storage, error) {
	const el = `sqlite.New`
	db, err := sql.Open(`sqlite3`, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)

	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS banners (
		id INTEGER PRIMARY KEY,
		feature_id INTEGER,
		tag_id INTEGER,
		json_data TEXT,
		active INTEGER,
		last_updated_time DATETIME
	);`) //todo: CREATE INDEX IF NOT EXISTS ...

	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)

	}

	// Создание таблицы tags для хранения информации о тегах
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY
   )`)
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы features для хранения информации о фичах
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS features (
	id INTEGER PRIMARY KEY
   )`)
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы admin_tokens для хранения админских токенов
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS admin_tokens (
	id INTEGER PRIMARY KEY,
	token TEXT
   )`)
	if err != nil {
		log.Fatal(err)
	}
	// Создание таблицы user_tokens для хранения пользовательских токенов
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS user_tokens (
	id INTEGER PRIMARY KEY,
	token TEXT,
	user_id INTEGER
   )`)
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
			`SELECT id, feature_id, tags, json_data, active, last_updated_time 
			FROM banners WHERE feature_id=? AND ? IN (SELECT tag_id FROM banner_tags WHERE banner_id=banners.id)`,
			featureID, tagID).Scan(&banner.ID, &banner.FeatureID, &banner.TagIDs, &banner.JSONData, &banner.Active, &banner.LastUpdatedTime)

	} else {
		err = s.db.QueryRow(
			`SELECT id, feature_id, tags, json_data, active, last_updated_time 
		FROM banners WHERE feature_id=? AND ? IN (SELECT tag_id FROM banner_tags WHERE banner_id=banners.id)`,
			featureID, tagID).Scan(&banner.ID, &banner.FeatureID, &banner.TagIDs, &banner.JSONData, &banner.Active, &banner.LastUpdatedTime)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	return &banner, nil
}

// // Получения баннера по идентификатору
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
	_, err := s.db.Exec(`INSERT INTO banners (feature_id, tags, json_data, active, last_updated_time) VALUES (?, ?, ?, ?, ?)`,
		banner.FeatureID, banner.TagIDs, banner.JSONData, banner.Active, banner.LastUpdatedTime)
	return err
}

func (s *Storage) UpdateBanner(banner Banner) error {
	_, err := s.db.Exec(`UPDATE banners SET feature_id=?, tags=?, json_data=?, active=?, last_updated_time=? WHERE id=?`,
		banner.FeatureID, banner.TagIDs, banner.JSONData, banner.Active, banner.LastUpdatedTime, banner.ID)
	return err
}

func (s *Storage) DeleteBanner(id int) error {
	_, err := s.db.Exec(`DELETE FROM banners WHERE id=?`, id)
	return err
}

func (s *Storage) AddTag(tag Tag) error {
	_, err := s.db.Exec(`INSERT INTO tags (id) VALUES (?)`, tag.ID)
	return err
}

func (s *Storage) UpdateTag(tag Tag) error {
	_, err := s.db.Exec(`UPDATE tags SET id=? WHERE id=?`, tag.ID, tag.ID)
	return err
}

func (s *Storage) DeleteTag(id int) error {
	_, err := s.db.Exec(`DELETE FROM tags WHERE id=?`, id)
	return err
}

func (s *Storage) AddFeature(ftr Feature) error {
	_, err := s.db.Exec(`INSERT INTO features (id) VALUES (?)`, ftr.ID)
	return err
}

func (s *Storage) UpdateFeature(ftr Feature) error {
	_, err := s.db.Exec(`UPDATE features SET id=? WHERE id=?`, ftr.ID, ftr.ID)
	return err
}

func (s *Storage) DeleteFeature(id int) error {
	_, err := s.db.Exec(`DELETE FROM features WHERE id=?`, id)
	return err
}

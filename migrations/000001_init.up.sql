BEGIN;

CREATE TABLE IF NOT EXISTS user_banner (
		id INTEGER PRIMARY KEY ,
		tag_id INTEGER NOT NULL,
		feature_id INTEGER NOT NULL,
		content JSON NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY ,
		is_admin VARCHAR(5) NOT NULL
	);

COMMIT;
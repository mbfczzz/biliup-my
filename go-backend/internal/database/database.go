package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Init(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	log.Println("数据库初始化完成 (MySQL)")
	return db, nil
}

func migrate(db *sqlx.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS live_streamers (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(512) NOT NULL UNIQUE,
			remark VARCHAR(255) DEFAULT '',
			filename_prefix VARCHAR(255) DEFAULT '{streamer}_%Y-%m-%d_%H_%M_%S',
			time_range VARCHAR(64) DEFAULT '',
			upload_template_id BIGINT DEFAULT 0,
			format VARCHAR(16) DEFAULT 'flv',
			segment_time INT DEFAULT 3600,
			file_size BIGINT DEFAULT 0,
			opt_args TEXT,
			excluded_keywords TEXT,
			paused TINYINT DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS upload_templates (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			template_name VARCHAR(255) NOT NULL,
			title VARCHAR(512) DEFAULT '{streamer} {title} %Y-%m-%d',
			tid INT DEFAULT 21,
			copyright TINYINT DEFAULT 2,
			copyright_source VARCHAR(512) DEFAULT '',
			cover_path VARCHAR(512) DEFAULT '',
			description TEXT,
			dynamic VARCHAR(512) DEFAULT '',
			dtime BIGINT DEFAULT 0,
			dolby TINYINT DEFAULT 0,
			hires TINYINT DEFAULT 0,
			no_reprint TINYINT DEFAULT 0,
			tags VARCHAR(1024) DEFAULT '[]',
			user_cookie TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS streamer_info (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) DEFAULT '',
			url VARCHAR(512) DEFAULT '',
			title VARCHAR(512) DEFAULT '',
			live_cover_path VARCHAR(512) DEFAULT '',
			status VARCHAR(32) DEFAULT 'idle',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS file_list (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			file_path VARCHAR(1024) NOT NULL,
			file_size BIGINT DEFAULT 0,
			streamer_info_id BIGINT,
			uploaded TINYINT DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (streamer_info_id) REFERENCES streamer_info(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS configuration (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			` + "`key`" + ` VARCHAR(128) NOT NULL UNIQUE,
			value TEXT
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(128) NOT NULL UNIQUE,
			cookie_data TEXT,
			uid BIGINT DEFAULT 0,
			nickname VARCHAR(128) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}

	for _, sql := range tables {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}

	// 插入默认配置
	defaults := map[string]string{
		"download_dir":        "./downloads",
		"segment_time":        "3600",
		"file_size":           "0",
		"filename_prefix":     "{streamer}_%Y-%m-%d_%H_%M_%S",
		"check_interval":      "60",
		"upload_line":         "AUTO",
		"upload_threads":      "3",
		"downloader":          "ffmpeg",
		"filtering_threshold": "10",
		"delay":               "0",
	}
	for k, v := range defaults {
		db.Exec("INSERT IGNORE INTO configuration (`key`, value) VALUES (?, ?)", k, v)
	}

	return nil
}

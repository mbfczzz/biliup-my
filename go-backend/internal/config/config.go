package config

import (
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
)

// Config 全局配置管理
type Config struct {
	mu   sync.RWMutex
	data map[string]string
	db   *sqlx.DB
}

func Load(db *sqlx.DB) (*Config, error) {
	cfg := &Config{
		data: make(map[string]string),
		db:   db,
	}

	rows, err := db.Queryx("SELECT `key`, value FROM configuration")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		cfg.data[key] = value
	}

	return cfg, nil
}

func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Config) GetInt(key string, defaultVal int) int {
	v := c.Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

func (c *Config) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	_, err := c.db.Exec(
		"INSERT INTO configuration (`key`, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = VALUES(value)",
		key, value,
	)
	return err
}

func (c *Config) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]string, len(c.data))
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

func (c *Config) SetAll(data map[string]string) error {
	for k, v := range data {
		if err := c.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

package db

import (
	"fmt"
	"inivoice/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(cfg *config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database))

	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Duration(cfg.Pool.MaxIdleConns))
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Pool.MaxOpenConns)

	return db, nil
}

func Release(db *sqlx.DB) error {
	return db.Close()
}

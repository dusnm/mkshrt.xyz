package database

import (
	"database/sql"
	"fmt"
	"github.com/dusnm/mkshrt.xyz/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func New(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return &sql.DB{}, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}

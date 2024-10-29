package database

import (
	"log"

	"github.com/Aritiaya50217/E-CommerceRESTAPIs/config"
	// _ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect by postgres
	// db, err := sqlx.Connect("pgx", cfg.Url())

	// Connect by mysql
	db, err := sqlx.Connect("mysql", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}

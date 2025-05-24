package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"idm/inner/common"
	"time"
)

// Временная переменная, которая будет ссылаться на подключение к базе данных
var DB *sqlx.DB

// ConnectDb получить конфиг и подключиться с ним к базе данных
func ConnectDb() *sqlx.DB {
	cfg := common.GetConfig(".env")
	return ConnectDbWithCfg(cfg)
}

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом
func ConnectDbWithCfg(cfg common.Config) *sqlx.DB {
	DB = sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(20)
	DB.SetConnMaxLifetime(1 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	return DB
}

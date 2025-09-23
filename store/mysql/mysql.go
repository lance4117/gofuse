package store

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/once"
	"xorm.io/xorm"
)

type Config struct {
	Driver          string // 数据库驱动名称 (mysql, postgres, sqlite3 ...)
	Source          string // Data Source Name (连接字符串)
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ShowSQL         bool
}

type DataBase struct {
	engine *xorm.Engine
}

var NewMySql = once.DoWithParam(func(cfg Config) *DataBase {
	engine, err := xorm.NewEngine(cfg.Driver, cfg.Source)
	if err != nil {
		logger.Fatal(err, "Init My Sql Fail")
		return nil
	}
	return &DataBase{engine}
})

func (db *DataBase) Insert(beans ...interface{}) (int64, error) {
	return db.engine.Insert(beans)
}

func (db *DataBase) Get(beans ...interface{}) (bool, error) {
	return db.engine.Get(beans)
}

func (db *DataBase) Ping() error {
	return db.engine.Ping()
}

func (db *DataBase) Close() error {
	return db.engine.Close()
}

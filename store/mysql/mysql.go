package store

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lance4117/gofuse/errs"
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
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
		return nil
	}
	engine.SetMaxOpenConns(cfg.MaxOpenConns)
	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	engine.ShowSQL(cfg.ShowSQL)

	err = engine.Ping()
	if err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
		return nil
	}
	return &DataBase{engine}
})

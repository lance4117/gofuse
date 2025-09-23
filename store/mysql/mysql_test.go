package store

import "testing"

func TestMysql(t *testing.T) {
	cfg := Config{
		Driver:          "MySql",
		Source:          "root:123123@/test?charset=utf8",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 0,
		ShowSQL:         true,
	}
	sql := NewMySql(cfg)
	table := sql.engine.Table("cities")
	err := table.AllCols().Ping()

	t.Log(err)
}

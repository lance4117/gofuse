package test

import (
	"os"
	"testing"

	"github.com/lance4117/gofuse/store/dbs"
)

func mysqlCfg(t *testing.T) dbs.Config {
	t.Helper()
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set TEST_MYSQL_DSN to run mysql integration tests")
	}
	driver := os.Getenv("TEST_MYSQL_DRIVER")
	if driver == "" {
		driver = "mysql"
	}
	return dbs.Config{
		Name:            "test-mysql",
		Driver:          driver,
		DSN:             dsn,
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 0,
		ShowSQL:         false,
	}
}

func TestMysql(t *testing.T) {
	cfg := mysqlCfg(t)

	repo, err := dbs.NewRepo[User](cfg)
	if err != nil {
		t.Skip("mysql not available:", err)
	}

	user, _, err := repo.GetByID(1939940059276906496)
	if err != nil {
		t.Skip("mysql not available for query:", err)
	}
	user.Name = "updated"
	err = repo.Delete(user)
	t.Log(err)
}

func TestSession(t *testing.T) {
	cfg := mysqlCfg(t)
	repo, err := dbs.NewRepo[User](cfg)
	if err != nil {
		t.Skip("mysql not available:", err)
	}

	err = repo.DoTx(func(txRepo *dbs.Repo[User]) error {
		if err := txRepo.Insert(&User{Name: "Bob"}); err != nil {
			return err
		}
		if err := txRepo.UpdateById(&User{Name: "31"}, 1); err != nil {
			return err
		}
		return nil
	})

	t.Log(err)
}

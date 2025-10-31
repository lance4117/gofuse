package test

import (
	"testing"

	"github.com/lance4117/gofuse/store/dbs"
)

func TestMysql(t *testing.T) {
	cfg := dbs.Config{
		Name:            "cfg1",
		Driver:          "mysql",
		DSN:             "admin:@2CSacf378*`@/node1?charset=utf8",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 0,
		ShowSQL:         true,
	}

	repo := dbs.NewRepo[User](cfg)

	user, _, err := repo.GetByID(1939940059276906496)
	if err != nil {
		t.Fatal(err)
	}
	user.Name = "updated"
	err = repo.Delete(user)
	t.Log(err)
}

func TestSession(t *testing.T) {
	cfg := dbs.Config{
		Name:            "cfg1",
		Driver:          "mysql",
		DSN:             "admin:@2CSacf378*`@/node1?charset=utf8",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 0,
		ShowSQL:         true,
	}
	repo := dbs.NewRepo[User](cfg)
	err := repo.DoTx(func(txRepo *dbs.Repo[User]) error {
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

package db_test

import (
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func Test1(t *testing.T) {
	sql, err := db.NewSqliteDriver("./test.sqlite")
	if err != nil {
		panic(err)
	}
	db := db.NewDB(sql, "./test.sqlite")

	db.Close()
}

package db_test

import (
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func TestInitDB(t *testing.T) {
	db.NewSqlite("../../test.sqlite")

}

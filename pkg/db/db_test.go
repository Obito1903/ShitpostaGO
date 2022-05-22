package db_test

import (
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func TestInitDB(t *testing.T) {
	db.NewSqlite("../../test.sqlite", db.DatabaseConfig{
		Folder: "./test/",
	})

}

func TestVideoTranscode(t *testing.T) {
	db.ImportFile("../../Bucket.webm", "../../", "video0")
}

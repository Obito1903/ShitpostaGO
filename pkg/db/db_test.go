package db_test

import (
	"path"
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func TestInitDB(t *testing.T) {
	db.NewSqlite(db.Database{
		Folder: "../../appTest/",
	})

}

func TestMediaImport(t *testing.T) {
	shitdb := db.NewSqlite(db.Database{
		Folder: "../../appTest/",
	})
	media, err := shitdb.NewMediaFromPath(path.Join(shitdb.Folder, "/import/1.mp4"))
	t.Log(media)

	if err != nil {
		t.Error("DB ERROR : ", err)
	}
}

func TestCategories(t *testing.T) {
	shitdb := db.NewSqlite(db.Database{
		Folder: "../../appTest/",
	})
	// cat := shitdb.NewCategory("test")
	// t.Log(cat)
	media := shitdb.GetMediaFromId(1)

	// media = shitdb.AddCategoryToMedia(media, cat)

	// media = shitdb.GetMediaFromId(1)
	t.Log(media)
	if media.Catergories == nil {
		t.Error("No cats")
	}
}

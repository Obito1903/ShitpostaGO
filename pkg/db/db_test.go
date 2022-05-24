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
	cat1 := shitdb.NewCategory("test1")
	cat2 := shitdb.NewCategory("test2")
	cat3 := shitdb.NewCategory("test3")

	media1, err := shitdb.NewMediaFromPath(path.Join(shitdb.Folder, "/import/1.mp4"))
	if err != nil {
		t.Error("DB ERROR : ", err)
	}
	media2, err := shitdb.NewMediaFromPath(path.Join(shitdb.Folder, "/import/2.mp4"))
	if err != nil {
		t.Error("DB ERROR : ", err)
	}
	media1 = shitdb.AddCategoryToMedia(media1, cat1)
	media1 = shitdb.AddCategoryToMedia(media1, cat2)
	shitdb.AddCategoryToMedia(media1, cat3)

	media2 = shitdb.AddCategoryToMedia(media2, cat1)
	shitdb.AddCategoryToMedia(media2, cat2)

	medias := shitdb.GetMediasFromCats([]db.Category{cat1, cat2}, 10, 0)

	t.Log(medias)
	if media1.Catergories == nil {
		t.Error("No cats")
	}
}

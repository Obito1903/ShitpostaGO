package db_test

import (
	"fmt"
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func Test1(t *testing.T) {
	sql, err := db.NewSqliteDriver("../../tmp/test.sqlite")
	if err != nil {
		panic(err)
	}
	fsdb := db.NewFsMediaDB("../../tmp")
	db_inst := db.NewDB(sql, fsdb, "../../tmp")

	fmt.Println("Creating media")
	media, err := db_inst.NewMedia(db.Metadata{
		OgName:    "test",
		Name:      "test",
		MediaType: db.MediaType_Image,
		FileType:  "jpg",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", media)

	fmt.Println("Getting media by id")
	media, err = db_inst.GetMedia(media.Id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", media)

	fmt.Println("Getting media by name")
	media, err = db_inst.GetMediaByName("test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", media)

	fmt.Println("Getting random media")
	media, err = db_inst.GetMediaRandom(db.MediaType_Image)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", media)

	fmt.Println("Removing media")
	err = db_inst.DeleteMedia(media.Id)
	if err != nil {
		panic(err)
	}

	db_inst.Close()
}

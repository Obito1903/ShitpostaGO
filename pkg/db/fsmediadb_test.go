package db_test

import (
	"testing"

	"github.com/Obito1903/shitpostaGo/pkg/db"
)

func TestConvertImage(t *testing.T) {

	err := db.ConvertImageMedia("../../test/test.png", "../../test/test2.jpg")
	if err != nil {
		panic(err)
	}

}

func TestConvertVideo(t *testing.T) {

	err := db.ConvertVideoMedia("../../test/testv.mov", "../../test/test2.mp4")
	if err != nil {
		panic(err)
	}

}

func TestAddMedia(t *testing.T) {
	fsdb := db.FsMediaDB{Path: "../../tmp"}
	fsdb.AddMedia()
}

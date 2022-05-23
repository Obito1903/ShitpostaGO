package db

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"time"
)

type MediaType string

const (
	Video   MediaType = "video"
	Image   MediaType = "image"
	Unknown MediaType = "unknown"
)

type Database struct {
	DatabaseInterface
	Folder string
}

type DatabaseInterface interface {
	GetConfig() Database

	NewMediaFromPath(path string) (Media, error)
	UpdateMedia(media Media) Media
	RemoveMedia(media Media)
	GetMediaFromId(id int64) Media
	GetMediasFromCats(categories []Category) []Media
	AddCategoryToMedia(media Media, category Category) Media
	RemoveCategoryFromMedia(media Media, category Category) Media

	NewCategory(name string) Category
	RemoveCategory(category Category)
	UpdateCategory(category Category) Category
	GetCategoryFromId(id int64) Category
	GetCategoriesFromId(ids []int64) []Category
	GetCategoriesFromMedia(media Media) []Category
	GetCategories() []Category
}

type Category struct {
	Id   int
	Name string
}

type categoryLink struct {
	Media_id    int64
	Category_id int64
}

type Media struct {
	Id          int64
	Og_name     string
	Name        string
	Path        string
	Date        time.Time
	Type_       MediaType
	Catergories []Category
}

func (media Media) String() string {
	out, _ := json.Marshal(media)
	return string(out)
}

func checkDBerr(err error) {
	if err != nil {
		log.Fatal("DB ERROR : ", err)
	}
}

func (db Database) initFs() {
	for _, folder := range [4]string{
		db.Folder,
		path.Join(db.Folder, "/videos/"),
		path.Join(db.Folder, "/images/"),
		path.Join(db.Folder, "/import/"),
	} {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, os.ModePerm)
		}
	}
}

func CheckFileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

package db

import (
	"log"
	"time"
)

type MediaType string

const (
	Video   MediaType = "video"
	Image   MediaType = "image"
	Unknown MediaType = "unknown"
)

type DatabaseConfig struct {
	Folder string
}

type Database interface {
	GetConfig() DatabaseConfig

	NewMediaFromPath(path string) Media

	UpdateMedia(media Media) Media
	NewCategory(name string) Category
	UpdateCategory(category Category) Category

	GetMediaFromId(id int64) Media
	GetMediaCats(media Media) []Category
	GetMediasFromCats(categories []Category) []Media

	GetCategoryFromId(id int64) Category
	GetCategories() []Category
}

type Category struct {
	id   int
	Name string
}

type Media struct {
	id          int
	Og_name     string
	Name        string
	Path        string
	Date        time.Time
	Type_       MediaType
	Catergories []Category
}

func checkDBerr(err error) {
	if err != nil {
		log.Fatal("DB ERROR : ", err)
	}
}

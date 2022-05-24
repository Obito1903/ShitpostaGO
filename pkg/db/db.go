package db

import (
	"encoding/json"
	"errors"
	"fmt"
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

type ErrorCode string

const (
	InitFailed    ErrorCode = "InitFailed"
	FileNotFound  ErrorCode = "FileNotFound"
	ImportFailed  ErrorCode = "ImportFailed"
	EntryNotFound ErrorCode = "EntryNotFound"
)

type DBError struct {
	Err     error
	Message string
	Code    ErrorCode
}

func (dberr DBError) Error() string {
	if dberr.Err != nil {
		return fmt.Sprintf("%s: %s\n  ⎣%s", dberr.Code, dberr.Message, dberr.Err)
	} else {
		return fmt.Sprintf("%s: %s", dberr.Code, dberr.Message)
	}
}

func checkErr(err error, message string, code ErrorCode) *DBError {
	if (err != nil) && (err != (*DBError)(nil)) {
		log.Println(err)
		return &DBError{
			err,
			message,
			code,
		}
	}
	return nil
}

type Database struct {
	DatabaseInterface
	Folder string
}

type DatabaseInterface interface {
	GetConfig() Database

	NewMediaFromPath(filePath string) (Media, *DBError)
	UpdateMedia(media Media) (Media, *DBError)
	RemoveMedia(media Media) *DBError
	GetMediaFromId(id int64) (Media, *DBError)
	GetMediasFromCats(categories []Category, limit int, offsetn int) ([]Media, *DBError)
	AddCategoryToMedia(media Media, category Category) (Media, *DBError)
	RemoveCategoryFromMedia(media Media, category Category) (Media, *DBError)

	GetRandomMedia(type_ MediaType) (Media, *DBError)

	NewCategory(name string) (Category, *DBError)
	UpdateCategory(category Category) (Category, *DBError)
	RemoveCategory(category Category) *DBError
	GetCategoryFromId(id int64) (Category, *DBError)
	GetCategoriesFromIds(ids []int64) ([]Category, *DBError)
	GetCategoriesFromMedia(media Media) ([]Category, *DBError)
	GetCategories() ([]Category, *DBError)
}

type Category struct {
	Id   int
	Name string
}

// type categoryLink struct {
// 	Media_id    int64
// 	Category_id int64
// }

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

func CheckFileExist(filePath string) (bool, *DBError) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, &DBError{err, fmt.Sprintf("FS error on : %s", filePath), FileNotFound}
}

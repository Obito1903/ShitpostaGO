package db

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(zerolog.DebugLevel).Output(zerolog.ConsoleWriter{Out: os.Stdout})
)

type MediaType int

const (
	MediaType_Unknown MediaType = iota
	MediaType_Image
	MediaType_Video
	MediaType_Audio
)

type PermissionLevel int

const (
	PermissionLevel_Unknown PermissionLevel = iota
	PermissionLevel_User
	PermissionLevel_Admin
	PermissionLevel_SuperAdmin
)

type Metadata struct {
	Id          int       `json:"id"`
	OgName      string    `json:"og_name"`
	Name        string    `json:"name"`
	MediaType   MediaType `json:"media_type"`
	FileType    string    `json:"file_type"`
	DateAdded   time.Time `json:"date_added"`
	DateCreated time.Time `json:"date_created"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	HashPass   string `json:"hash_pass"`
	Permission int    `json:"permission"`
}

type metaDB interface {
	// Medias
	// Get media by id
	GetMedia(id int) (Metadata, error)
	// Get media by name
	GetMediaByName(name string) (Metadata, error)
	// Get random media of spcified type
	GetMediaRandom(mediatype MediaType) (Metadata, error)
	// Get Redom media of specified type and category
	GetMediaRandomByCategory(mediatype MediaType, categoryid int) (Metadata, error)
	// Get all media faved by user
	GetMediasFromUser(user User) ([]Metadata, error)
	// Get all media of type
	GetMediasByType(mediaType MediaType) ([]Metadata, error)
	// Get all media of type and category
	GetMediasByTypeAndCategory(mediaType MediaType, categoryid int) ([]Metadata, error)

	// Search media by name
	SearchMediaByName(name string) ([]Metadata, error)

	// Create a new media, retuns metadata with new id
	NewMedia(metadata Metadata) (Metadata, error)
	// Update media metadata, only update the name
	UpdateMedia(metadata Metadata) error
	// Delete media
	DeleteMedia(id int) error

	// Categories
	// Get category by id
	GetCategory(id int) (Category, error)
	// Get category by name
	GetCategoryByName(name string) (Category, error)
	// Get all categories of a media
	GetCategoriesFromMedia(mediaid int) ([]Category, error)
	// Get all categories
	GetCategories() ([]Category, error)

	// Create a new category, retuns category with new id
	NewCategory(category Category) (Category, error)
	// Update category, only update the name
	UpdateCategory(category Category) error
	// Delete category
	DeleteCategory(id int) error

	// Users
	// Get user by id
	GetUser(id int) (User, error)
	// Get user by name
	GetUserByName(name string) (User, error)

	// New user, retuns user with new id
	NewUser(user User) (User, error)
	// Update user, only update the name
	UpdateUser(user User) error
	// Delete user
	DeleteUser(id int) error

	// Categories relations to media
	// Add category to media
	AddCategoryToMedia(mediaid int, categoryid int) error
	// Remove category from media
	RemoveCategoryFromMedia(mediaid int, categoryid int) error
	// Remova all categories from media
	RemoveAllCategoriesFromMedia(mediaid int) error
	// Remove all media from category
	RemoveAllMediaFromCategory(categoryid int) error

	// Users relations to media
	// Add Media to user
	AddMediaToUser(mediaid int, userid int) error
	// Remove Media from user
	RemoveMediaFromUser(mediaid int, userid int) error
	// Remove all Media from user
	RemoveAllMediaFromUser(userid int) error
	// Remove all Users from media
	RemoveAllUsersFromMedia(mediaid int) error
	Close() error
}

type DB struct {
	metaDB
	path string
}

func NewDB(dbmeta metaDB, path string) *DB {

	return &DB{dbmeta, path}
}

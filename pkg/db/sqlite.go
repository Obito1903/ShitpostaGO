package db

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteDB struct {
	Database
	connection *sql.DB
	path       string
}

// Open a new sqlite connection
func NewSqlite(config Database) sqliteDB {
	var err error
	var sqlite sqliteDB

	sqlite.Database = config
	sqlite.path = path.Join(config.Folder, "/db.sqlite")
	sqlite.initFs()

	sqlite.connection, err = sql.Open("sqlite3", sqlite.path)

	checkDBerr(err)

	sqlite.checkDB()
	return sqlite
}

// Return the number of rows in table
func countRows(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkDBerr(err)
	}
	return count
}

// Check if the MediaType table exist and if not init it
func (db sqliteDB) checkMediaTypeTable() {
	mediaTypeTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'mediaTypes';")
	checkDBerr(err)

	if countRows(mediaTypeTable) == 0 {
		stmt, err := db.connection.Prepare(`CREATE TABLE mediaTypes (
			id TEXT PRIMARY KEY
			);`)
		checkDBerr(err)

		_, err = stmt.Exec()
		checkDBerr(err)

		stmt, err = db.connection.Prepare("INSERT INTO mediaTypes (id) VALUES(?);")
		checkDBerr(err)

		_, err = stmt.Exec(Video)
		checkDBerr(err)

		_, err = stmt.Exec(Image)
		checkDBerr(err)
	}
}

// Check if the Medias table exist and if not init it
func (db sqliteDB) checkMediaTable() {
	mediaTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'medias';")
	checkDBerr(err)
	if countRows(mediaTable) == 0 {
		stmt, err := db.connection.Prepare(`CREATE TABLE medias (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			og_name TEXT NOT NULL,
			name TEXT ,
			path TEXT,
			date datetime NOT NULL,
			type TEXT NOT NULL,

			FOREIGN KEY(type) REFERENCES mediaTypes(id)
			);`)
		checkDBerr(err)

		_, err = stmt.Exec()
		checkDBerr(err)
	}
}

// Check if the Medias table exist and if not init it
func (db sqliteDB) checkCategoriesTable() {
	caterogiesTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'categories';")
	checkDBerr(err)
	if countRows(caterogiesTable) == 0 {
		stmt, err := db.connection.Prepare(`CREATE TABLE categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
			);`)
		checkDBerr(err)

		_, err = stmt.Exec()
		checkDBerr(err)
	}
}

// Check if the Medias Caterories relation table exist and if not init it
func (db sqliteDB) checkMediaCategoriesTable() {
	mediaCaterogiesTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'mediaCategories';")
	checkDBerr(err)
	if countRows(mediaCaterogiesTable) == 0 {
		stmt, err := db.connection.Prepare(`CREATE TABLE mediaCategories (
			media_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,

			PRIMARY KEY (media_id,category_id),

			FOREIGN KEY (media_id) REFERENCES medias(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
			);`)
		checkDBerr(err)

		_, err = stmt.Exec()
		checkDBerr(err)
	}
}

// Check if the db contains all the necessary tables and if not inits them
func (db sqliteDB) checkDB() {
	db.checkMediaTypeTable()
	db.checkMediaTable()
	db.checkCategoriesTable()
	db.checkMediaCategoriesTable()
}

func (db sqliteDB) GetConfig() Database {
	return db.Database
}

func (db sqliteDB) NewMediaFromPath(filePath string) (Media, error) {
	exist, err := CheckFileExist(filePath)
	println(exist)
	if !exist {
		if err != nil {
			return Media{}, err
		}
		return Media{}, errors.New("no such file")
	}

	var destFolder string
	mediaType, _, err := FindMediaType(filePath)
	checkDBerr(err)
	switch mediaType {
	case Video:
		destFolder = path.Join(db.Folder, "/videos/")
	case Image:
		destFolder = path.Join(db.Folder, "/images/")
	}

	u, err := uuid.NewRandom()
	checkDBerr(err)
	destFile, err := ImportFile(filePath, destFolder, u.String())
	checkDBerr(err)

	exist, err = CheckFileExist(destFile)
	if !exist {
		if err != nil {
			return Media{}, err
		}
		return Media{}, errors.New("import failed")
	}

	stmt, err := db.connection.Prepare("INSERT INTO medias (og_name,name,path,date,type) VALUES(?,?,?,?,?)")
	checkDBerr(err)
	res, err := stmt.Exec(path.Base(filePath), path.Base(filePath), destFile, time.Now(), mediaType)
	checkDBerr(err)
	id, err := res.LastInsertId()
	checkDBerr(err)

	return db.GetMediaFromId(id), nil
}

func (db sqliteDB) UpdateMedia(media Media) Media {
	stmt, err := db.connection.Prepare("UPDATE medias SET name = ?,path = ? WHERE id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(media.Name, media.Path, media.Id)
	checkDBerr(err)
	return media
}

func (db sqliteDB) RemoveMedia(media Media) {
	stmt, err := db.connection.Prepare("DELETE FROM medias WHERE id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(media.Id)
	checkDBerr(err)

	stmt, err = db.connection.Prepare("DELETE FROM mediaCategories WHERE media_id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(media.Id)
	checkDBerr(err)
}

func (db sqliteDB) GetMediaFromId(id int64) Media {
	var media Media
	stmt := db.connection.QueryRow("SELECT id,og_name,name,path,date,type FROM medias WHERE id = ?", id)
	err := stmt.Scan(&media.Id, &media.Og_name, &media.Name, &media.Path, &media.Date, &media.Type_)
	checkDBerr(err)

	fmt.Println(db.GetCategoriesFromMedia(media))

	media.Catergories = db.GetCategoriesFromMedia(media)

	return media
}

func (db sqliteDB) AddCategoryToMedia(media Media, category Category) Media {
	db.addCaterogyLink(media, category)
	media.Catergories = append(media.Catergories, category)
	return media
}

func (db sqliteDB) RemoveCategoryFromMedia(media Media, category Category) Media {
	db.removeCaterogyLink(media, category)
	return db.GetMediaFromId(media.Id)
}

//-----------------------------------------
// CategoryLinks
//-----------------------------------------

func (db sqliteDB) addCaterogyLink(media Media, category Category) {
	stmt, err := db.connection.Prepare("INSERT INTO mediaCategories (media_id,category_id) VALUES(?,?)")
	checkDBerr(err)
	_, err = stmt.Exec(media.Id, category.Id)
	checkDBerr(err)
}

func (db sqliteDB) removeCaterogyLink(media Media, category Category) {
	stmt, err := db.connection.Prepare("DELETE FROM mediaCategories WHERE media_id = ? AND category_id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(media.Id, category.Id)
	checkDBerr(err)
}

//-----------------------------------------
// Categories
//-----------------------------------------

func (db sqliteDB) NewCategory(name string) Category {
	stmt, err := db.connection.Prepare("INSERT INTO categories (name) VALUES(?)")
	checkDBerr(err)
	res, err := stmt.Exec(name)
	checkDBerr(err)
	id, err := res.LastInsertId()
	checkDBerr(err)

	return db.GetCategoryFromId(id)
}

func (db sqliteDB) RemoveCategory(category Category) {
	stmt, err := db.connection.Prepare("DELETE FROM categories WHERE id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(category.Id)
	checkDBerr(err)

	stmt, err = db.connection.Prepare("DELETE FROM mediaCategories WHERE category_id = ?")
	checkDBerr(err)
	_, err = stmt.Exec(category.Id)
	checkDBerr(err)
}

func (db sqliteDB) GetCategoryFromId(id int64) Category {
	var category Category
	stmt := db.connection.QueryRow("SELECT id,name FROM categories WHERE id = ?", id)
	err := stmt.Scan(&category.Id, &category.Name)
	checkDBerr(err)

	return category
}

func (db sqliteDB) GetCategoriesFromId(ids []int64) []Category {
	vals := []interface{}{}
	sqlStr := "SELECT id,name FROM categories WHERE id in ("
	for _, v := range ids {
		vals = append(vals, strconv.Itoa(int(v)))
		sqlStr += "?,"
	}
	sqlStr = strings.TrimRight(sqlStr, ",")
	sqlStr += ")"
	fmt.Println(sqlStr)
	stmt, err := db.connection.Prepare(sqlStr)
	checkDBerr(err)
	defer stmt.Close()

	rows, err := stmt.Query(vals...)
	checkDBerr(err)

	var categories []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name)
		checkDBerr(err)
		categories = append(categories, category)
	}

	return categories
}

func (db sqliteDB) GetCategoriesFromMedia(media Media) []Category {
	rows, err := db.connection.Query("SELECT category_id FROM mediaCategories WHERE media_id = ?", media.Id)
	checkDBerr(err)

	var categoryIds []int64
	for rows.Next() {
		var link int64
		err = rows.Scan(&link)
		checkDBerr(err)
		categoryIds = append(categoryIds, link)
	}
	fmt.Println(categoryIds)
	return db.GetCategoriesFromId(categoryIds)
}

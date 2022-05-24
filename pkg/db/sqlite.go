package db

import (
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SQLRowNotFound    ErrorCode = "SQLRowNotFound"
	SQLColumnNotFound ErrorCode = "SQLColumnNotFound"
	SQLScan           ErrorCode = "SQLScan"
	SQLQuery          ErrorCode = "SQLQuery"
	SQLPrepare        ErrorCode = "SQLPrepare"
	SQLExec           ErrorCode = "SQLExec"
	SQLTableInit      ErrorCode = "SQLTableInit"
)

type sqliteDB struct {
	Database
	connection *sql.DB
	path       string
}

// Open a new sqlite connection
func NewSqlite(config Database) (sqlite sqliteDB, dberr *DBError) {
	var err error

	sqlite.Database = config
	sqlite.path = path.Join(config.Folder, "/db.sqlite")
	sqlite.initFs()

	sqlite.connection, err = sql.Open("sqlite3", sqlite.path)
	dberr = checkErr(err, "Failed to open sqlite3 connection", InitFailed)

	sqlite.checkDB()
	return sqlite, dberr
}

// Return the number of rows in table
func countRows(rows *sql.Rows) (count int, dberr *DBError) {
	for rows.Next() {
		err := rows.Scan(&count)
		dberr = checkErr(err, "Failed to scan Count result", SQLScan)
	}
	return count, dberr
}

// Check if the MediaType table exist and if not init it
func (db sqliteDB) checkMediaTypeTable() (dberr *DBError) {
	mediaTypeTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'mediaTypes';")
	if dberr = checkErr(err, "Failed to Query sqlite_master table", SQLQuery); dberr != (*DBError)(nil) {
		return dberr
	}

	if count, err := countRows(mediaTypeTable); count == 0 {
		if err != nil {
			return err
		}
		stmt, err := db.connection.Prepare(`CREATE TABLE mediaTypes (
			id TEXT PRIMARY KEY
			);`)
		if dberr = checkErr(err, "Failed to prepare 'mediaTypes' init request", SQLPrepare); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec()
		if dberr = checkErr(err, "Failed to execute 'mediaTypes' init request", SQLExec); dberr != (*DBError)(nil) {
			return dberr
		}

		stmt, err = db.connection.Prepare("INSERT INTO mediaTypes (id) VALUES(?);")
		if dberr = checkErr(err, "Failed to prepare 'mediaTypes' insert request", SQLPrepare); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec(Video)
		if dberr = checkErr(err, "Failed to insert Video in 'mediaTypes'", SQLExec); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec(Image)
		if dberr = checkErr(err, "Failed to insert Image in 'mediaTypes'", SQLExec); dberr != (*DBError)(nil) {
			return dberr
		}
	}

	return nil
}

// Check if the Medias table exist and if not init it
func (db sqliteDB) checkMediaTable() (dberr *DBError) {
	mediaTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'medias';")
	checkDBerr(err)
	if count, err := countRows(mediaTable); count == 0 {
		if err != nil {
			return err
		}
		stmt, err := db.connection.Prepare(`CREATE TABLE medias (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			og_name TEXT NOT NULL,
			name TEXT ,
			path TEXT,
			date datetime NOT NULL,
			type TEXT NOT NULL,

			FOREIGN KEY(type) REFERENCES mediaTypes(id)
			);`)
		if dberr = checkErr(err, "Failed to prepare 'medias' init request", SQLPrepare); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec()
		if dberr = checkErr(err, "Failed to execute 'medias' init request", SQLExec); err != nil {
			return dberr
		}
	}
	return nil
}

// Check if the Medias table exist and if not init it
func (db sqliteDB) checkCategoriesTable() (dberr *DBError) {
	caterogiesTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'categories';")
	checkDBerr(err)
	if count, err := countRows(caterogiesTable); count == 0 {
		if err != nil {
			return err
		}
		stmt, err := db.connection.Prepare(`CREATE TABLE categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
			);`)
		if dberr = checkErr(err, "Failed to prepare 'categories' init request", SQLPrepare); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec()
		if dberr = checkErr(err, "Failed to execute 'categories' init request", SQLExec); dberr != (*DBError)(nil) {
			return dberr
		}
	}
	return nil
}

// Check if the Medias Caterories relation table exist and if not init it
func (db sqliteDB) checkMediaCategoriesTable() (dberr *DBError) {
	mediaCaterogiesTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'mediaCategories';")
	checkDBerr(err)
	if count, err := countRows(mediaCaterogiesTable); count == 0 {
		if err != nil {
			return err
		}
		stmt, err := db.connection.Prepare(`CREATE TABLE mediaCategories (
			media_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,

			PRIMARY KEY (media_id,category_id),

			FOREIGN KEY (media_id) REFERENCES medias(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
			);`)
		if dberr = checkErr(err, "Failed to prepare 'mediaCategories' init request", SQLPrepare); dberr != (*DBError)(nil) {
			return dberr
		}

		_, err = stmt.Exec()
		if dberr = checkErr(err, "Failed to execute 'mediaCategories' init request", SQLExec); dberr != (*DBError)(nil) {
			return dberr
		}
	}
	return nil
}

// Check if the db contains all the necessary tables and if not inits them
func (db sqliteDB) checkDB() (dberr *DBError) {
	if dberr = checkErr(db.checkMediaTypeTable(), "Failed MediaType table initialisation", SQLTableInit); dberr != (*DBError)(nil) {
		return dberr
	}
	if dberr = checkErr(db.checkMediaTable(), "Failed Media table initialisation", SQLTableInit); dberr != (*DBError)(nil) {
		return dberr
	}
	if dberr = checkErr(db.checkCategoriesTable(), "Failed Categories table initialisation", SQLTableInit); dberr != (*DBError)(nil) {
		return dberr
	}
	if dberr = checkErr(db.checkMediaCategoriesTable(), "Failed Categories Links table initialisation", SQLTableInit); dberr != (*DBError)(nil) {
		return dberr
	}

	return nil
}

func (db sqliteDB) GetConfig() Database {
	return db.Database
}

func (db sqliteDB) NewMediaFromPath(filePath string) (Media, *DBError) {
	exist, dberr := CheckFileExist(filePath)
	if (!exist) || (dberr != (*DBError)(nil)) {
		return Media{}, &DBError{dberr, fmt.Sprintf("No such file : %s", filePath), FileNotFound}
	}

	var destFolder string
	mediaType, _, dberr := FindMediaType(filePath)
	if dberr != (*DBError)(nil) {
		return Media{}, dberr
	}
	switch mediaType {
	case Video:
		destFolder = path.Join(db.Folder, "/videos/")
	case Image:
		destFolder = path.Join(db.Folder, "/images/")
	}

	u, err := uuid.NewRandom()
	if dberr = checkErr(err, "Failed to generate new  UUID", ImportFailed); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}
	destFile, err := ImportFile(filePath, destFolder, u.String())
	if dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	exist, err = CheckFileExist(destFile)
	if (!exist) || (dberr != (*DBError)(nil)) {
		return Media{}, &DBError{dberr, fmt.Sprintf("No such file : %s", filePath), ImportFailed}
	}

	stmt, err := db.connection.Prepare("INSERT INTO medias (og_name,name,path,date,type) VALUES(?,?,?,?,?)")
	if dberr = checkErr(err, "Failed to prepare 'medias' insert request", SQLPrepare); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}
	res, err := stmt.Exec(path.Base(filePath), path.Base(filePath), destFile, time.Now(), mediaType)
	if dberr = checkErr(err, "Failed to insert media", SQLExec); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}
	id, err := res.LastInsertId()
	if dberr = checkErr(err, "Failed to retrieve last insert id", SQLExec); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	return db.GetMediaFromId(id)
}

func (db sqliteDB) UpdateMedia(media Media) (Media, *DBError) {
	stmt, err := db.connection.Prepare("UPDATE medias SET name = ?,path = ? WHERE id = ?")
	if dberr := checkErr(err, "Failed to prepare 'medias' update request", SQLPrepare); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}
	_, err = stmt.Exec(media.Name, media.Path, media.Id)
	if dberr := checkErr(err, "Failed to update media", SQLExec); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	return media, nil
}

func (db sqliteDB) RemoveMedia(media Media) (dberr *DBError) {
	stmt, err := db.connection.Prepare("DELETE FROM medias WHERE id = ?")
	if dberr = checkErr(err, "Failed to prepare 'medias' delete request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(media.Id)
	if dberr := checkErr(err, "Failed to delete media", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}

	stmt, err = db.connection.Prepare("DELETE FROM mediaCategories WHERE media_id = ?")
	if dberr = checkErr(err, "Failed to prepare 'mediaCategories' delete request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(media.Id)
	if dberr := checkErr(err, "Failed to delete mediaCategories", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}

	return nil
}

func (db sqliteDB) GetMediaFromId(id int64) (media Media, dberr *DBError) {
	stmt := db.connection.QueryRow("SELECT id,og_name,name,path,date,type FROM medias WHERE id = ?", id)
	err := stmt.Scan(&media.Id, &media.Og_name, &media.Name, &media.Path, &media.Date, &media.Type_)
	if dberr = checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	media.Catergories, dberr = db.GetCategoriesFromMedia(media)
	if dberr := checkErr(dberr, fmt.Sprintf("Failed to query categories for : %d", media.Id), SQLQuery); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	return media, nil
}

func (db sqliteDB) GetMediasFromCats(categories []Category, limit int, offset int) (medias []Media, dberr *DBError) {
	values := []interface{}{}
	sqlStmt := "SELECT media_id FROM mediaCategories GROUP BY media_id HAVING "
	for i, cat := range categories {
		values = append(values, strconv.Itoa(int(cat.Id)))
		sqlStmt += "sum(CASE WHEN category_id = ? then 1 else 0 end) > 0"
		if i < len(categories)-1 {
			sqlStmt += " and\n"
		}
	}
	sqlStmt += "\n LIMIT ? OFFSET ?;"
	values = append(values, limit, offset)
	stmt, err := db.connection.Prepare(sqlStmt)
	if dberr = checkErr(err, "Failed to prepare 'mediaCategories' select request", SQLPrepare); dberr != (*DBError)(nil) {
		return []Media{}, dberr
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if dberr = checkErr(err, "Failed to Query mediaCategories table", SQLQuery); dberr != (*DBError)(nil) {
		return []Media{}, dberr
	}

	for rows.Next() {
		var media_id int64
		err = rows.Scan(&media_id)
		if dberr = checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
			return []Media{}, dberr
		}
		media, dberr := db.GetMediaFromId(media_id)
		if dberr != (*DBError)(nil) {
			return []Media{}, dberr
		}
		medias = append(medias, media)
	}
	return medias, nil
}

func (db sqliteDB) AddCategoryToMedia(media Media, category Category) (Media, *DBError) {
	dberr := db.addCaterogyLink(media, category)
	if dberr != (*DBError)(nil) {
		return media, dberr
	}

	media.Catergories = append(media.Catergories, category)
	return media, nil
}

func (db sqliteDB) RemoveCategoryFromMedia(media Media, category Category) (Media, *DBError) {
	dberr := db.removeCaterogyLink(media, category)
	if dberr != (*DBError)(nil) {
		return media, dberr
	}
	return db.GetMediaFromId(media.Id)
}

func (db sqliteDB) GetRandomMedia(type_ MediaType) (media Media, dberr *DBError) {

	stmt := db.connection.QueryRow("SELECT id,og_name,name,path,date,type FROM medias ORDER BY RANDOM() LIMIT 1;")
	err := stmt.Scan(&media.Id, &media.Og_name, &media.Name, &media.Path, &media.Date, &media.Type_)
	if dberr = checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	media.Catergories, dberr = db.GetCategoriesFromMedia(media)
	if dberr != (*DBError)(nil) {
		return Media{}, dberr
	}

	return media, nil
}

//-----------------------------------------
// CategoryLinks
//-----------------------------------------

func (db sqliteDB) addCaterogyLink(media Media, category Category) (dberr *DBError) {
	stmt, err := db.connection.Prepare("INSERT INTO mediaCategories (media_id,category_id) VALUES(?,?)")
	if dberr = checkErr(err, "Failed to prepare 'mediaCategories' insert request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(media.Id, category.Id)
	if dberr := checkErr(err, "Failed to insert link", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}
	return nil
}

func (db sqliteDB) removeCaterogyLink(media Media, category Category) (dberr *DBError) {
	stmt, err := db.connection.Prepare("DELETE FROM mediaCategories WHERE media_id = ? AND category_id = ?")
	if dberr = checkErr(err, "Failed to prepare 'mediaCategories' delete request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(media.Id, category.Id)
	if dberr := checkErr(err, "Failed to delete link", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}

	return nil
}

//-----------------------------------------
// Categories
//-----------------------------------------

func (db sqliteDB) NewCategory(name string) (Category, *DBError) {
	stmt, err := db.connection.Prepare("INSERT INTO categories (name) VALUES(?)")
	if dberr := checkErr(err, "Failed to prepare 'categories' insert request", SQLPrepare); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}
	res, err := stmt.Exec(name)
	if dberr := checkErr(err, "Failed to insert category", SQLExec); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}
	id, err := res.LastInsertId()
	if dberr := checkErr(err, "Failed to retrieve last insert id", SQLExec); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}

	return db.GetCategoryFromId(id)
}

func (db sqliteDB) UpdateCategory(category Category) (Category, *DBError) {
	stmt, err := db.connection.Prepare("UPDATE categories SET name = ? WHERE id = ?")
	if dberr := checkErr(err, "Failed to prepare 'categories' update request", SQLPrepare); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}
	_, err = stmt.Exec(category.Name, category.Id)
	if dberr := checkErr(err, "Failed to update category", SQLExec); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}

	return category, nil
}

func (db sqliteDB) RemoveCategory(category Category) (dberr *DBError) {
	stmt, err := db.connection.Prepare("DELETE FROM categories WHERE id = ?")
	if dberr = checkErr(err, "Failed to prepare 'categories' delete request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(category.Id)
	if dberr := checkErr(err, "Failed to delete category", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}

	stmt, err = db.connection.Prepare("DELETE FROM mediaCategories WHERE category_id = ?")
	if dberr = checkErr(err, "Failed to prepare 'mediaCategories' delete request", SQLPrepare); dberr != (*DBError)(nil) {
		return dberr
	}
	_, err = stmt.Exec(category.Id)
	if dberr := checkErr(err, "Failed to delete mediaCategories", SQLExec); dberr != (*DBError)(nil) {
		return dberr
	}

	return nil
}

func (db sqliteDB) GetCategoryFromId(id int64) (category Category, dberr *DBError) {
	stmt := db.connection.QueryRow("SELECT id,name FROM categories WHERE id = ?", id)
	err := stmt.Scan(&category.Id, &category.Name)
	if dberr = checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
		return Category{}, dberr
	}

	return category, nil
}

func (db sqliteDB) GetCategoriesFromIds(ids []int64) (categories []Category, dberr *DBError) {
	values := []interface{}{}
	//Prepare the statement
	sqlStmt := "SELECT id,name FROM categories WHERE id in ("
	for _, v := range ids {
		values = append(values, strconv.Itoa(int(v)))
		sqlStmt += "?,"
	}
	sqlStmt = strings.TrimRight(sqlStmt, ",")
	sqlStmt += ")"
	stmt, err := db.connection.Prepare(sqlStmt)
	if dberr = checkErr(err, "Failed to prepare 'categories' select request", SQLPrepare); dberr != (*DBError)(nil) {
		return []Category{}, dberr
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if dberr = checkErr(err, "Failed to Query categories table", SQLQuery); dberr != (*DBError)(nil) {
		return []Category{}, dberr
	}

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name)
		if dberr := checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
			return []Category{}, dberr
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db sqliteDB) GetCategoriesFromMedia(media Media) ([]Category, *DBError) {
	rows, err := db.connection.Query("SELECT category_id FROM mediaCategories WHERE media_id = ?", media.Id)
	if dberr := checkErr(err, "Failed to Query mediaCategories table", SQLQuery); dberr != (*DBError)(nil) {
		return []Category{}, dberr
	}

	var categoryIds []int64
	for rows.Next() {
		var link int64
		err = rows.Scan(&link)
		if dberr := checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
			return []Category{}, dberr
		}
		categoryIds = append(categoryIds, link)
	}
	return db.GetCategoriesFromIds(categoryIds)
}

func (db sqliteDB) GetCategories() (categories []Category, dberr *DBError) {
	rows, err := db.connection.Query("SELECT id,name FROM categories;")
	if dberr := checkErr(err, "Failed to Query categories table", SQLQuery); dberr != (*DBError)(nil) {
		return []Category{}, dberr
	}

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name)
		if dberr := checkErr(err, "Failed to scan request result", SQLScan); dberr != (*DBError)(nil) {
			return []Category{}, dberr
		}
		categories = append(categories, category)
	}

	return categories, nil
}

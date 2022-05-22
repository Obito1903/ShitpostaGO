package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteDB struct {
	connection *sql.DB
	path       string
}

// Open a new sqlite connection
func NewSqlite(path string) sqliteDB {
	var err error
	var sqlite sqliteDB

	sqlite.path = path
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
	mediaTypeTable, err := db.connection.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'mediaType';")
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
			name TEXT,
			path TEXT NOT NULL,
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
			name TEXT NOT NULL
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

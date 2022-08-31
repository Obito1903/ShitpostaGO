package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteinitscript = `
CREATE TABLE IF NOT EXISTS metadatas (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	og_name TEXT NOT NULL UNIQUE,
	name TEXT UNIQUE,
	media_type INTEGER NOT NULL,
	file_type TEXT NOT NULL,
	date_added DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_created DATETIME
);

CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	hash_pass TEXT NOT NULL,
	permission INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata_category (
	metadata_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,

	CONSTRAINT meta_cat_pk PRIMARY KEY (metadata_id,category_id),
	UNIQUE(metadata_id,category_id) ON CONFLICT REPLACE

	FOREIGN KEY (metadata_id) REFERENCES metadatas(id),
	FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS user_fav (
	user_id INTEGER NOT NULL,
	meta_id INTEGER NOT NULL,

	CONSTRAINT  user_fav_pk PRIMARY KEY (user_id,meta_id),
	UNIQUE(user_id,meta_id) ON CONFLICT REPLACE


	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (meta_id) REFERENCES metadatas(id)
);
`
)

type SqliteDriver struct {
	db *sql.DB
}

func NewSqliteDriver(dbpath string) (*SqliteDriver, error) {
	needinit := false
	if _, err := os.Stat(dbpath); errors.Is(err, os.ErrNotExist) {
		needinit = true
	}

	Log.Info().Msgf("Connecting to sqlite database: %s", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	if needinit {
		Log.Info().Msg("Initializing sqlite database")
		_, err := db.Exec(sqliteinitscript)
		if err != nil {
			return nil, err
		}
	}

	return &SqliteDriver{db}, nil
}

func (d *SqliteDriver) Close() error {
	return d.db.Close()
}

//-------------------------------
// 			  Media
//-------------------------------

func scanMediaRow(row *sql.Row) (Metadata, error) {
	var metadata Metadata
	err := row.Scan(
		&metadata.Id,
		&metadata.OgName,
		&metadata.Name,
		&metadata.MediaType,
		&metadata.FileType,
		&metadata.DateAdded,
		&metadata.DateCreated,
	)
	return metadata, err
}

func scanMediaRows(rows *sql.Rows) ([]Metadata, error) {
	var medias []Metadata
	for rows.Next() {
		var metadata Metadata
		err := rows.Scan(
			&metadata.Id,
			&metadata.OgName,
			&metadata.Name,
			&metadata.MediaType,
			&metadata.FileType,
			&metadata.DateAdded,
			&metadata.DateCreated,
		)
		if err != nil {
			return nil, err
		}
		medias = append(medias, metadata)
	}
	rows.Close()
	return medias, nil
}

func (d *SqliteDriver) GetMedia(id int) (Metadata, error) {
	row := d.db.QueryRow("SELECT * FROM metadatas WHERE id = ?", id)
	return scanMediaRow(row)
}

func (d *SqliteDriver) GetMediaByName(name string) (Metadata, error) {
	row := d.db.QueryRow("SELECT * FROM metadatas WHERE name = ?", name)
	return scanMediaRow(row)
}

func (d *SqliteDriver) GetMediaRandom(mediatype MediaType) (Metadata, error) {
	row := d.db.QueryRow("SELECT * FROM metadatas WHERE media_type = ? ORDER BY RANDOM() LIMIT 1;", mediatype)
	return scanMediaRow(row)
}

func (d *SqliteDriver) GetMediaRandomByCategory(mediatype MediaType, categoryid int) (Metadata, error) {
	row := d.db.QueryRow("SELECT m.* FROM metadata_category mc INNER JOIN metadatas m ON mc.metadata_id = m.id WHERE mc.category_id IN (?) GROUP BY metadata_id ORDER BY RANDOM() LIMIT 1;", mediatype, categoryid)
	return scanMediaRow(row)
}

func (d *SqliteDriver) GetMediasFromUser(user User) ([]Metadata, error) {
	rows, err := d.db.Query("SELECT m.* FROM user_fav uf INNER JOIN metadatas m ON uf.meta_id = m.id WHERE uf.user_id = ?;", user.Id)
	if err != nil {
		return nil, err
	}
	return scanMediaRows(rows)
}

func (d *SqliteDriver) GetMediasByType(mediatype MediaType) ([]Metadata, error) {
	rows, err := d.db.Query("SELECT * FROM metadatas WHERE media_type = ?;", mediatype)
	if err != nil {
		return nil, err
	}
	return scanMediaRows(rows)
}

func (d *SqliteDriver) GetMediasByTypeAndCategory(mediaType MediaType, categoryid int) ([]Metadata, error) {
	rows, err := d.db.Query("SELECT m.* FROM metadata_category mc INNER JOIN metadatas m ON mc.metadata_id = m.id WHERE mc.category_id IN (?);", mediaType, categoryid)
	if err != nil {
		return nil, err
	}
	return scanMediaRows(rows)
}

func (d *SqliteDriver) SearchMediaByName(name string) ([]Metadata, error) {
	rows, err := d.db.Query("SELECT * FROM metadatas WHERE name LIKE ?;", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	return scanMediaRows(rows)
}

func (d *SqliteDriver) NewMedia(metadata Metadata) (Metadata, error) {
	stmt, err := d.db.Prepare("INSERT INTO metadatas (og_name, name, media_type, file_type, date_created) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return metadata, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		metadata.OgName,
		metadata.Name,
		metadata.MediaType,
		metadata.FileType,
		metadata.DateCreated,
	)
	if err != nil {
		return metadata, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return metadata, err
	}
	metadata.Id = int(id)
	return metadata, err
}

func (d *SqliteDriver) UpdateMedia(metadata Metadata) error {
	stmt, err := d.db.Prepare("UPDATE metadatas SET name = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		metadata.Name,
		metadata.Id,
	)

	return err
}

func (d *SqliteDriver) DeleteMedia(id int) error {
	err := d.RemoveAllCategoriesFromMedia(id)
	if err != nil {
		return err
	}
	err = d.RemoveAllUsersFromMedia(id)
	if err != nil {
		return err
	}

	stmt, err := d.db.Prepare("DELETE FROM metadatas WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

//-------------------------------
// 			  Categories
//-------------------------------

func scanCategoryRow(row *sql.Row) (Category, error) {
	var category Category
	err := row.Scan(
		&category.Id,
		&category.Name,
	)
	return category, err
}

func scanCategoryRows(rows *sql.Rows) ([]Category, error) {
	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	rows.Close()
	return categories, nil
}

func (d *SqliteDriver) GetCategory(id int) (Category, error) {
	row := d.db.QueryRow("SELECT * FROM categories WHERE id = ?", id)
	return scanCategoryRow(row)
}

func (d *SqliteDriver) GetCategoryByName(name string) (Category, error) {
	row := d.db.QueryRow("SELECT * FROM categories WHERE name = ?", name)
	return scanCategoryRow(row)
}

func (d *SqliteDriver) GetCategoriesFromMedia(mediaid int) ([]Category, error) {
	rows, err := d.db.Query("SELECT c.* FROM metadata_category mc INNER JOIN categories c ON mc.category_id = c.id WHERE mc.metadata_id = ?;", mediaid)
	if err != nil {
		return nil, err
	}
	return scanCategoryRows(rows)
}

func (d *SqliteDriver) GetCategories() ([]Category, error) {
	rows, err := d.db.Query("SELECT * FROM categories;")
	if err != nil {
		return nil, err
	}
	return scanCategoryRows(rows)
}

func (d *SqliteDriver) NewCategory(category Category) (Category, error) {
	stmt, err := d.db.Prepare("INSERT INTO categories (name) VALUES (?);")
	if err != nil {
		return category, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(category.Name)
	if err != nil {
		return category, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return category, err
	}
	category.Id = int(id)
	return category, err
}

func (d *SqliteDriver) UpdateCategory(category Category) error {
	stmt, err := d.db.Prepare("UPDATE categories SET name = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(category.Name, category.Id)
	return err
}

func (d *SqliteDriver) DeleteCategory(id int) error {
	err := d.RemoveAllMediaFromCategory(id)
	if err != nil {
		return err
	}

	stmt, err := d.db.Prepare("DELETE FROM categories WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

//-------------------------------
// 			  Users
//-------------------------------

func scanUserRow(row *sql.Row) (User, error) {
	var user User
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.HashPass,
		&user.Permission,
	)
	return user, err
}

func (d *SqliteDriver) GetUser(id int) (User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	return scanUserRow(row)
}

func (d *SqliteDriver) GetUserByName(name string) (User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE name = ?", name)
	return scanUserRow(row)
}

func (d *SqliteDriver) NewUser(user User) (User, error) {
	stmt, err := d.db.Prepare("INSERT INTO users (name, hash_pass, permission) VALUES (?, ?, ?);")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.HashPass, user.Permission)
	if err != nil {
		return user, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return user, err
	}
	user.Id = int(id)
	return user, err
}

func (d *SqliteDriver) UpdateUser(user User) error {
	stmt, err := d.db.Prepare("UPDATE users SET name = ?, hash_pass = ?, permission = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name, user.HashPass, user.Permission, user.Id)
	return err
}

func (d *SqliteDriver) DeleteUser(id int) error {
	err := d.RemoveAllMediaFromUser(id)
	if err != nil {
		return err
	}

	stmt, err := d.db.Prepare("DELETE FROM users WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

//-------------------------------
// Catergories to media relations
//-------------------------------

func (d *SqliteDriver) AddCategoryToMedia(mediaid int, categoryid int) error {
	stmt, err := d.db.Prepare("INSERT OR REPLACE INTO metadata_category (metadata_id,category_id) VALUES (?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid, categoryid)
	return err
}

func (d *SqliteDriver) RemoveCategoryFromMedia(mediaid int, categoryid int) error {
	stmt, err := d.db.Prepare("DELETE FROM metadata_category WHERE metadata_id = ? AND category_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid, categoryid)
	return err
}

func (d *SqliteDriver) RemoveAllCategoriesFromMedia(mediaid int) error {
	stmt, err := d.db.Prepare("DELETE FROM metadata_category WHERE metadata_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid)
	return err
}

func (d *SqliteDriver) RemoveAllMediaFromCategory(categoryid int) error {
	stmt, err := d.db.Prepare("DELETE FROM metadata_category WHERE category_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(categoryid)
	return err
}

//-------------------------------
// User to media relations
//-------------------------------

func (d *SqliteDriver) AddMediaToUser(mediaid int, userid int) error {
	stmt, err := d.db.Prepare("INSERT OR REPLACE INTO user_fav (meta_id,user_id) VALUES (?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid, userid)
	return err
}

func (d *SqliteDriver) RemoveMediaFromUser(mediaid int, userid int) error {
	stmt, err := d.db.Prepare("DELETE FROM user_fav WHERE meta_id = ? AND user_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid, userid)
	return err
}

func (d *SqliteDriver) RemoveAllMediaFromUser(userid int) error {
	stmt, err := d.db.Prepare("DELETE FROM user_fav WHERE user_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userid)
	return err
}

func (d *SqliteDriver) RemoveAllUsersFromMedia(mediaid int) error {
	stmt, err := d.db.Prepare("DELETE FROM user_fav WHERE meta_id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mediaid)
	return err
}

// TODO: Add relations support

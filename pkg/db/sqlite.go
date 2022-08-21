package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDriver struct {
	db *sql.DB
}

func NewSqliteDriver(dbpath string) (*SqliteDriver, error) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
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
		&metadata.Type,
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
			&metadata.Type,
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
	stmt, err := d.db.Prepare("DELETE FROM metadatas WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type MediaInfo struct {
	Id        int
	Extension string
	Title     string
	Type      string
}

var db *sql.DB

func checkDBerr(err error) {
	if err != nil {
		log.Fatal("SQL ERROR : ", err)
	}
}

func countRows(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkDBerr(err)
	}
	return count
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./data/shit.db")
	checkDBerr(err)
	videoRows, err := db.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'videos';")
	checkDBerr(err)
	if countRows(videoRows) == 0 {
		stmt, err := db.Prepare("CREATE TABLE videos (id INTEGER PRIMARY KEY AUTOINCREMENT, extension TEXT, title TEXT, og_filename TEXT);")
		checkDBerr(err)
		_, err = stmt.Exec()
		checkDBerr(err)

		stmt, err = db.Prepare("CREATE UNIQUE INDEX idx_videos_id ON videos (id);")
		checkDBerr(err)
		_, err = stmt.Exec()
		checkDBerr(err)
	}
	videoRows.Close()
	imagesRows, err := db.Query("SELECT count(*) FROM sqlite_master WHERE type='table' AND name = 'images';")
	checkDBerr(err)
	if countRows(imagesRows) == 0 {
		stmt, err := db.Prepare("CREATE TABLE images (id INTEGER PRIMARY KEY AUTOINCREMENT, extension TEXT, title TEXT, og_filename TEXT);")
		checkDBerr(err)
		_, err = stmt.Exec()
		checkDBerr(err)

		stmt, err = db.Prepare("CREATE UNIQUE INDEX idx_images_id ON images (id);")
		checkDBerr(err)
		_, err = stmt.Exec()
		checkDBerr(err)
	}
	imagesRows.Next()
	checkDBerr(err)

}

func addFileToDb(db_table string, file string, ext string) int64 {

	stmt, err := db.Prepare("INSERT INTO " + db_table + "(extension, og_filename, title) VALUES(?,?,?);")
	checkDBerr(err)
	res, err := stmt.Exec(ext, file, "")
	checkDBerr(err)
	id, err := res.LastInsertId()
	checkDBerr(err)
	return id
}

func AddVideo(file string, newExt string) int64 {
	id := addFileToDb("videos", file, newExt)
	fmt.Println("Added video", file, "with id", id)
	return id
}

func AddImage(file string, newExt string) int64 {
	id := addFileToDb("images", file, newExt)
	fmt.Println("Added image", file, "with id", id)
	return id
}

func GetTableCount(table string) int {
	var nbRows int
	rows, err := db.Query("SELECT COUNT(*) FROM " + table + ";")
	checkDBerr(err)
	for rows.Next() {
		err = rows.Scan(&nbRows)
		checkDBerr(err)
	}
	rows.Close()
	return nbRows
}

func GetMediaInfo(dbName string, id int) (MediaInfo, error) {
	var info MediaInfo
	info.Type = dbName
	row := db.QueryRow("SELECT id,extension,title FROM "+dbName+" WHERE id = ?", id)
	err := row.Scan(&info.Id, &info.Extension, &info.Title)
	if err != nil {
		log.Print(err)
	}
	return info, err
}

func GetVideoInfo(id int) (MediaInfo, error) {
	return GetMediaInfo("videos", id)
}

func GetImageInfo(id int) (MediaInfo, error) {
	return GetMediaInfo("images", id)
}

func GetVideosCount() int {
	return GetTableCount("videos")
}

func GetImagesCount() int {
	return GetTableCount("images")
}

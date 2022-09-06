package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type FsMediaDB struct {
	fsMediaDB_interface
	Path         string
	RemoveSource bool
}

func NewFsMediaDB(path string) *FsMediaDB {
	return &FsMediaDB{Path: path, RemoveSource: true}
}

func (db *FsMediaDB) GetMedia(media Metadata) ([]byte, error) {
	return os.ReadFile(filepath.Join(db.Path, fmt.Sprintf("%d%s", media.Id, media.FileType)))
}

func (db *FsMediaDB) AddMedia(id int, source string) (Metadata, error) {
	mt := MediaType_Unknown
	var destination string
	mtype, err := mimetype.DetectFile(source)
	if err != nil {
		return Metadata{}, err
	}
	Log.Debug().Msgf("Adding media %s (%s)", filepath.Base(source), mtype.String())

	if strings.Contains(mtype.String(), "video") && !mtype.Is("video/mp4") {
		mt = MediaType_Video
		Log.Info().Msgf("Converting video to mp4 : %s", filepath.Base(source))
		destination = filepath.Join(db.Path, fmt.Sprintf("%d.mp4", id))
		err := ConvertVideoMedia(source, destination)
		if err != nil {
			return Metadata{}, err
		}
	}

	if strings.Contains(mtype.String(), "image") && !mtype.Is("image/jpeg") && !mtype.Is("image/gif") {
		mt = MediaType_Image
		Log.Info().Msgf("Converting image to jpeg : %s", filepath.Base(source))
		destination = filepath.Join(db.Path, fmt.Sprintf("%d.jpg", id))
		err := ConvertImageMedia(source, destination)
		if err != nil {
			return Metadata{}, err
		}
	} else {
		destination = filepath.Join(db.Path, fmt.Sprintf("%d%s", id, filepath.Ext(source)))
		err := os.Rename(source, destination)
		if err != nil {
			return Metadata{}, err
		}
	}

	if db.RemoveSource {
		err = os.Remove(source)
		if err != nil {
			return Metadata{}, err
		}
	}

	return Metadata{
		Id:        id,
		OgName:    filepath.Base(source),
		Name:      filepath.Base(source[:len(source)-len(filepath.Ext(source))]),
		MediaType: mt,
		FileType:  filepath.Ext(destination),
	}, nil
}

func (db *FsMediaDB) RemoveMedia(media Metadata) error {
	return os.Remove(filepath.Join(db.Path, fmt.Sprintf("%d%s", media.Id, media.FileType)))
}

// transcode video to x265 mp4
func ConvertVideoMedia(source string, destination string) error {

	err := ffmpeg_go.Input(source).
		Output(destination, ffmpeg_go.KwArgs{
			"c:v": "libx265",
			"crf": "20",
			"b:a": "128k",
			"c:a": "aac",
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	return err
}

// transcode image to jpeg
func ConvertImageMedia(source string, destination string) error {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(source)
	if err != nil {
		return err
	}
	mw.ResetIterator()

	err = mw.SetImageFormat("jpeg")
	if err != nil {
		return err
	}

	err = mw.SetImageCompressionQuality(90)
	if err != nil {
		return err
	}
	err = mw.WriteImage(destination)
	if err != nil {
		return err
	}

	return nil
}

func (db *FsMediaDB) GetThumbnail(media Metadata) ([]byte, error) {
	return nil, nil
}

package db

import (
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type FsMediaDB struct {
	fsMediaDB_interface
	Path string
}

func NewFsMediaDB(path string) *FsMediaDB {
	return &FsMediaDB{Path: path}
}

func (db *FsMediaDB) GetMedia(media Metadata) ([]byte, error) {
	return nil, nil
}

func (db *FsMediaDB) AddMedia(source string) (Metadata, error) {

	return Metadata{}, nil
}

func (db *FsMediaDB) RemoveMedia(media Metadata) error {
	return nil
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

package db

import (
	"errors"
	"log"
	"os"
	"path"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func FindMediaType(filePath string) (MediaType, bool, error) {
	ext := path.Ext(filePath)
	var err error
	var mediaType MediaType
	needConvert := false
	switch ext {
	case ".mp4", ".webm":
		mediaType = Video
	case ".mov", ".gif":
		mediaType = Video
		needConvert = true
	case ".jpg", ".png", "webp":
		mediaType = Image
	case ".jpeg", ".jpg_large", ".jfif":
		mediaType = Image
		needConvert = true
	default:
		err = errors.New("unknown file format")
	}

	return mediaType, needConvert, err
}

func convertVideo(sourceFile string, destFolder string, destName string) string {

	outPath := path.Join(destFolder, destName+".webm")
	err := ffmpeg_go.Input(sourceFile).
		Output(outPath, ffmpeg_go.KwArgs{
			"c:v": "libvpx-vp9",
			"crf": "10",
			"b:a": "128k",
			"c:a": "libopus",
		}).
		OverWriteOutput().ErrorToStdOut().Run()
	checkConverterr(err)
	return outPath
}

func convertImage(sourceFile string, destFolder string, destName string) string {
	ext := path.Ext(sourceFile)
	switch ext {
	case ".jpeg", ".jpg_large", ".jfif":
		ext = ".jpg"

	}
	outPath := path.Join(destFolder, destName+ext)
	err := os.Rename(sourceFile, outPath)
	if err != nil {
		log.Fatal(err)
	}
	checkConverterr(err)
	return outPath
}

func moveFile(sourceFile string, destFile string) {
	err := os.Rename(sourceFile, destFile)
	if err != nil {
		log.Fatal(err)
	}
	checkConverterr(err)
}

func ImportFile(sourceFile string, destFolder string, destName string) (string, error) {
	mediaType, needConvert, _ := FindMediaType(sourceFile)
	outPath := path.Join(destFolder, destName+path.Ext(sourceFile))
	var err error
	switch mediaType {
	case Video:
		if needConvert {
			outPath = convertVideo(sourceFile, destFolder, destName)
		} else {
			moveFile(sourceFile, outPath)
		}
	case Image:
		if needConvert {
			outPath = convertImage(sourceFile, destFolder, destName)
		} else {
			moveFile(sourceFile, outPath)
		}
	case Unknown:
		err = errors.New("unknown file format")
	}
	return outPath, err
}

func checkConverterr(err error) {
	if err != nil {
		log.Fatal("CONVERTION ERROR : ", err)
	}
}

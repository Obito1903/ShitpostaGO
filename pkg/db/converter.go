package db

import (
	"fmt"
	"os"
	"path"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

const (
	CVFileNotFound    ErrorCode = "CVFileNotFound"
	CVUnkownMediaType ErrorCode = "CVUnkownMediaType"
	CVFFmpeg          ErrorCode = "CVFFmpeg"
	CVFileMove        ErrorCode = "CVFileMove"
)

/*
 *  Find if the file is a video or an image
 *
 * Possible improvement :
 *   - Use mime types
 */
func FindMediaType(filePath string) (mediaType MediaType, needConvert bool, dberr *DBError) {
	ext := path.Ext(filePath)
	needConvert = false
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
		dberr = &DBError{
			nil,
			fmt.Sprintf("Unkown type for : %s", filePath),
			CVUnkownMediaType,
		}
	}

	return mediaType, needConvert, nil
}

// Convert video to webm format if the source file is not in a web friendly format
func convertVideo(sourceFile string, destFolder string, destName string) (string, *DBError) {

	outPath := path.Join(destFolder, destName+".webm")
	err := ffmpeg_go.Input(sourceFile).
		Output(outPath, ffmpeg_go.KwArgs{
			"c:v": "libvpx-vp9",
			"crf": "10",
			"b:a": "128k",
			"c:a": "libopus",
		}).
		OverWriteOutput().ErrorToStdOut().Run()
	return outPath, checkErr(err, "FFmpeg transcode failed", CVFFmpeg)
}

// Normalize image extension name like .jpeg, .jpg_large, .jfif to .jpeg
func convertImage(sourceFile string, destFolder string, destName string) (outPath string, dberr *DBError) {
	ext := path.Ext(sourceFile)
	switch ext {
	case ".jpeg", ".jpg_large", ".jfif":
		ext = ".jpg"
	}
	outPath = path.Join(destFolder, destName+ext)
	return outPath, moveFile(destFolder, outPath)
}

// Move file on the fs
func moveFile(sourceFile string, destFile string) (dberr *DBError) {
	err := os.Rename(sourceFile, destFile)
	return checkErr(err, fmt.Sprintf("Could not move file '%s' to '%s'", sourceFile, destFile), CVFileMove)
}

// Import file into the database by renaming and transcoding it if necessary
func ImportFile(sourceFile string, destFolder string, destName string) (outPath string, dberr *DBError) {
	mediaType, needConvert, _ := FindMediaType(sourceFile)
	outPath = path.Join(destFolder, destName+path.Ext(sourceFile))
	switch mediaType {
	case Video:
		if needConvert {
			outPath, dberr = convertVideo(sourceFile, destFolder, destName)
			if dberr != nil {
				return "", dberr
			}
		} else {
			if moveFile(sourceFile, outPath) != nil {
				return "", dberr
			}
		}
	case Image:
		if needConvert {
			outPath, dberr = convertImage(sourceFile, destFolder, destName)
			if dberr != nil {
				return "", dberr
			}
		} else {
			if moveFile(sourceFile, outPath) != nil {
				return "", dberr
			}
		}
	case Unknown:
		dberr = &DBError{
			nil,
			fmt.Sprintf("Unkown type for : %s", sourceFile),
			CVUnkownMediaType,
		}
	}
	return outPath, dberr
}

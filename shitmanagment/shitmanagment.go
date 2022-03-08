package shitmanagment

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Obito1903/shitpostaGo/shitmanagment/db"
)

func init() {
	for _, folder := range [2]string{"./data/img", "./data/video"} {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, os.ModePerm)
		}
	}
	ScanForNewShit()
}

func moveFile(file fs.DirEntry, destPath string) {
	err := os.Rename("./data/new/"+file.Name(), destPath)
	if err != nil {
		log.Fatal(err)
	}
}

func convertGifandMove(file fs.DirEntry, outPath string) {
	o, err := exec.Command(
		"ffmpeg",
		"-i", "./data/new/"+file.Name(),
		"-y",
		"-vcodec", "libx264",
		"-pix_fmt", "yuv420p",
		"-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2",
		"-profile:v", "baseline",
		"-x264opts", "cabac=0:bframes=0:ref=1:weightp=0:level=30:bitrate=2000:vbv_maxrate=2200:vbv_bufsize=3000",
		"-movflags", "faststart",
		"-pass", "1",
		"-strict", "experimental",
		outPath,
	).CombinedOutput()
	if err != nil {
		log.Println(err, o)
	}
	os.Remove("./data/new/" + file.Name())
}

func moveFileBasedOnExtension(file fs.DirEntry) bool {
	fileMoved := true
	ext := filepath.Ext(file.Name())
	switch ext {
	case ".mp4":
		moveFile(file, "./data/video/"+fmt.Sprintf("%d%v", db.AddVideo(file.Name(), ext), ext))
	case ".webm":
		moveFile(file, "./data/video/"+fmt.Sprintf("%d%v", db.AddVideo(file.Name(), ext), ext))
	case ".mov":
		moveFile(file, "./data/video/"+fmt.Sprintf("%d.mp4", db.AddVideo(file.Name(), ".mp4")))
	case ".gif":
		convertGifandMove(file, "./data/video/"+fmt.Sprintf("%d.mp4", db.AddVideo(file.Name(), ".mp4")))
	case ".jpg":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d%v", db.AddImage(file.Name(), ext), ext))
	case ".jpeg":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d.jpg", db.AddImage(file.Name(), ".jpg")))
	case ".jpg_large":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d.jpg", db.AddImage(file.Name(), ".jpg")))
	case ".jfif":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d.jpg", db.AddImage(file.Name(), ".jpg")))
	case ".png":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d%v", db.AddImage(file.Name(), ext), ext))
	case ".webp":
		moveFile(file, "./data/img/"+fmt.Sprintf("%d%v", db.AddImage(file.Name(), ext), ext))
	default:
		fileMoved = false
	}

	return fileMoved
}

func ScanForNewShit() {
	fmt.Println("Scaning New folder")
	folder, err := os.Open("./data/new")
	if err != nil {
		log.Fatal(err)
	}
	files, err := folder.ReadDir(-1)
	folder.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Type().IsRegular() {
			moveFileBasedOnExtension(file)
		}
	}
}

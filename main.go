package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Obito1903/shitpostaGo/shitmanagment"
	"github.com/Obito1903/shitpostaGo/shitmanagment/db"
)

func shit(w http.ResponseWriter, req *http.Request) {
	shitIDstr := req.URL.Query().Get("id")
	shitType := req.URL.Query().Get("type")
	if shitIDstr == "" || shitType == "" {
		//Get not set, send a 400 bad request
		http.Error(w, "Get 'id' or 'type not specified in url.", 400)
		return
	}

	mediaId, _ := strconv.Atoi(shitIDstr)

	mediaFile, err := db.GetMediaInfo(shitType, mediaId)
	if err != nil {
		http.Error(w, "No such file.", 400)
		return
	}
	fmt.Println("Client requests: " + shitIDstr)
	switch shitType {
	case "videos":
		http.ServeFile(w, req, "./data/video/"+shitIDstr+mediaFile.Extension)
	case "images":
		http.ServeFile(w, req, "./data/img/"+shitIDstr+mediaFile.Extension)
	default:
		http.Error(w, "Get 'id' or 'type not specified in url.", 400)
	}

}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func shitCount(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	shitType := req.FormValue("type")
	var count int
	switch shitType {
	case "videos":
		count = db.GetVideosCount()
	case "images":
		count = db.GetImagesCount()
	default:
		http.Error(w, "Get 'id' or 'type not specified in url.", 400)
		return
	}
	fmt.Fprintf(w, "%d", count)

}

func periodicScan() {
	for {
		shitmanagment.ScanForNewShit()
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	shitmanagment.ScanForNewShit()

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", fs)
	http.HandleFunc("/shit", shit)
	http.HandleFunc("/shitCount", shitCount)
	go periodicScan()
	http.ListenAndServe(":80", nil)
}

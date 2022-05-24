package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Obito1903/shitpostaGo/pkg/db"
	"github.com/gorilla/mux"
)

type HandlerContext struct {
	shitdb db.DatabaseInterface
}

type ApiError struct {
	ResponsCode int
	Message     string
}

func newApiErrorFromDbErr(dberr db.DBError) (apiErr ApiError) {
	log.Println(dberr)
	switch dberr.Code {
	case db.SQLScan:
		return ApiError{
			http.StatusNotFound,
			"Not Found",
		}
	default:
		return ApiError{
			http.StatusInternalServerError,
			"Unknown error",
		}
	}
}

func sendError(w http.ResponseWriter, req *http.Request, apiErr ApiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.ResponsCode)
	json.NewEncoder(w).Encode(apiErr)
}

func (ctx HandlerContext) sendMedia(w http.ResponseWriter, req *http.Request, media db.Media) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, req, media.Path)
}

func (ctx HandlerContext) getMediaById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])
	media, dberr := ctx.shitdb.GetMediaFromId(int64(id))
	if dberr == nil {
		ctx.sendMedia(w, req, media)
	} else {
		sendError(w, req, newApiErrorFromDbErr(*dberr))
	}
}

func (ctx HandlerContext) getRandomMedia(w http.ResponseWriter, req *http.Request) {
	media, dberr := ctx.shitdb.GetRandomMedia(db.Video)
	if dberr == nil {
		ctx.sendMedia(w, req, media)
	} else {
		sendError(w, req, newApiErrorFromDbErr(*dberr))
	}
}

func (ctx HandlerContext) getCategories(w http.ResponseWriter, req *http.Request) {
	categories, _ := ctx.shitdb.GetCategories()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (ctx HandlerContext) Register() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/shit/{id:[0-9]+}", ctx.getMediaById)
	router.HandleFunc("/shit/getRandom", ctx.getRandomMedia)
	router.HandleFunc("/shit/getCategories", ctx.getCategories)

	return router
}

func Start(folder string) {
	shitdb, _ := db.NewSqlite(db.Database{
		Folder: folder,
	})

	ctx := HandlerContext{
		shitdb: shitdb,
	}

	router := ctx.Register()
	// http.Handle("/", router)

	http.ListenAndServe(":8090", router)
}

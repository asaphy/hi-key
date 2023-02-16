package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/go-openapi/runtime/middleware"
)

// swagger:response Song
type Song struct {
	// in: string
	Title string `json:"title"`
	Key string `json:"key"`
	HighNote string `json:"high_note"`
	FirstChord string `json:"first_chord"`
}

// swagger:response CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

var Songs []Song


func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/v0/all", returnAllSongs).Methods(http.MethodGet)
	myRouter.HandleFunc("/api/v0/song/{title}", returnSingleSong).Methods(http.MethodGet)
	myRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	myRouter.Handle("/docs", sh)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

// swagger:route GET /api/v0/song/{title} SongRequest
// Get song
// 
// security:
// - apiKey: []
//
// responses:
//  200: Song
//  400: CommonError
func returnSingleSong(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	title := vars["title"]
	// parse title to change - to spaces and lowercase everything
	for _, song := range Songs {
		if song.Title == title {
				json.NewEncoder(w).Encode(song)
		}
	}
}

// swagger:parameters SongRequest
type SongRequest struct {
	// Title of the song
	// in: path
	Title string `json:"title"validate:"required,min=2,max=100,alpha_space"`
}

// swagger:route GET /api/v0/all
// Get all songs
// 
// security:
// - apiKey: []
//
// responses:
//  200: Song
//  400: CommonError
func returnAllSongs(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllSongs")
	json.NewEncoder(w).Encode(Songs)
}


func main() {
		Songs = []Song{
				Song{Title: "Hello", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello1", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello2", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello3", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello4", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello4", Key: "A", HighNote: "A", FirstChord: "A"},
				Song{Title: "Hello5", Key: "A", HighNote: "A", FirstChord: "A"},
		}
    handleRequests()
}

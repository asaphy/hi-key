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


func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

//TODO: version the API
func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllSongs)
	myRouter.HandleFunc("/song/{title}", returnSingleSong).Methods(http.MethodGet)
	// documentation for developers
	myRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	myRouter.Handle("/docs", sh)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

// swagger:route GET /song/{title} SongRequest
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

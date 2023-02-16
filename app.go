package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/exp/slices"
	"strings"
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

var Songs = []Song{
	Song{Title: "Goodness of God", Key: "G", HighNote: "C", FirstChord: "G"},
	Song{Title: "Great Are You Lord", Key: "G", HighNote: "D", FirstChord: "C"},
	Song{Title: "Build My Life", Key: "A", HighNote: "F#", FirstChord: "A"},
}

var AllowedKeys = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

var KeyChords = map[string][]string{
	"C": []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"},
	"C#": []string{"C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C"},
	"D": []string{"D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#"},
	"D#": []string{"D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D"},
	"E": []string{"E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#"},
	"F": []string{"F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E"},
	"F#": []string{"F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F"},
	"G": []string{"G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#"},
	"G#": []string{"G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G"},
	"A": []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"},
	"A#": []string{"A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A"},
	"B": []string{"B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#"},
}


func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/v0/song/", returnSingleSong).Methods(http.MethodGet)
	myRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	myRouter.Handle("/docs", sh)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

// swagger:route GET /api/v0/song/ SongRequest
// Get song
// 
// security:
// - apiKey: []
//
// responses:
//  200: Song
//  400: CommonError
func returnSingleSong(w http.ResponseWriter, r *http.Request){
	title := r.URL.Query().Get("title")
	high_note := strings.ToUpper(r.URL.Query().Get("high_note"))
	if !slices.Contains(AllowedKeys, high_note) {
		errorMessage := fmt.Sprintln("The high_note provided:", high_note, "is not valid. Note: Currently, we do not support flats.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonError{Status: 400, Message: errorMessage})
		return
	}

	// parse title to change - to spaces and lowercase everything
	for _, song := range Songs {
		if strings.ToLower(song.Title) == strings.ToLower(title) {
				var step_change = getStepChange(high_note, song.HighNote)
				json.NewEncoder(w).Encode(Song{Title: song.Title, Key: getKeyFromStepChange(step_change, song.Key), HighNote: high_note, FirstChord: getFirstChordFromStepChange(step_change, song.FirstChord)})
			  return
		}
	}
	errorMessage := fmt.Sprintln("The song title provided", title, "is not in our library.")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(CommonError{Status: 400, Message: errorMessage})
}

func getStepChange(high_note string, original_high_note string) int {
  var step_change = 0
  var chords_for_key = KeyChords[original_high_note]
	for _, chord := range chords_for_key { 
		if chord == high_note {
			return step_change
		}
		step_change++
	}
	return 0
}

func getKeyFromStepChange(step_change int, original_key string) string {
	var key = KeyChords[original_key]
  return key[step_change]
}

func getFirstChordFromStepChange(step_change int, original_first_chord string) string {
	var key = KeyChords[original_first_chord]
	return key[step_change]
}


// swagger:parameters SongRequest
type SongRequest struct {
	// Title of the song
	// in: query
	Title string `json:"title"validate:"required,min=2,max=100,alpha_space"`
  // High note of singer
	// in: query
	HighNote string `json:"high_note"validate:"required,min=1,max=1,alpha_space"`
}


func main() {
		fmt.Println("Starting..")
    handleRequests()
}

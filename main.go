package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {

	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "424", Title: "Movie one", Director: &Director{Firstname: "johm", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "4242", Title: "Movie two", Director: &Director{Firstname: "samet", Lastname: "körük"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Sunucu başlatılıyor: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Film bulunamadı", http.StatusNotFound)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			movie.ID = params["id"]
			movies = append(movies, movie)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "Film bulunamadı", http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

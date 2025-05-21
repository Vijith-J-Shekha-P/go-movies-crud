package main

import (
	"encoding/json" //json encoding
	"fmt"           //fmt for printing
	"log"           //logging
	"math/rand"     //random number generation
	"net/http"      //http server
	"strconv"       //string conversion

	"github.com/gorilla/mux" //router - external package
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"` //unique movie identifier
	Title    string    `json:"title"`
	Director *Director `json:"director"` //pointer to Director struct
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie //slice of Movie structs

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set the content type to json
	json.NewEncoder(w).Encode(movies)                  //encode the movies slice to json and write it to the response
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set the content type to json
	params := mux.Vars(r)                              //get the URL parameters
	for _, movie := range movies {                     //loop through the movies slice
		if movie.ID == params["id"] { //check if the movie ID matches the URL parameter
			json.NewEncoder(w).Encode(movie) //encode the movie to json and write it to the response
			return                           //exit the function
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set the content type to json
	var movie Movie                                    //create a new Movie struct
	_ = json.NewDecoder(r.Body).Decode(&movie)         //decode the request body to the movie struct
	movie.ID = strconv.Itoa(rand.Intn(1000000))        //generate a random ID for the movie
	movies = append(movies, movie)                     //add the movie to the movies slice
	json.NewEncoder(w).Encode(movie)                   //encode the movie to json and write it to the response
	// w.WriteHeader(http.StatusCreated) //set the response status to 201 Created
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set the content type to json
	params := mux.Vars(r)                              //get the URL parameters
	for index, movie := range movies {                 //loop through the movies slice
		if movie.ID == params["id"] { //check if the movie ID matches the URL parameter
			movies = append(movies[:index], movies[index+1:]...) //remove the movie from the slice
			var movie Movie                                      //create a new Movie struct
			_ = json.NewDecoder(r.Body).Decode(&movie)           //decode the request body to the movie struct
			movie.ID = params["id"]                              //set the ID of the new movie to the ID from the URL parameter
			movies = append(movies, movie)                       //add the new movie to the movies slice
			json.NewEncoder(w).Encode(movie)                     //encode the new movie to json and write it to the response
			return                                               //exit the function
		}
	}
	json.NewEncoder(w).Encode(movies) //encode the updated movies slice to json and write it to the response
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //set the content type to json
	params := mux.Vars(r)                              //get the URL parameters
	for index, movie := range movies {                 //loop through the movies slice
		if movie.ID == params["id"] { //check if the movie ID matches the URL parameter
			movies = append(movies[:index], movies[index+1:]...) //remove the movie from the slice
			break                                                //exit the loop
		}
	}
	json.NewEncoder(w).Encode(movies) //encode the updated movies slice to json and write it to the response
}

func main() {
	r := mux.NewRouter() //create a new router

	//initialize the movies slice with some data
	movies = append(movies, Movie{ID: "1", Isbn: "438-1234567890", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})   //add a movie to the slice
	movies = append(movies, Movie{ID: "2", Isbn: "438-1234567891", Title: "Movie Two", Director: &Director{Firstname: "Jane", Lastname: "Smith"}}) //add another movie

	r.HandleFunc("/movies", getMovies).Methods("GET")           //handle GET request for /movies
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")       //handle GET request for /movies/{id}
	r.HandleFunc("/movies", createMovie).Methods("POST")        //handle POST request for /movies
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")    //handle PUT request for /movies/{id}
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") //handle DELETE request for /movies/{id}

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r)) //start the server on port 8000
}

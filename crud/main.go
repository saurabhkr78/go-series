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

// here we are building a movie management CRUD application
// here we will be using struct/slices as a database to perform CRUD operations
// we will be using external gorilla/mux package as it's not part of main package, to create routes and handle requests
// struct is like an object in other languages in which we can define multiple fields

// Exported field = starts with capital letter
// if kept samall letter it will be unexported and won't be accessible outside the package
type Movie struct {
	Title    string    `json:"title"`
	Year     int       `json:"year"`
	Director *Director `json:"director"`
	ID       string    `json:"id"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// this is our in-memory database
var movies []Movie

func main() {
	// initializing the router
	r := mux.NewRouter()

	// append a movie object to movies slice
	movies = append(movies, Movie{Title: "Inception", Year: 2020, Director: &Director{Firstname: "John", Lastname: "Doe"}, ID: "1"})
	movies = append(movies, Movie{Title: "Interstellar", Year: 2019, Director: &Director{Firstname: "Jane", Lastname: "Doe"}, ID: "2"})
	movies = append(movies, Movie{Title: "The Dark Knight", Year: 2018, Director: &Director{Firstname: "Johnny", Lastname: "Doe"}, ID: "3"})

	// route handlers / endpoints
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	//log and start the server with error handling
	fmt.Println("Starting server at port 8000")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

// funcs to handle requests

/*
Enginnering Side:
when a client hits an endpoint, the corresponding handler func is called and behind the scene http server spawns a goroutine to handle that request
it provides:
Go’s net/http server accepts a TCP connection and parse Parses the incoming HTTP request.

1.w: an interface to build/write the http response and it's network output stream
2.r: a full struct containing all the info about the http request(headers,path,body,method,host etc)
so all these are handled concurrently using goroutines for each request

and
W is responsewriter a concrete struct called * response
it contains a internal map of headers map[string][]string
set () writes into that map before writing the response body
and Setting Content-Type tells the browser:"Interpret the following bytes as JSON."
If you do w.Write() or json.Encoder.Encode() before setting headers → Go automatically writes default headers and you lose control.

json.NewEncoder(w):Creates a json.Encoder struct with target output your response writer w
.Encode(movies):Go reflects on movies Uses reflection (reflect package) to inspect its type and Iterates through each element and Converts each field to valid JSON
It writes JSON chunk-by-chunk to the network socket via via w.Write() i.e ResponseWriter

this brings efficiency no temp buffer,no full string building,low memory usage and suits for large data sets

After your handler finishes Go flushes writes,Closes or reuses the TCP connection depending on KeepAlive, and free the goroutine
*/
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) //encode the response we get from movies slice to json and send it as response

}

/*
above enginnering side applies to all below handler funcs too
Gorilla mux stores path variables in r.Context()
mux.Vars(r) retrieves those path variables as a map[string]string this uses This uses context.WithValue() under the hood and it's o(1) lookup

Each item is a copy of the struct element (not a pointer unless movie is pointer-type)
movies[:index] Up to (but not including) the element at index
movies[index+1:] Elements after the deleted index.
if append founds that If capacity is enough, reuse same array

	Go garbage collector cleans the  deleted items it later when unreachable
*/
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies { //here we dont need of index so we use _ blank identifier
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}

	}
	json.NewEncoder(w).Encode(&Movie{}) //if no movie found return empty movie object

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)   //decode the incoming json request body to movie struct
	movie.ID = strconv.Itoa(rand.Intn(10000000)) //mock id - generate random id
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie) //return the created movie as response
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			//delete the existing movie by appending with rest of the movies except the one to be updated
			movies = append(movies[:index], movies[index+1:]...)

			//now we gonna add the updated movie
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]        //keep the same id
			movies = append(movies, movie) //add the updated movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

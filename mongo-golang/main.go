package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"mongo-golang/controllers"
)

func main() {

	// Step 1: Create a new HTTP router instance
	r := httprouter.New()

	// Step 2: Create UserController with MongoDB session
	// getSession() establishes connection with MongoDB
	uc := controllers.NewUserController(getSession())

	// Step 3: Register API routes and map them to controller methods
	r.GET("/user/:id", uc.GetUser)       // fetch user by id
	r.POST("/user", uc.CreateUser)       // create new user
	r.DELETE("/user/:id", uc.DeleteUser) // delete user by id

	// Step 4: Start HTTP server and listen on port 9000
	http.ListenAndServe("localhost:9000", r)
}

// Step 5: Create and return MongoDB session
func getSession() *mgo.Session {

	// Step 6: Connect to MongoDB server
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err) // stop application if DB connection fails
	}

	// Step 7: Return MongoDB session
	return s
}

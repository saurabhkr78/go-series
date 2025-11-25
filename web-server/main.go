package main

import (
	"fmt"
	"log"
	"net/http"
)

// This function handles POST form submission coming from the client.
func formHandler(w http.ResponseWriter, r *http.Request) {

	// 1️.Parse the incoming form data from the request body
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form %v", err), http.StatusBadRequest)
		return
	}

	// 2️.Send a simple confirmation message back to the client
	fmt.Fprintf(w, "POST request successful\n")

	// 3️. Read form values (name and address) sent by the user
	name := r.FormValue("name")
	address := r.FormValue("address")

	// 4️. Send the extracted values back in the response
	fmt.Fprintf(w, "Hello %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

// ===================== HELLO HANDLER =====================
// This function handles GET request on /hello URL.
func helloHandler(w http.ResponseWriter, r *http.Request) {

	// 1️. Check if the client is requesting the correct path
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// 2️. Allow only GET method; reject all others
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}

	// 3️. If everything is valid, send a message back
	fmt.Fprintf(w, "Hello from server!")
}

// ===================== MAIN FUNCTION (SERVER SETUP) =====================
func main() {

	// 1️. Serve static files (HTML, CSS, JS) from the ./static directory
	fileserver := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileserver)

	// 2️. Map the /form URL to the formHandler function
	http.HandleFunc("/form", formHandler)

	// 3️. Map the /hello URL to the helloHandler function
	http.HandleFunc("/hello", helloHandler)

	// 4️. Print message that server is starting
	fmt.Println("starting server at port 8080")

	// 5️. Start the server on port 8080
	//     - http.ListenAndServe listens for incoming HTTP requests
	//     - If error occurs, handle it by exiting the app
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

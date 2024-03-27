package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Middleware function to log requests
func Logger(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)		
        next(w, r, ps)
    }
}

func Logger2(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Printf("Request2: %s %s\n", r.Method, r.URL.Path)		
        next(w, r, ps)
    }
}

// Handler function for the home route
func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "Welcome to the home page!")
}

// Handler function for the about route
func AboutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "About page")
}

func main() {
    // Create a new router
    router := httprouter.New()

    // Define route handlers
    router.GET("/", Logger2(Logger(HomeHandler)))
    router.GET("/about", Logger(AboutHandler))

    // Start the HTTP server
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8087", router)
}

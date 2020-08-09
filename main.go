package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)


// Function for each address
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	fmt.Fprint(w, "Welcome, it's home page")
}

// Main Function
func main(){
	// Declare router variable
	router := httprouter.New()

	// Connect each function to each address
	router.GET("/", Index)

	log.Fatal(http.ListenAndServe(":8080", router))
}

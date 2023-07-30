package main

import (
	"log"
	"net/http"

	"scet.com/routes"
)

func main() {
	r := routes.SetupRoutes()

	log.Println("Server started on :1414")
	log.Fatal(http.ListenAndServe(":1414", r))
}

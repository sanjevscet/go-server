package routes

import (
	"github.com/gorilla/mux"
	"scet.com/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlers.GetAllUsersHandler).Methods("GET")

	return r
}

package routes

import (
	"github.com/gorilla/mux"
	"scet.com/controllers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", controllers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", controllers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", controllers.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users", controllers.GetAllUsersHandler).Methods("GET")

	return r
}

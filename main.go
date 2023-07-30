package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var users = []User{}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newUser.ID = uuid.New().String()
	newUser.CreatedAt = time.Now()
	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	for _, user := range users {
		if user.ID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/users", getAllUsersHandler).Methods("GET")

	log.Println("Server started on :1414")
	log.Fatal(http.ListenAndServe(":1414", r))
}

package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"scet.com/utils"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name" validate:"required,min=3,max=15"`
	LastName  string    `json:"last_name" validate:"required,min=3,max=15"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

var users = []User{}

func validateUser(user User) validator.ValidationErrors {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the User struct
	validationErrors := validateUser(newUser)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Format validation errors as JSON with field names as keys
		errorMap := make(map[string]string)
		for _, err := range validationErrors {
			fieldName := err.Field()
			errorMessage := err.Error()
			errorMap[fieldName] = errorMessage
		}

		json.NewEncoder(w).Encode(errorMap)
		return
	}

	newUser.ID = utils.GetUUID()
	newUser.CreatedAt = time.Now()
	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var indexToRemove = -1
	for i, user := range users {
		if user.ID == userID {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		http.Error(w, "Invalid user id", http.StatusNotFound)
		return

	}

	users = append(users[:indexToRemove], users[indexToRemove+1:]...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"scet.com/myValidators"
	"scet.com/structs"
	"scet.com/utils"
)

var users = []structs.User{}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser structs.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the User struct
	validate := validator.New()
	if err := validate.Struct(newUser); err != nil {
		errorMessages := myValidators.ExtractUserValidationErrors(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessages)
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
	response := structs.DeleteUserResponse{
		Success: true,
		Msg:     fmt.Sprintf("User with %s removed successfully", userID),
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

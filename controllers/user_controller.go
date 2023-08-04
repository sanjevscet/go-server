package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"scet.com/db"
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

	if userID == "" {
		http.Error(w, "UserId not passed", http.StatusBadRequest)
		return
	}

	query := "SELECT id, first_name, last_name, email, created_at FROM go_users where id = ?"

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user structs.User
	var createdAt []uint8 // Temporary variable to store created_at as []uint8

	err = db.QueryRow(query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse createdAt ([]uint8) into time.Time
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the user struct to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	if userID == "" {
		http.Error(w, "UserId not passed", http.StatusBadRequest)
		return
	}

	query := "DELETE from go_users where id = ?"
	// Log the executed query
	log.Printf("Executing query: %s with userID: %s\n", query, userID)

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute the DELETE query
	result, err := db.Exec(query, userID)
	if err != nil {
		log.Fatal("Error executing the DELETE query: ", err)
	}

	// Check the number of rows affected by the DELETE operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Error getting the number of rows affected: ", err)
	}

	fmt.Printf("%d row(s) deleted.\n", rowsAffected)
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := structs.DeleteUserResponse{
		Success: true,
		Msg:     fmt.Sprintf("User with %s removed successfully", userID),
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, first_name, last_name, email, created_at FROM go_users"

	db, err := db.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute the SQL query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the users
	var users []structs.User

	// Iterate through the rows and scan the data into User structs
	for rows.Next() {
		var user structs.User
		var createdAt []uint8 // Temporary variable to store created_at as []uint8
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &createdAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse createdAt ([]uint8) into time.Time
		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for any errors during row iteration
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

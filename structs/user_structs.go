package structs

import "time"

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name" validate:"required,min=3,max=15"`
	LastName  string    `json:"last_name" validate:"required,min=3,max=15"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

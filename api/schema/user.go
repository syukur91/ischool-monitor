package schema

import (
	"time"
)

// CreateUserRequest ...
type CreateUserRequest struct {
	Nama     string `json:"nama" validate:"required"`
	Alamat   string `json:"alamat" validate:"required"`
	Password string `json:"password" validate:"required"`
	Telepon  string `json:"telepon" validate:"required"`
}

// UserResponse ...
type UserResponse struct {
	ID        int        `json:"id" db:"id"`
	Nama      string     `json:"nama" db:"nama"`
	Alamat    string     `json:"alamat" db:"alamat"`
	Password  string     `json:"password" db:"password"`
	Telepon   string     `json:"telepon" db:"telepon"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateUserRequest ...
type UpdateUserRequest struct {
	Nama     string `json:"nama"`
	Alamat   string `json:"alamat"`
	Password string `json:"password"`
	Telepon  string `json:"telepon"`
}

package schema

import (
	"time"
)

// CreateWali_KelasRequest ...
type CreateWali_KelasRequest struct {
	Nama   string `json:"nama" validate:"required"`
	Alamat string `json:"alamat" validate:"required"`
	Telpon string `json:"telpon" validate:"required"`
}

// Wali_KelasResponse ...
type Wali_KelasResponse struct {
	ID        int        `json:"id" db:"id"`
	Nama      string     `json:"nama" db:"nama"`
	Alamat    string     `json:"alamat" db:"alamat"`
	Telpon    string     `json:"telpon" db:"telpon"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateWali_KelasRequest ...
type UpdateWali_KelasRequest struct {
	Nama   string `json:"nama"`
	Alamat string `json:"alamat"`
	Telpon string `json:"telpon"`
}

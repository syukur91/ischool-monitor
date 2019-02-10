package schema

import (
	"time"
)

// CreateKelasRequest ...
type CreateKelasRequest struct {
	Nama    string `json:"nama" validate:"required"`
	Tingkat int    `json:"tingkat" validate:"required"`
}

// KelasResponse ...
type KelasResponse struct {
	ID        int        `json:"id" db:"id"`
	Nama      string     `json:"nama" db:"nama"`
	Tingkat   int        `json:"tingkat" db:"tingkat"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateKelasRequest ...
type UpdateKelasRequest struct {
	Nama    string `json:"nama"`
	Tingkat int    `json:"tingkat"`
}

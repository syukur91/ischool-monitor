package schema

import (
	"time"
)

// CreateMata_PelajaranRequest ...
type CreateMata_PelajaranRequest struct {
	Nama    string `json:"nama" validate:"required"`
	Kode    string `json:"kode" validate:"required"`
	Tingkat int    `json:"tingkat" validate:"required"`
}

// Mata_PelajaranResponse ...
type Mata_PelajaranResponse struct {
	ID        int        `json:"id" db:"id"`
	Nama      string     `json:"nama" db:"nama"`
	Kode      string     `json:"kode" db:"kode"`
	Tingkat   int        `json:"tingkat" db:"tingkat"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateMata_PelajaranRequest ...
type UpdateMata_PelajaranRequest struct {
	Nama    string `json:"nama"`
	Kode    string `json:"kode"`
	Tingkat int    `json:"tingkat"`
}

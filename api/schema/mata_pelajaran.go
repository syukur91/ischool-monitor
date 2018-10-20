package schema

import (
	"time"
)

// CreateMata_PelajaranRequest ...
type CreateMata_PelajaranRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Mata_PelajaranResponse ...
type Mata_PelajaranResponse struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateMata_PelajaranRequest ...
type UpdateMata_PelajaranRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

package legislativebody

import "errors"

// ErrLegislativeBodyNotFound is returned when a LegislativeBody is not found by ID.
var ErrLegislativeBodyNotFound = errors.New("legislative body not found")

// LegislativeBodyRepository persists and retrieves LegislativeBodies.
type LegislativeBodyRepository interface {
	GetByID(id string) (*LegislativeBody, error)
	Save(lb *LegislativeBody) error
}

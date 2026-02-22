package act

import "errors"

// ErrActNotFound is returned when an Act is not found by ID.
var ErrActNotFound = errors.New("act not found")

// ActRepository persists and retrieves Acts.
type ActRepository interface {
	GetByID(id string) (*Act, error)
	Save(a *Act) error
}

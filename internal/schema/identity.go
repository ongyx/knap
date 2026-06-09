package schema

import "github.com/google/uuid"

// Represents the identity of a user in Outline.
type Identity struct {
	// The user's UUID.
	ID uuid.UUID
	// The name of the user.
	Name string
	// The email of the user.
	Email string
}

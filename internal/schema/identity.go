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

// Applies defaults to zero values in the identity.
func (i *Identity) Defaults() *Identity {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

	if i.Name == "" {
		i.Name = "test"
	}

	if i.Email == "" {
		i.Email = "test@test.invalid"
	}

	return i
}

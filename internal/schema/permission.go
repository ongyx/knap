package schema

import (
	"encoding/json"
	"errors"
)

const (
	// Read access is enabled.
	PermissionRead Permission = iota
	// Read and write access is enabled.
	PermissionReadWrite
	// Admin operations are enabled.
	PermissionAdmin
)

var permissionToString = map[Permission]string{
	PermissionRead:      "read",
	PermissionReadWrite: "read_write",
	PermissionAdmin:     "admin",
}

var stringToPermission = make(map[string]Permission, len(permissionToString))

// interface asserts
var _ json.Marshaler = (*Permission)(nil)
var _ json.Unmarshaler = (*Permission)(nil)

// Error returned from Permission.MarshalJSON()/UnmarshalJSON() if the value is invalid.
var ErrPermissionInvalid = errors.New("sort direction is invalid")

func init() {
	for k, v := range permissionToString {
		stringToPermission[v] = k
	}
}

// The access permissions set on a collection.
type Permission int

func (p Permission) String() string {
	return permissionToString[p]
}

func (p Permission) MarshalJSON() ([]byte, error) {
	s := p.String()
	if s == "" {
		return nil, ErrPermissionInvalid
	}

	return []byte(s), nil
}

func (p *Permission) UnmarshalJSON(b []byte) error {
	var temp string
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	if v, ok := stringToPermission[temp]; ok {
		*p = v
		return nil
	}

	return ErrPermissionInvalid
}

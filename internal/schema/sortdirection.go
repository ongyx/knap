package schema

import (
	"encoding/json"
	"errors"
)

const (
	// Sort in ascending order.
	SortAscending SortDirection = iota
	// Sort in descending order.
	SortDescending
)

// interface asserts
var (
	_ json.Marshaler   = (*SortDirection)(nil)
	_ json.Unmarshaler = (*SortDirection)(nil)
)

// Error returned from SortDirection.MarshalJSON()/UnmarshalJSON() if the value is invalid.
var ErrSortDirectionInvalid = errors.New("sort direction is invalid")

// SortDirection is a direction to sort documents by in a collection.
type SortDirection int

func (sd SortDirection) String() string {
	switch sd {
	case SortAscending:
		return "asc"
	case SortDescending:
		return "desc"
	default:
		return ""
	}
}

func (sd SortDirection) MarshalJSON() ([]byte, error) {
	s := sd.String()
	if s == "" {
		return nil, ErrSortDirectionInvalid
	}

	return json.Marshal(s)
}

func (sd *SortDirection) UnmarshalJSON(b []byte) error {
	var temp string
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	switch temp {
	case "asc":
		*sd = SortAscending
	case "desc":
		*sd = SortDescending
	default:
		return ErrSortDirectionInvalid
	}

	return nil
}

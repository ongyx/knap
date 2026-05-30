package schema

import (
	"encoding/json"
	"errors"
)

const (
	// Sort in ascending order.
	DirectionAscending SortDirection = iota
	// Sort in descending order.
	DirectionDescending
)

// interface asserts
var _ json.Marshaler = (*SortDirection)(nil)
var _ json.Unmarshaler = (*SortDirection)(nil)

// Error returned from SortDirection.MarshalJSON()/UnmarshalJSON() if the value is invalid.
var ErrSortDirectionInvalid = errors.New("sort direction is invalid")

// The direction to sort documents in for a collection.
type SortDirection int

func (sd SortDirection) String() string {
	switch sd {
	case DirectionAscending:
		return "asc"
	case DirectionDescending:
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

	return []byte(s), nil
}

func (sd *SortDirection) UnmarshalJSON(b []byte) error {
	var temp string
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	switch temp {
	case "asc":
		*sd = DirectionAscending
	case "desc":
		*sd = DirectionDescending
	default:
		return ErrSortDirectionInvalid
	}

	return nil
}

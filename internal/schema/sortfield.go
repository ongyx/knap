package schema

import (
	"encoding/json"
	"errors"
)

const (
	// Sort by document title.
	SortTitle SortField = iota
	// Sort by document index.
	SortIndex
)

// interface asserts
var (
	_ json.Marshaler   = (*SortField)(nil)
	_ json.Unmarshaler = (*SortField)(nil)
)

// Error returned from SortField.MarshalJSON()/UnmarshalJSON() if the value is invalid.
var ErrSortFieldInvalid = errors.New("sort Index is invalid")

// SortField is a field to sort documents by in a collection.
type SortField int

func (sf SortField) String() string {
	switch sf {
	case SortTitle:
		return "title"
	case SortIndex:
		return "index"
	default:
		return ""
	}
}

func (sf SortField) MarshalJSON() ([]byte, error) {
	s := sf.String()
	if s == "" {
		return nil, ErrSortFieldInvalid
	}

	return json.Marshal(s)
}

func (sf *SortField) UnmarshalJSON(b []byte) error {
	var temp string
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	switch temp {
	case "title":
		*sf = SortTitle
	case "index":
		*sf = SortIndex
	default:
		return ErrSortFieldInvalid
	}

	return nil
}

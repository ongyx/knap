package schema

// Sort defines the sorting criteria for documents within a collection.
type Sort struct {
	// The field to sort with.
	Field SortField `json:"field"`
	// The direction to sort in.
	Direction SortDirection `json:"direction"`
}

// Returns the default values for a sort.
func NewSort() Sort {
	return Sort{Field: SortIndex, Direction: SortAscending}
}

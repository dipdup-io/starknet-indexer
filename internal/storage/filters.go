package storage

// IntegerFilter -
type IntegerFilter struct {
	Eq      uint64
	Neq     uint64
	Gt      uint64
	Gte     uint64
	Lt      uint64
	Lte     uint64
	Between *BetweenFilter
}

// BetweenFilter -
type BetweenFilter struct {
	From uint64
	To   uint64
}

// TimeFilter -
type TimeFilter struct {
	Gt      uint64
	Gte     uint64
	Lt      uint64
	Lte     uint64
	Between *BetweenFilter
}

// EnumFilter -
type EnumFilter struct {
	Eq    uint64
	Neq   uint64
	In    []uint64
	Notin []uint64
}

// StringFilter -
type StringFilter struct {
	Eq string
	In []string
}

// EqualityFilter -
type EqualityFilter struct {
	Eq  string
	Neq string
}

// BytesFilter -
type BytesFilter struct {
	Eq []byte
	In [][]byte
}

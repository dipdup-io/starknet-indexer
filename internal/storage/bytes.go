package storage

import (
	"encoding/hex"
	stdJSON "encoding/json"
	"errors"
	"strings"
)

// Bytes -
type Bytes stdJSON.RawMessage

// MarshalJSON returns b as the JSON encoding of b.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	return b, nil
}

// UnmarshalJSON sets *b to a copy of data.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if b == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*b = append((*b)[0:0], data...)
	return nil
}

// MustNewBytes -
func MustNewBytes(str string) Bytes {
	raw, _ := hex.DecodeString(str)
	return Bytes(raw)
}

// Implements the Unmarshaler interface of the yaml pkg.
func (b *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var seq string
	if err := unmarshal(&seq); err != nil {
		return err
	}

	seq = strings.TrimPrefix(seq, "0x")
	data, err := hex.DecodeString(seq)
	if err != nil {
		return err
	}
	*b = append((*b)[0:0], data...)
	return nil
}

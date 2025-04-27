package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Hex []byte

var nullBytes = "null"

func HexFromString(hexStr string) (Hex, error) {
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	resultBytes := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid hex string")
		}
		resultBytes[i/2] = byte(b)
	}

	return resultBytes, nil
}

func (h *Hex) UnmarshalJSON(data []byte) error {
	if h == nil {
		return nil
	}

	if nullBytes == string(data) {
		*h = nil
		return nil
	}
	length := len(data)
	if length%2 == 1 {
		return errors.Errorf("odd hex length: %d %v", length, data)
	}
	if data[0] != '"' || data[length-1] != '"' {
		return errors.Errorf("hex should be quotted string: got=%s", data)
	}

	data = bytes.Trim(data, `"`)
	*h = make(Hex, hex.DecodedLen(length-1))
	if length-1 == 0 {
		return nil
	}
	_, err := hex.Decode(*h, data)
	return err
}

func (h Hex) MarshalJSON() ([]byte, error) {
	if len(h) == 0 {
		return []byte(nullBytes), nil
	}
	hexStr := hex.EncodeToString(h)
	if len(hexStr) > 0 && hexStr[0] == '0' {
		hexStr = hexStr[1:]
	}

	return []byte(strconv.Quote("0x" + hexStr)), nil
}

func (h *Hex) Scan(src interface{}) (err error) {
	switch val := src.(type) {
	case []byte:
		*h = make(Hex, len(val))
		_ = copy(*h, val)
	case nil:
		*h = make(Hex, 0)
	default:
		return errors.Errorf("unknown hex database type: %T", src)
	}
	return nil
}

var _ driver.Valuer = (*Hex)(nil)

func (h Hex) Value() (driver.Value, error) {
	return []byte(h), nil
}

func (h Hex) Bytes() []byte {
	return []byte(h)
}

func (h Hex) String() string {
	return strings.ToUpper(hex.EncodeToString([]byte(h)))
}

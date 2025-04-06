package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

func BytesToFormattedHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hexStr := hex.EncodeToString(data)
	if len(hexStr) > 0 && hexStr[0] == '0' {
		hexStr = hexStr[1:]
	}

	return "0x" + hexStr
}

func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	hexStr = fmt.Sprintf("0%s", hexStr)

	if len(hexStr)%2 != 0 {
		return nil, errors.Errorf("hex string must have even length")
	}

	bytes := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid hex character")
		}
		bytes[i/2] = byte(b)
	}

	return bytes, nil
}

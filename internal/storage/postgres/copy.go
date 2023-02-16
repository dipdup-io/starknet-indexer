package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/shopspring/decimal"
)

func writeUint64(w *strings.Builder, i uint64) error {
	_, err := w.WriteString(fmt.Sprintf("%d", i))
	return err
}

func writeUint64Pointer(w *strings.Builder, i *uint64) (err error) {
	if i == nil {
		_, err = w.WriteString("")
	} else {
		_, err = w.WriteString(fmt.Sprintf("%d", *i))
	}
	return err
}

func writeBytes(w *strings.Builder, b []byte) error {
	_, err := w.WriteString(fmt.Sprintf("\\x%x", b))
	return err
}

func writeTime(w *strings.Builder, t time.Time) error {
	_, err := w.WriteString(fmt.Sprintf("'%s'", t.Format(time.RFC3339)))
	return err
}

func writeString(w *strings.Builder, s string) error {
	if err := w.WriteByte('"'); err != nil {
		return err
	}
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return w.WriteByte('"')
}

func writeStringArray(w *strings.Builder, arr ...string) error {
	if err := w.WriteByte('"'); err != nil {
		return err
	}
	if err := w.WriteByte('{'); err != nil {
		return err
	}
	for i := range arr {
		if i > 0 {
			if err := w.WriteByte(','); err != nil {
				return err
			}
		}
		if err := writeString(w, arr[i]); err != nil {
			return err
		}
	}
	if err := w.WriteByte('}'); err != nil {
		return err
	}
	return w.WriteByte('"')
}

func writeMap(w *strings.Builder, m map[string]any) error {
	if len(m) == 0 || m == nil {
		return nil
	}
	b, err := json.MarshalWithOption(m, json.UnorderedMap(), json.DisableNormalizeUTF8())
	if err != nil {
		return err
	}
	if _, err := w.WriteString(strconv.Quote(string(b))); err != nil {
		return err
	}
	return nil
}

func writeDecimal(w *strings.Builder, d decimal.Decimal) error {
	_, err := w.WriteString(d.String())
	return err
}

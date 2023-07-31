package filters

import (
	"context"
	"strings"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/pkg/errors"
)

func validInteger(f *pb.IntegerFilter, i uint64) bool {
	if f == nil {
		return true
	}

	switch typ := f.Filter.(type) {
	case *pb.IntegerFilter_Between:
		return i >= typ.Between.From && i < typ.Between.To
	case *pb.IntegerFilter_Eq:
		return i == typ.Eq
	case *pb.IntegerFilter_Neq:
		return i != typ.Neq
	case *pb.IntegerFilter_Gt:
		return i > typ.Gt
	case *pb.IntegerFilter_Gte:
		return i >= typ.Gte
	case *pb.IntegerFilter_Lt:
		return i < typ.Lt
	case *pb.IntegerFilter_Lte:
		return i <= typ.Lte
	}

	return false
}

func validTime(f *pb.TimeFilter, t time.Time) bool {
	if f == nil {
		return true
	}

	unixTime := uint64(t.UTC().Unix())

	switch typ := f.Filter.(type) {
	case *pb.TimeFilter_Between:
		return unixTime > typ.Between.From && unixTime < typ.Between.To
	case *pb.TimeFilter_Gt:
		return unixTime > typ.Gt
	case *pb.TimeFilter_Gte:
		return unixTime >= typ.Gte
	case *pb.TimeFilter_Lt:
		return unixTime < typ.Lt
	case *pb.TimeFilter_Lte:
		return unixTime <= typ.Lte
	}

	return false
}

func validString(f *pb.StringFilter, s string) bool {
	if f == nil {
		return true
	}

	switch typ := f.Filter.(type) {
	case *pb.StringFilter_Eq:
		return s == typ.Eq
	case *pb.StringFilter_In:
		for i := range typ.In.Arr {
			if typ.In.Arr[i] == s {
				return true
			}
		}
	}

	return false
}

func validEquality(f *pb.EqualityFilter, s string) bool {
	if f == nil {
		return true
	}

	switch typ := f.Filter.(type) {
	case *pb.EqualityFilter_Eq:
		return typ.Eq == s
	case *pb.EqualityFilter_Neq:
		return typ.Neq != s
	}

	return false
}

func validMap(f map[string]string, m map[string]any) bool {
	if f == nil {
		return true
	}
	if len(m) == 0 {
		return false
	}

	for key, value := range f {
		path := strings.Split(key, ".")
		if found, ok := findPathInMap(path, m); ok {
			return found == value
		} else {
			return false
		}
	}

	return false
}

func findPathInMap(path []string, m map[string]any) (string, bool) {
	pathLen := len(path)
	if pathLen == 0 {
		return "", false
	}
	head := path[0]
	val, ok := m[head]
	if !ok {
		return "", false
	}

	switch typ := val.(type) {
	case string:
		if pathLen > 1 {
			return "", false
		}
		return typ, true
	case map[string]any:
		if pathLen == 1 {
			return "", false
		}
		return findPathInMap(path[1:], typ)
	}

	return "", false
}

func validEnum(f *pb.EnumFilter, val uint64) bool {
	if f == nil {
		return true
	}

	switch typ := f.Filter.(type) {
	case *pb.EnumFilter_Eq:
		return typ.Eq == val
	case *pb.EnumFilter_Neq:
		return typ.Neq != val
	case *pb.EnumFilter_In:
		for i := range typ.In.Arr {
			if typ.In.Arr[i] == val {
				return true
			}
		}
	case *pb.EnumFilter_Notin:
		for i := range typ.Notin.Arr {
			if typ.Notin.Arr[i] == val {
				return false
			}
		}
	}

	return false
}

func validEnumString(f *pb.EnumStringFilter, val string) bool {
	if f == nil {
		return true
	}

	switch typ := f.Filter.(type) {
	case *pb.EnumStringFilter_Eq:
		return typ.Eq == val
	case *pb.EnumStringFilter_Neq:
		return typ.Neq != val
	case *pb.EnumStringFilter_In:
		for i := range typ.In.Arr {
			if typ.In.Arr[i] == val {
				return true
			}
		}
	case *pb.EnumStringFilter_Notin:
		for i := range typ.Notin.Arr {
			if typ.Notin.Arr[i] == val {
				return false
			}
		}
	}

	return false
}

type ids map[uint64]struct{}

// In -
func (m ids) In(i uint64) bool {
	_, ok := m[i]
	return ok
}

func fillAddressMapFromBytesFilter(ctx context.Context, address storage.IAddress, f *pb.BytesFilter, out ids) error {
	if f == nil {
		return nil
	}

	switch typ := f.Filter.(type) {
	case *pb.BytesFilter_Eq:
		a, err := address.GetByHash(ctx, typ.Eq)
		if err != nil {
			return errors.Wrapf(err, "%x", typ.Eq)
		}
		out[a.ID] = struct{}{}
	case *pb.BytesFilter_In:
		if typ.In == nil {
			return nil
		}
		for i := range typ.In.Arr {
			a, err := address.GetByHash(ctx, typ.In.Arr[i])
			if err != nil {
				return errors.Wrapf(err, "%x", typ.In.Arr[i])
			}
			out[a.ID] = struct{}{}
		}
	}
	return nil
}

func fillClassMapFromBytesFilter(ctx context.Context, class storage.IClass, f *pb.BytesFilter, out ids) error {
	if f == nil {
		return nil
	}

	switch typ := f.Filter.(type) {
	case *pb.BytesFilter_Eq:
		a, err := class.GetByHash(ctx, typ.Eq)
		if err != nil {
			return errors.Wrapf(err, "%x", typ.Eq)
		}
		out[a.ID] = struct{}{}
	case *pb.BytesFilter_In:
		if typ.In == nil {
			return nil
		}
		for i := range typ.In.Arr {
			a, err := class.GetByHash(ctx, typ.In.Arr[i])
			if err != nil {
				return errors.Wrapf(err, "%x", typ.In.Arr[i])
			}
			out[a.ID] = struct{}{}
		}
	}
	return nil
}

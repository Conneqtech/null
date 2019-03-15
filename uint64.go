package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Uint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Uint64 is not NULL
}

func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Uint64From(i uint64) Uint64 {
	return NewUint64(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Uint64FromPtr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}
	return NewUint64(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint64) ValueOrZero() uint64 {
	if !i.Valid {
		return 0
	}
	return i.Uint64
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Uint64)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			i.Uint64 = uint64(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Uint64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Uint64) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Uint64) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Uint64, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		i.Uint64 = uint64(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Uint64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint64), 10)), nil
}

func (i Uint64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint64), 10)), nil
}

func (i *Uint64) SetValid(n uint64) {
	i.Uint64 = n
	i.Valid = true
}

func (i Uint64) Ptr() *uint64 {
	if !i.Valid {
		return nil
	}
	return &i.Uint64
}

func (i Uint64) IsZero() bool {
	return !i.Valid
}

package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Int8 struct {
	Int8  int8
	Valid bool // Valid is true if Int8 is not NULL
}

func NewInt8(i int8, valid bool) Int8 {
	return Int8{
		Int8:  i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Int8From(i int8) Int8 {
	return NewInt8(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Int8FromPtr(i *int8) Int8 {
	if i == nil {
		return NewInt8(0, false)
	}
	return NewInt8(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int8) ValueOrZero() int8 {
	if !i.Valid {
		return 0
	}
	return i.Int8
}

func (i *Int8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Int8)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			i.Int8 = int8(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Int8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Int8) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Int8) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Int8, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		i.Int8 = int8(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int8), 10)), nil
}

func (i Int8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int8), 10)), nil
}

func (i *Int8) SetValid(n int8) {
	i.Int8 = n
	i.Valid = true
}

func (i Int8) Ptr() *int8 {
	if !i.Valid {
		return nil
	}
	return &i.Int8
}

func (i Int8) IsZero() bool {
	return !i.Valid
}

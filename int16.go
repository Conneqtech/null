package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Int16 struct {
	Int16 int16
	Valid bool // Valid is true if Int16 is not NULL
}

func NewInt16(i int16, valid bool) Int16 {
	return Int16{
		Int16: i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Int16From(i int16) Int16 {
	return NewInt16(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Int16FromPtr(i *int16) Int16 {
	if i == nil {
		return NewInt16(0, false)
	}
	return NewInt16(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int16) ValueOrZero() int16 {
	if !i.Valid {
		return 0
	}
	return i.Int16
}

func (i *Int16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Int16)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			i.Int16 = int16(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Int16)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Int16) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Int16) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Int16, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		i.Int16 = int16(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Int16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int16), 10)), nil
}

func (i Int16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int16), 10)), nil
}

func (i *Int16) SetValid(n int16) {
	i.Int16 = n
	i.Valid = true
}

func (i Int16) Ptr() *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

func (i Int16) IsZero() bool {
	return !i.Valid
}

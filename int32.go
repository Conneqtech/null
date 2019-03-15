package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Int32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		Int32: i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Int32From(i int32) Int32 {
	return NewInt32(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	return NewInt32(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int32) ValueOrZero() int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

func (i *Int32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Int32)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			i.Int32 = int32(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Int32)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Int32) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Int32) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Int32, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		i.Int32 = int32(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int32), 10)), nil
}

func (i Int32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Int32), 10)), nil
}

func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = true
}

func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

func (i Int32) IsZero() bool {
	return !i.Valid
}

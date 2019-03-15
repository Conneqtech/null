package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Uint8 struct {
	Uint8 uint8
	Valid bool // Valid is true if Uint8 is not NULL
}

func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		Uint8: i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	return NewUint8(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint8) ValueOrZero() uint8 {
	if !i.Valid {
		return 0
	}
	return i.Uint8
}

func (i *Uint8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Uint8)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			i.Uint8 = uint8(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Uint8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Uint8) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Uint8) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Uint8, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		i.Uint8 = uint8(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Uint8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint8), 10)), nil
}

func (i Uint8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint8), 10)), nil
}

func (i *Uint8) SetValid(n uint8) {
	i.Uint8 = n
	i.Valid = true
}

func (i Uint8) Ptr() *uint8 {
	if !i.Valid {
		return nil
	}
	return &i.Uint8
}

func (i Uint8) IsZero() bool {
	return !i.Valid
}

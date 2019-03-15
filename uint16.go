package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Uint16 struct {
	Uint16 uint16
	Valid  bool // Valid is true if Uint16 is not NULL
}

func NewUint16(i uint16, valid bool) Uint16 {
	return Uint16{
		Uint16: i,
		Valid:  valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Uint16From(i uint16) Uint16 {
	return NewUint16(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Uint16FromPtr(i *uint16) Uint16 {
	if i == nil {
		return NewUint16(0, false)
	}
	return NewUint16(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint16) ValueOrZero() uint16 {
	if !i.Valid {
		return 0
	}
	return i.Uint16
}

func (i *Uint16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Uint16)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			i.Uint16 = uint16(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Uint16)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Uint16) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Uint16) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Uint16, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		i.Uint16 = uint16(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Uint16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint16), 10)), nil
}

func (i Uint16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint16), 10)), nil
}

func (i *Uint16) SetValid(n uint16) {
	i.Uint16 = n
	i.Valid = true
}

func (i Uint16) Ptr() *uint16 {
	if !i.Valid {
		return nil
	}
	return &i.Uint16
}

func (i Uint16) IsZero() bool {
	return !i.Valid
}

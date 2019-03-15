package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"reflect"
	"strconv"
)

type Uint32 struct {
	Uint32 uint32
	Valid  bool // Valid is true if Uint32 is not NULL
}

func NewUint32(i uint32, valid bool) Uint32 {
	return Uint32{
		Uint32: i,
		Valid:  valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func Uint32From(i uint32) Uint32 {
	return NewUint32(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func Uint32FromPtr(i *uint32) Uint32 {
	if i == nil {
		return NewUint32(0, false)
	}
	return NewUint32(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint32) ValueOrZero() uint32 {
	if !i.Valid {
		return 0
	}
	return i.Uint32
}

func (i *Uint32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Uint32)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		parsedInt, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			i.Uint32 = uint32(parsedInt)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.Uint32)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

func (s *Uint32) SetBSON(raw bson.Raw) error {
	return s.UnmarshalJSON(raw.Data)
}

func (s Uint32) GetBSON() (interface{}, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Uint32, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	parsedInt, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		i.Uint32 = uint32(parsedInt)
	}
	i.Valid = err == nil
	return err
}

func (i Uint32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

func (i Uint32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

func (i *Uint32) SetValid(n uint32) {
	i.Uint32 = n
	i.Valid = true
}

func (i Uint32) Ptr() *uint32 {
	if !i.Valid {
		return nil
	}
	return &i.Uint32
}

func (i Uint32) IsZero() bool {
	return !i.Valid
}

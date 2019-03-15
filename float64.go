package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"reflect"
	"strconv"
)

// Float is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float64 struct {
	sql.NullFloat64
}

// NewFloat creates a new Float
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

// FloatFrom creates a new Float that will always be valid.
func Float64From(f float64) Float64 {
	return NewFloat64(f, true)
}

// FloatFromPtr creates a new Float that be null if f is nil.
func Float64FromPtr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}
	return NewFloat64(*f, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (f Float64) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float64.
// It also supports unmarshalling a sql.NullFloat64.
func (f *Float64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Float64 = float64(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			f.Valid = false
			return nil
		}
		f.Float64, err = strconv.ParseFloat(str, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &f.NullFloat64)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float", reflect.TypeOf(v).Name())
	}
	f.Valid = err == nil
	return err
}



// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float64 is null.
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(f.Float64),
			Str:   strconv.FormatFloat(f.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (f *Float64) SetBSON(raw bson.Raw) error {
	return f.UnmarshalJSON(raw.Data)
}

func (f Float64) GetBSON() (interface{}, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float is null.
func (f Float64) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// SetValid changes this Float's value and also sets it to be non-null.
func (f *Float64) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float's value, or a nil pointer if this Float is null.
func (f Float64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero returns true for invalid Floats, for future omitempty support (Go 1.4?)
// A non-null Float with a 0 value will not be considered zero.
func (f Float64) IsZero() bool {
	return !f.Valid
}

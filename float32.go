package null

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"reflect"
	"strconv"
)

// Float3232 is a nullable float32.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float32 struct {
	Float32 float32
	Valid   bool // Valid is true if Float32 is not NULL
}

// NewFloat3232 creates a new Float3232
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		Float32: f,
		Valid:   valid,
	}
}

// Float3232From creates a new Float3232 that will always be valid.
func Float32From(f float32) Float32 {
	return NewFloat32(f, true)
}

// Float32FromPtr creates a new Float32 that be null if f is nil.
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (f Float32) ValueOrZero() float32 {
	if !f.Valid {
		return 0
	}
	return f.Float32
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float32.
// It also supports unmarshalling a sql.NullFloat3232.
func (f *Float32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float32:
		f.Float32 = float32(x)
	case float64:
		f.Float32 = float32(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			f.Valid = false
			return nil
		}
		parsedFloat, err := strconv.ParseFloat(str, 32)
		if err != nil {
			f.Float32 = float32(parsedFloat)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &f.Float32)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float32", reflect.TypeOf(v).Name())
	}
	f.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	parsedFloat, err := strconv.ParseFloat(string(text), 32)
	if err != nil {
		f.Float32 = float32(parsedFloat)
	}
	f.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float32 is null.
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	castedFloat := float64(f.Float32)
	if math.IsInf(castedFloat, 0) || math.IsNaN(castedFloat) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(f.Float32),
			Str:   strconv.FormatFloat(castedFloat, 'g', -1, 32),
		}
	}
	return []byte(strconv.FormatFloat(castedFloat, 'f', -1, 32)), nil
}

func (t *Float32) SetBSON(raw bson.Raw) error {
	return t.UnmarshalJSON(raw.Data)
}

func (t Float32) GetBSON() (interface{}, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Float32, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float32 is null.
func (f Float32) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// SetValid changes this Float32's value and also sets it to be non-null.
func (f *Float32) SetValid(n float32) {
	f.Float32 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float32's value, or a nil pointer if this Float32 is null.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero returns true for invalid Float32s, for future omitempty support (Go 1.4?)
// A non-null Float32 with a 0 value will not be considered zero.
func (f Float32) IsZero() bool {
	return !f.Valid
}

package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Float64String is a custom type that can unmarshal both float64 and string
type Float64String float64

func (f *Float64String) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64 first
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*f = Float64String(floatVal)
		return nil
	}

	// If that fails, try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}

	// Empty string should be treated as nil
	if strVal == "" || strVal == "-" {
		return nil
	}

	// Parse string to float64
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return err
	}
	*f = Float64String(floatVal)
	return nil
}

// BoolString is a custom type that can unmarshal both bool and string
type BoolString bool

func (b *BoolString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as bool first
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolString(boolVal)
		return nil
	}

	// If that fails, try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}

	// Convert string to bool
	switch strVal {
	case "true", "True", "TRUE", "1":
		*b = BoolString(true)
	case "false", "False", "FALSE", "0", "":
		*b = BoolString(false)
	default:
		return fmt.Errorf("cannot convert string %q to bool", strVal)
	}
	return nil
}

// Int64String is a custom type that can unmarshal both int64 and string
type Int64String int64

func (i *Int64String) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int64 first
	var intVal int64
	if err := json.Unmarshal(data, &intVal); err == nil {
		*i = Int64String(intVal)
		return nil
	}

	// If that fails, try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}

	// Empty string should be treated as nil
	if strVal == "" || strVal == "-" {
		return nil
	}

	// Parse string to int64
	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64String(intVal)
	return nil
}

// Helper functions for converting custom types to standard types

// ToFloat64Ptr converts Float64String pointer to float64 pointer
func ToFloat64Ptr(f *Float64String) *float64 {
	if f == nil {
		return nil
	}
	val := float64(*f)
	return &val
}

// ToInt64Ptr converts Float64String pointer to int64 pointer
func ToInt64Ptr(f *Float64String) *int64 {
	if f == nil {
		return nil
	}
	val := int64(*f)
	return &val
}

// ToInt64PtrFromInt64String converts Int64String pointer to int64 pointer
func ToInt64PtrFromInt64String(i *Int64String) *int64 {
	if i == nil {
		return nil
	}
	val := int64(*i)
	return &val
}

// ToBool converts BoolString to bool
func ToBool(b BoolString) bool {
	return bool(b)
}

// Float64StringWithDash is a custom type that can unmarshal float64, string, and treats "-" as nil
type Float64StringWithDash struct {
	value    *float64
	isUndefined bool
}

func (f *Float64StringWithDash) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64 first
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		f.value = &floatVal
		f.isUndefined = false
		return nil
	}

	// If that fails, try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}

	// Handle special cases
	if strVal == "-" {
		f.value = nil
		f.isUndefined = true
		return nil
	}
	
	if strVal == "" {
		f.value = nil
		f.isUndefined = false
		return nil
	}

	// Parse string to float64
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return err
	}
	f.value = &floatVal
	f.isUndefined = false
	return nil
}

// ToFloat64Ptr returns the float64 pointer value
func (f *Float64StringWithDash) ToFloat64Ptr() *float64 {
	return f.value
}

// StringWithDash is a custom type that handles "-" as a special value
type StringWithDash struct {
	value string
}

func (s *StringWithDash) UnmarshalJSON(data []byte) error {
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}
	s.value = strVal
	return nil
}

// ToStringPtr returns the string pointer value
func (s *StringWithDash) ToStringPtr() *string {
	if s.value == "" {
		return nil
	}
	return &s.value
}
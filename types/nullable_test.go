package types

import (
	"encoding/json"
	"testing"
)

func TestNullableFloat64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		want             *float64
		wantUndetermined bool
		wantErr          bool
	}{
		{name: "number", input: `123.45`, want: float64Ptr(123.45)},
		{name: "numeric string", input: `"678.9"`, want: float64Ptr(678.9)},
		{name: "empty string", input: `""`, want: nil},
		{name: "dash", input: `"-"`, want: nil, wantUndetermined: true},
		{name: "null", input: `null`, want: nil},
		{name: "invalid string", input: `"abc"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n NullableFloat64
			err := json.Unmarshal([]byte(tt.input), &n)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			got := n.Ptr()
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Errorf("Ptr() = %v, want %v", got, tt.want)
			}
			if n.IsUndetermined() != tt.wantUndetermined {
				t.Errorf("IsUndetermined() = %v, want %v", n.IsUndetermined(), tt.wantUndetermined)
			}
		})
	}
}

func TestNullableInt64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *int64
	}{
		{name: "integer", input: `62000000`, want: int64Ptr(62000000)},
		{name: "decimal number", input: `123.0`, want: int64Ptr(123)},
		{name: "numeric string", input: `"456"`, want: int64Ptr(456)},
		{name: "empty string", input: `""`, want: nil},
		{name: "dash", input: `"-"`, want: nil},
		{name: "null", input: `null`, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n NullableInt64
			if err := json.Unmarshal([]byte(tt.input), &n); err != nil {
				t.Fatalf("UnmarshalJSON(%s) error = %v", tt.input, err)
			}
			got := n.Ptr()
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Errorf("Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		want             *string
		wantUndetermined bool
	}{
		{name: "string", input: `"2024-03-31"`, want: strPtr("2024-03-31")},
		{name: "empty string", input: `""`, want: nil},
		{name: "dash", input: `"-"`, want: nil, wantUndetermined: true},
		{name: "null", input: `null`, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n NullableString
			if err := json.Unmarshal([]byte(tt.input), &n); err != nil {
				t.Fatalf("UnmarshalJSON(%s) error = %v", tt.input, err)
			}
			got := n.Ptr()
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Errorf("Ptr() = %v, want %v", got, tt.want)
			}
			if n.IsUndetermined() != tt.wantUndetermined {
				t.Errorf("IsUndetermined() = %v, want %v", n.IsUndetermined(), tt.wantUndetermined)
			}
		})
	}
}

func TestNullable_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		value NullableFloat64
		want  string
	}{
		{name: "value", value: NewNullable(12.5), want: `12.5`},
		{name: "absent", value: NullableFloat64{}, want: `null`},
		{name: "undetermined", value: NewUndetermined[float64](), want: `"-"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNullable_RoundTrip(t *testing.T) {
	// unmarshal -> marshal -> unmarshal で意味が保存されること
	for _, input := range []string{`123.45`, `null`, `"-"`} {
		var n1 NullableFloat64
		if err := json.Unmarshal([]byte(input), &n1); err != nil {
			t.Fatalf("first unmarshal(%s): %v", input, err)
		}
		data, err := json.Marshal(n1)
		if err != nil {
			t.Fatalf("marshal(%s): %v", input, err)
		}
		var n2 NullableFloat64
		if err := json.Unmarshal(data, &n2); err != nil {
			t.Fatalf("second unmarshal(%s): %v", input, err)
		}
		p1, p2 := n1.Ptr(), n2.Ptr()
		if (p1 == nil) != (p2 == nil) || (p1 != nil && *p1 != *p2) || n1.IsUndetermined() != n2.IsUndetermined() {
			t.Errorf("round trip changed meaning for %s: %v/%v -> %v/%v",
				input, p1, n1.IsUndetermined(), p2, n2.IsUndetermined())
		}
	}
}

func TestNullable_Accessors(t *testing.T) {
	n := NewNullable(42.0)
	if v, ok := n.Get(); !ok || v != 42.0 {
		t.Errorf("Get() = %v, %v; want 42.0, true", v, ok)
	}
	if n.Or(0) != 42.0 {
		t.Errorf("Or(0) = %v, want 42.0", n.Or(0))
	}

	var absent NullableFloat64
	if _, ok := absent.Get(); ok {
		t.Error("Get() on absent value should return ok=false")
	}
	if absent.Or(99.9) != 99.9 {
		t.Errorf("Or(99.9) = %v, want 99.9", absent.Or(99.9))
	}

	// Ptr はコピーを返すため、書き換えても内部状態に影響しない
	p := n.Ptr()
	*p = 0
	if v, _ := n.Get(); v != 42.0 {
		t.Errorf("mutating Ptr() result changed internal value: %v", v)
	}
}

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

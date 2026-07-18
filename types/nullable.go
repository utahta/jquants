package types

import (
	"encoding/json"
	"strconv"
)

// NullableValue はNullableが保持できる値の型集合です。
type NullableValue interface {
	float64 | int64 | string
}

// Nullable はJ-Quants APIの「値が存在しない可能性のあるフィールド」を表します。
// APIはエンドポイントによって欠損を null・空文字（非設定）・"-"（未定）で表現し、
// 数値を数値型と文字列型の両方で返すため、それらすべてを単一の型で受けます。
//
//   - JSON値（数値・文字列）: 値として保持
//   - 数値の文字列（"678.9"）: 数値型ではパースして保持
//   - null / ""（非設定）: 値なし
//   - "-"（未定）: 値なし、IsUndetermined() が true
type Nullable[T NullableValue] struct {
	value        *T
	undetermined bool
}

// NewNullable は値を持つNullableを作成します。
func NewNullable[T NullableValue](v T) Nullable[T] {
	return Nullable[T]{value: &v}
}

// NewUndetermined は「未定（-）」を表すNullableを作成します。
func NewUndetermined[T NullableValue]() Nullable[T] {
	return Nullable[T]{undetermined: true}
}

// Get は値と、値が存在するかどうかを返します。
func (n Nullable[T]) Get() (T, bool) {
	if n.value == nil {
		var zero T
		return zero, false
	}
	return *n.value, true
}

// Ptr は値のコピーへのポインタを返します。値が存在しない場合はnilです。
func (n Nullable[T]) Ptr() *T {
	if n.value == nil {
		return nil
	}
	v := *n.value
	return &v
}

// Or は値を返します。値が存在しない場合はdefを返します。
func (n Nullable[T]) Or(def T) T {
	if n.value == nil {
		return def
	}
	return *n.value
}

// IsUndetermined は元の値が「未定（-）」だったかどうかを返します。
// 「非設定（空文字・null）」との区別に使用します。どちらの場合も値は存在しません。
func (n Nullable[T]) IsUndetermined() bool {
	return n.undetermined
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.value = nil
	n.undetermined = false

	if string(data) == "null" {
		return nil
	}

	// 文字列は特別扱い: "" は非設定、"-" は未定、数値型では数値文字列を受ける
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch s {
		case "":
			return nil
		case "-":
			n.undetermined = true
			return nil
		}
		v, err := valueFromString[T](s)
		if err != nil {
			return err
		}
		n.value = &v
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		// int64はJSON上小数表現（例: 123.0）で届くことがあるためfloat64経由で受ける
		if p, ok := any(&v).(*int64); ok {
			var f float64
			if err2 := json.Unmarshal(data, &f); err2 == nil {
				*p = int64(f)
				n.value = &v
				return nil
			}
		}
		return err
	}
	n.value = &v
	return nil
}

// MarshalJSON implements json.Marshaler.
// 値なしはnull、未定は"-"として出力します。
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.undetermined {
		return json.Marshal("-")
	}
	if n.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.value)
}

func valueFromString[T NullableValue](s string) (T, error) {
	var v T
	switch p := any(&v).(type) {
	case *string:
		*p = s
	case *float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return v, err
		}
		*p = f
	case *int64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return v, err
		}
		*p = int64(f)
	}
	return v, nil
}

// NullableFloat64 は欠損しうるfloat64値です。
type NullableFloat64 = Nullable[float64]

// NullableInt64 は欠損しうるint64値です。
type NullableInt64 = Nullable[int64]

// NullableString は欠損しうるstring値です。
type NullableString = Nullable[string]

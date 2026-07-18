package types

import (
	"encoding/json"
	"fmt"
	"math/big"
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
//   - null / ""（非設定） / "*"（該当なし。例: ETFの上場比）: 値なし
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

	// 文字列は特別扱い: "" は非設定、"*" は該当なし、"-" は未定、
	// 数値型では数値文字列を受ける
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch s {
		case "", "*":
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
		// int64はJSON上小数・指数表現（例: 123.0, 1.23e5）で届くことがある
		if p, ok := any(&v).(*int64); ok {
			var f float64
			if err2 := json.Unmarshal(data, &f); err2 == nil {
				// 数値であることを確認済み。値の評価はfloat64の丸めを避けて
				// 元のテキストから厳密に行う
				i, err3 := int64FromNumericText(string(data))
				if err3 != nil {
					return err3
				}
				*p = i
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
		// 数値文法の検証はfloat64パースで行い、値の評価は丸め誤差を避けるため
		// 元のテキストから厳密に行う
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return v, err
		}
		i, err := int64FromNumericText(s)
		if err != nil {
			return v, err
		}
		*p = i
	}
	return v, nil
}

// int64FromNumericText は10進の数値テキスト（整数・小数・指数表記）を
// 丸め誤差なくint64へ変換します。float64を経由すると2^53を超える整数で
// 精度が失われるため、big.Ratで元のテキストを厳密に評価します。
// 非整数値やint64で表現できない値はエラーになります。
func int64FromNumericText(s string) (int64, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return 0, fmt.Errorf("invalid numeric value %q", s)
	}
	if !r.IsInt() {
		return 0, fmt.Errorf("value %q is not an integer", s)
	}
	if !r.Num().IsInt64() {
		return 0, fmt.Errorf("value %q overflows int64", s)
	}
	return r.Num().Int64(), nil
}

// NullableFloat64 は欠損しうるfloat64値です。
type NullableFloat64 = Nullable[float64]

// NullableInt64 は欠損しうるint64値です。
type NullableInt64 = Nullable[int64]

// NullableString は欠損しうるstring値です。
type NullableString = Nullable[string]

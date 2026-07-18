package jquants

import (
	"encoding/json"
	"testing"

	"github.com/utahta/jquants/types"
)

func TestDividendResponse_UnmarshalJSON_MissingValueVariants(t *testing.T) {
	// 値あり・未定（-）・非設定（空文字）が区別して取得できること
	jsonData := `{
		"data": [
			{
				"Code": "72030",
				"StatCode": "1",
				"DivRate": "25.5",
				"CommDivRate": "",
				"SpecDivRate": "-",
				"PayDate": "2024-03-25"
			},
			{
				"Code": "99840",
				"StatCode": "1",
				"DivRate": "-",
				"CommDivRate": "",
				"SpecDivRate": "",
				"PayDate": "-"
			}
		],
		"pagination_key": ""
	}`

	var resp DividendResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("data count = %d, want 2", len(resp.Data))
	}

	// 1件目: 配当額あり、記念配当は非設定、特別配当は未定、支払日あり
	d := resp.Data[0]
	if v, ok := d.DivRate.Get(); !ok || v != 25.5 {
		t.Errorf("DivRate = %v, %v; want 25.5, true", v, ok)
	}
	if _, ok := d.CommDivRate.Get(); ok || d.CommDivRate.IsUndetermined() {
		t.Error("CommDivRate should be absent and not undetermined for empty string")
	}
	if _, ok := d.SpecDivRate.Get(); ok || !d.SpecDivRate.IsUndetermined() {
		t.Error("SpecDivRate should be absent and undetermined for dash")
	}
	if !d.HasPayableDate() || d.IsPayableDateUndecided() {
		t.Errorf("PayDate = %v; want present and decided", d.PayDate)
	}
	if d.IsDividendRateUndecided() {
		t.Error("IsDividendRateUndecided() should be false when DivRate has a value")
	}

	// 2件目: 配当額・支払日ともに未定
	d = resp.Data[1]
	if !d.IsDividendRateUndecided() {
		t.Error("IsDividendRateUndecided() should be true for dash")
	}
	if !d.IsPayableDateUndecided() || d.HasPayableDate() {
		t.Error("PayDate should be undecided for dash")
	}
	if got := d.GetTotalDividendRate(); got != nil {
		t.Errorf("GetTotalDividendRate() = %v, want nil when DivRate is undecided", *got)
	}
}

func TestDividend_GetTotalDividendRate(t *testing.T) {
	d := Dividend{
		DivRate:     types.NewNullable(30.0),
		CommDivRate: types.NewNullable(5.0),
		SpecDivRate: types.NewNullable(10.0),
	}
	if got := d.GetTotalDividendRate(); got == nil || *got != 45.0 {
		t.Errorf("GetTotalDividendRate() = %v, want 45.0", got)
	}
	if got := d.GetOrdinaryDividendRate(); got == nil || *got != 15.0 {
		t.Errorf("GetOrdinaryDividendRate() = %v, want 15.0", got)
	}

	// 記念・特別が非設定でも合計は計算できる
	d2 := Dividend{DivRate: types.NewNullable(30.0)}
	if got := d2.GetTotalDividendRate(); got == nil || *got != 30.0 {
		t.Errorf("GetTotalDividendRate() = %v, want 30.0", got)
	}
}

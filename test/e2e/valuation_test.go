//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
)

func TestValuationEndpoint(t *testing.T) {
	data, err := jq.Valuation.GetValuationsByCodeAndDate(context.Background(), "7203", testDate)
	if err != nil {
		if isSubscriptionLimited(err) {
			t.Skip("Skipping due to subscription limitation")
		}
		t.Fatal(err)
	}
	if len(data) != 1 {
		t.Fatalf("got %d records, want 1 for 7203 on %s", len(data), testDate)
	}
	v := data[0]
	if v.Code != "72030" || strings.ReplaceAll(v.Date, "-", "") != testDate {
		t.Errorf("got %s on %s, want 72030 on %s", v.Code, v.Date, testDate)
	}
	if v.EPS == nil || v.BPS == nil || v.MktCap == nil {
		t.Errorf("missing valuation data for 7203: %+v", v)
	}
}

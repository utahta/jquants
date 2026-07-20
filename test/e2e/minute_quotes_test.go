//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"
)

// TestMinuteQuotesEndpoint は/equities/bars/minuteエンドポイントのテスト
// 株価分足アドオン契約が必要。
func TestMinuteQuotesEndpoint(t *testing.T) {
	t.Run("GetMinuteQuotes_ByCodeAndDate", func(t *testing.T) {
		date := getTestDate()
		quotes, err := jq.MinuteQuotes.GetMinuteQuotesByCodeAndDate(context.Background(), "7203", date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (minute bars add-on required)")
			}
			t.Fatalf("Failed to get minute quotes: %v", err)
		}
		if len(quotes) == 0 {
			t.Skipf("No minute quotes for 7203 on %s", date)
		}

		for i, q := range quotes {
			if q.Time == "" {
				t.Errorf("Quote[%d]: Time is empty", i)
			}
			if q.O <= 0 || q.H <= 0 || q.L <= 0 || q.C <= 0 {
				t.Errorf("Quote[%d] %s: OHLC = %v/%v/%v/%v, want > 0", i, q.Time, q.O, q.H, q.L, q.C)
			}
			if q.H < q.L {
				t.Errorf("Quote[%d] %s: H (%v) < L (%v)", i, q.Time, q.H, q.L)
			}
			if q.Vo <= 0 {
				t.Errorf("Quote[%d] %s: Vo = %v, want > 0 (約定のない分足は返らない仕様)", i, q.Time, q.Vo)
			}
			if i >= 30 {
				break
			}
		}
		t.Logf("Retrieved %d minute bars for 7203 on %s", len(quotes), date)
	})
}

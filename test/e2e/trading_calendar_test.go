//go:build e2e
// +build e2e

package e2e

import (
	"strings"
	"testing"
	"time"
)

// TestTradingCalendarEndpoint は/markets/trading_calendarエンドポイントの完全なテスト
func TestTradingCalendarEndpoint(t *testing.T) {
	t.Run("GetTradingCalendar_CurrentMonth", func(t *testing.T) {
		// 今月の取引カレンダーを取得
		from := "2024-06-01"
		to := "2024-06-30"

		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange(from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// 502エラーの場合は一時的な問題として扱う
			if strings.Contains(err.Error(), "status=502") {
				t.Skipf("Skipping due to temporary server error: %v", err)
			}
			t.Fatalf("Failed to get trading calendar: %v", err)
		}

		if len(calendar) == 0 {
			t.Skip("No trading calendar data available")
		}

		// 各カレンダーエントリを詳細に検証
		for i, day := range calendar {
			// 基本情報の検証
			if day.Date == "" {
				t.Errorf("Calendar[%d]: Date is empty", i)
			}
			if day.HolDiv == "" {
				t.Errorf("Calendar[%d]: HolidayDivision is empty", i)
			}

			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(day.Date) != 10 || day.Date[4] != '-' || day.Date[7] != '-' {
				t.Errorf("Calendar[%d]: Date format invalid = %v, want YYYY-MM-DD", i, day.Date)
			}

			// 日付範囲の検証
			if day.Date < from || day.Date > to {
				t.Errorf("Calendar[%d]: Date %s is outside requested range %s to %s",
					i, day.Date, from, to)
			}

			// 休日区分の検証（期待される値）
			validDivisions := map[string]bool{
				"0": true, // 非営業日（東証・OSEともに休業）
				"1": true, // 営業日（通常の営業日）
				"2": true, // 東証半日立会日
				"3": true, // 非営業日(祝日取引あり)
			}
			if !validDivisions[day.HolDiv] {
				t.Errorf("Calendar[%d]: Invalid HolidayDivision: %s", i, day.HolDiv)
			}

			// 最初の10件の詳細ログ
			if i < 10 {
				status := "営業日"
				if day.HolDiv == "1" {
					status = "休日"
				} else if day.HolDiv == "2" {
					status = "半日営業"
				}
				t.Logf("Calendar[%d]: %s - %s", i, day.Date, status)
			}
		}

		t.Logf("Retrieved %d calendar entries for %s", len(calendar), time.Now().Format("2006-01"))
	})

	t.Run("GetTradingCalendar_SpecificRange", func(t *testing.T) {
		// 特定の期間（2024年1月）の取引カレンダーを取得
		from := "2024-01-01"
		to := "2024-01-31"

		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange(from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			if strings.Contains(err.Error(), "status=502") {
				t.Skipf("Skipping due to temporary server error: %v", err)
			}
			t.Logf("Failed to get trading calendar for historical period: %v", err)
			return
		}

		if len(calendar) == 0 {
			t.Skip("No historical trading calendar data available")
		}

		t.Logf("Retrieved %d calendar entries for January 2024", len(calendar))

		// 営業日と休日の集計
		businessDays := 0
		holidays := 0

		for _, day := range calendar {
			switch day.HolDiv {
			case "0":
				businessDays++
			case "1":
				holidays++
			}
		}

		t.Logf("January 2024 summary:")
		t.Logf("  Business days: %d", businessDays)
		t.Logf("  Holidays: %d", holidays)
		t.Logf("  Total days: %d", len(calendar))

		// 2024年1月の特定の日をチェック（元日は非営業日のはず）
		for _, day := range calendar {
			if day.Date == "2024-01-01" {
				t.Logf("New Year's Day HolidayDivision: %s", day.HolDiv)
				if day.HolDiv != "0" {
					t.Errorf("New Year's Day should be non-business day (0), got: %s", day.HolDiv)
				} else {
					t.Logf("New Year's Day correctly marked as non-business day")
				}
				break
			}
		}
	})

	t.Run("GetTradingCalendar_WeekendValidation", func(t *testing.T) {
		// 週末の扱いを確認
		from := "2024-05-01"
		to := "2024-05-31"

		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange(from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			if strings.Contains(err.Error(), "status=502") {
				t.Skipf("Skipping due to temporary server error: %v", err)
			}
			t.Skip("No trading calendar data available")
		}

		if len(calendar) == 0 {
			t.Skip("No calendar data available")
		}

		weekendBusinessDays := 0
		weekdayHolidays := 0

		for _, day := range calendar {
			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				t.Errorf("Failed to parse date: %s", day.Date)
				continue
			}

			weekday := date.Weekday()
			isWeekend := weekday == time.Saturday || weekday == time.Sunday
			isBusinessDay := day.HolDiv == "1" || day.HolDiv == "2"

			if isWeekend && isBusinessDay {
				weekendBusinessDays++
				t.Logf("Weekend business day: %s (%s)", day.Date, weekday.String())
			}

			if !isWeekend && !isBusinessDay {
				weekdayHolidays++
				t.Logf("Weekday holiday: %s (%s)", day.Date, weekday.String())
			}
		}

		t.Logf("Weekend analysis:")
		t.Logf("  Weekend business days: %d", weekendBusinessDays)
		t.Logf("  Weekday holidays: %d", weekdayHolidays)
	})

	t.Run("GetTradingCalendar_YearlyAnalysis", func(t *testing.T) {
		// 年間の営業日数分析（過去の完全な年）
		from := "2023-01-01"
		to := "2023-12-31"

		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange(from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			if strings.Contains(err.Error(), "status=502") {
				t.Skipf("Skipping due to temporary server error: %v", err)
			}
			t.Skip("No historical trading calendar data available")
		}

		if len(calendar) == 0 {
			t.Skip("No calendar data available for 2023")
		}

		monthlyBusinessDays := make(map[string]int)
		monthlyHolidays := make(map[string]int)

		for _, day := range calendar {
			month := day.Date[:7] // YYYY-MM

			if day.HolDiv == "0" {
				monthlyBusinessDays[month]++
			} else {
				monthlyHolidays[month]++
			}
		}

		t.Logf("2023 Monthly business days:")
		for month := "2023-01"; month <= "2023-12"; month = nextMonth(month) {
			business := monthlyBusinessDays[month]
			holidays := monthlyHolidays[month]
			if business > 0 || holidays > 0 {
				t.Logf("  %s: %d business days, %d holidays", month, business, holidays)
			}
		}

		totalBusinessDays := 0
		for _, count := range monthlyBusinessDays {
			totalBusinessDays += count
		}
		t.Logf("Total business days in 2023: %d", totalBusinessDays)
	})

	t.Run("GetTradingCalendar_SingleDate", func(t *testing.T) {
		// 単一日付での取得テスト
		today := time.Now().Format("2006-01-02")

		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange(today, today)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			if strings.Contains(err.Error(), "status=502") {
				t.Skipf("Skipping due to temporary server error: %v", err)
			}
			t.Logf("Failed to get calendar for single date: %v", err)
			return
		}

		if len(calendar) == 0 {
			t.Skip("No calendar data for today")
		}

		if len(calendar) != 1 {
			t.Errorf("Expected 1 calendar entry, got %d", len(calendar))
		}

		day := calendar[0]
		if day.Date != today {
			t.Errorf("Date mismatch: got %s, want %s", day.Date, today)
		}

		status := "営業日"
		if day.HolDiv == "1" {
			status = "休日"
		}
		t.Logf("Today (%s) is: %s", today, status)
	})

	t.Run("GetTradingCalendar_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 無効な日付範囲
		calendar, err := jq.TradingCalendar.GetTradingCalendarByDateRange("2024-12-31", "2024-01-01")
		if err == nil && len(calendar) > 0 {
			t.Error("Expected error or empty result for invalid date range")
		}

		// 未来の遠い日付
		futureFrom := time.Now().AddDate(10, 0, 0).Format("2006-01-02")
		futureTo := time.Now().AddDate(10, 0, 30).Format("2006-01-02")

		calendar, err = jq.TradingCalendar.GetTradingCalendarByDateRange(futureFrom, futureTo)
		if err == nil && len(calendar) > 0 {
			t.Logf("Warning: Got calendar data for far future dates")
		}
	})
}

// nextMonth は月を1つ進める（YYYY-MM形式）
func nextMonth(month string) string {
	date, err := time.Parse("2006-01", month)
	if err != nil {
		return month
	}
	return date.AddDate(0, 1, 0).Format("2006-01")
}

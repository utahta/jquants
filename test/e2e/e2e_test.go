//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/utahta/jquants"
	"github.com/utahta/jquants/client"
)

var (
	jq          *jquants.JQuantsAPI
	testDate    string   // テストで使用する直近営業日（YYYYMMDD形式）
	tradingDays []string // 直近の営業日一覧（YYYYMMDD形式、新しい順）
)

func TestMain(m *testing.M) {
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("JQUANTS_API_KEY")
	if apiKey == "" {
		panic("JQUANTS_API_KEY environment variable is not set")
	}

	// クライアント作成（v2 APIではAPIキーを直接使用）
	c := client.NewClient(apiKey)

	// JQuantsAPI作成
	jq = jquants.NewJQuantsAPI(c)

	// テストで使用する営業日を決定
	resolveTestDate()

	// テスト実行
	os.Exit(m.Run())
}

// testLagDays はテスト基準日を何日前にするかと、環境変数
// JQUANTS_TEST_LAG_DAYS で明示的に指定されたかどうかを返す。
// デフォルトはデータの公表ラグを考慮した4日。不正な値が設定されている場合は
// タイポを黙って無視しないようpanicする。
func testLagDays() (int, bool) {
	v := os.Getenv("JQUANTS_TEST_LAG_DAYS")
	if v == "" {
		return 4, false
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		panic(fmt.Sprintf("invalid JQUANTS_TEST_LAG_DAYS %q: must be a positive integer", v))
	}
	return n, true
}

// resolveTestDate はテストで使用する営業日を決定する。
// まずデフォルトのラグで営業日を選び、その日付のデータが実際に取得できるかを
// 確認する。データ遅延のあるプラン（Freeプランは12週遅延）で取得できない場合は
// 13週前のラグに自動で切り替える。JQUANTS_TEST_LAG_DAYS が明示されている場合は
// そのラグをそのまま使用する。
func resolveTestDate() {
	lag, explicit := testLagDays()
	resolveTradingDaysWithLag(lag)

	if explicit {
		return
	}
	available, err := dataAvailable(testDate)
	if err != nil {
		// プローブのエラーは遅延プランの空データとは別問題（ネットワーク・認証・
		// クエリのリグレッション等）なので、フォールバックで隠さず即失敗させる
		panic(fmt.Sprintf("e2e: failed to probe data availability for %s: %v", testDate, err))
	}
	if !available {
		const delayedPlanLagDays = 13 * 7
		fmt.Printf("e2e: data for %s is not accessible on this plan; retrying with %d-day lag\n",
			testDate, delayedPlanLagDays)
		resolveTradingDaysWithLag(delayedPlanLagDays)
	}
}

// resolveTradingDaysWithLag は取引カレンダーAPIから「lagDays日前まで」の直近の
// 営業日一覧を取得してtradingDaysとtestDateに設定する。カレンダーの取得に
// 失敗した場合は、土日を除いた直近の平日にフォールバックする。
func resolveTradingDaysWithLag(lagDays int) {
	tradingDays = nil
	to := time.Now().AddDate(0, 0, -lagDays)
	from := to.AddDate(0, 0, -21)

	days, err := jq.TradingCalendar.GetTradingDays(context.Background(),
		from.Format("20060102"), to.Format("20060102"))
	if err == nil && len(days) > 0 {
		dates := make([]string, 0, len(days))
		for _, d := range days {
			dates = append(dates, strings.ReplaceAll(d.Date, "-", ""))
		}
		sort.Sort(sort.Reverse(sort.StringSlice(dates)))
		tradingDays = dates
		testDate = tradingDays[0]
		fmt.Printf("e2e test date: %s (from trading calendar)\n", testDate)
		return
	}

	// フォールバック: 土日を除いた直近の平日（祝日は考慮できない）
	for d := to; len(tradingDays) < 10; d = d.AddDate(0, 0, -1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			tradingDays = append(tradingDays, d.Format("20060102"))
		}
	}
	testDate = tradingDays[0]
	fmt.Printf("e2e test date: %s (weekday fallback; calendar error: %v)\n", testDate, err)
}

// dataAvailable は指定日の株価データが現在のプランで取得可能かを確認する。
// 株価四本値APIは全プランで利用可能なため、プローブとして使用する。
// 有効なクエリでデータがない場合、APIはエラーではなく空データを返すため、
// 「エラーなし・空データ」だけが取得不可（データ遅延プラン）を意味する。
func dataAvailable(date string) (bool, error) {
	resp, err := jq.Quotes.GetDailyQuotes(context.Background(),
		jquants.DailyQuotesParams{Code: "72030", Date: date})
	if err != nil {
		return false, err
	}
	return resp != nil && len(resp.Data) > 0, nil
}

// findPopulatedTestDate は直近の営業日から順に最大maxDays日分countFnを呼び、
// データが存在した最初の営業日を返す。開示・公表がある日にのみデータが存在する
// イベント駆動のデータセットのテストで使用する。見つからない場合は空文字を返す。
func findPopulatedTestDate(t *testing.T, maxDays int, countFn func(date string) (int, error)) string {
	t.Helper()
	for i, date := range tradingDays {
		if i >= maxDays {
			break
		}
		n, err := countFn(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to fetch data for %s: %v", date, err)
		}
		if n > 0 {
			return date
		}
	}
	return ""
}

// isSubscriptionLimited はサブスクリプションプランによる制限があるかをチェックする
func isSubscriptionLimited(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// エラーメッセージにサブスクリプション制限の文言が含まれているかチェック
	return strings.Contains(errMsg, "This API is not available on your subscription") ||
		strings.Contains(errMsg, "not available on your subscription")
}

// getTestDate はテストで使用する直近営業日をYYYYMMDD形式で取得する
func getTestDate() string {
	return testDate
}

// getTestDateFormatted はテストで使用する直近営業日をYYYY-MM-DD形式で取得する
func getTestDateFormatted() string {
	return formatDate(testDate)
}

// formatDate はYYYYMMDD形式の日付をYYYY-MM-DD形式に変換する
func formatDate(date string) string {
	return date[:4] + "-" + date[4:6] + "-" + date[6:]
}

// getTestDateDaysAgo はテスト基準日のn日前をYYYYMMDD形式で取得する
func getTestDateDaysAgo(n int) string {
	d, err := time.Parse("20060102", testDate)
	if err != nil {
		panic(fmt.Sprintf("invalid testDate %q: %v", testDate, err))
	}
	return d.AddDate(0, 0, -n).Format("20060102")
}

// fmtPrice はnil許容の価格フィールドをログ表示用に整形する
func fmtPrice(p *float64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *p)
}

package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/utahta/jquants/client"
)

func TestValuationResponse_UnmarshalJSON(t *testing.T) {
	raw := `{"data":[
		{"Date":"2023-03-24","Code":"86970","EPS":89.4,"FwdEPS":87.9,
		 "BPS":590.65,"ROE":0.1534,"FwdROE":0.1488,"PER":22.88,
		 "FwdPER":23.26,"PBR":3.46,"MktCap":1077137.0},
		{"Date":"2023-03-27","Code":"86970","EPS":"-12.3","FwdEPS":"0",
		 "BPS":"590.65","ROE":"-0.1","FwdROE":"0","PER":null,
		 "FwdPER":"","PBR":"3.46","MktCap":"1077137.0"},
		{"Date":"2023-03-24","Code":"13060","EPS":null,"FwdEPS":null,
		 "BPS":null,"ROE":null,"FwdROE":null,"PER":null,
		 "FwdPER":null,"PBR":null,"MktCap":null},
		{"Date":"2023-03-24","Code":"99990"}
	],"pagination_key":"next"}`
	want := []Valuation{
		{
			Date: "2023-03-24", Code: "86970",
			EPS: floatPtr(89.4), FwdEPS: floatPtr(87.9), BPS: floatPtr(590.65),
			ROE: floatPtr(0.1534), FwdROE: floatPtr(0.1488), PER: floatPtr(22.88),
			FwdPER: floatPtr(23.26), PBR: floatPtr(3.46), MktCap: floatPtr(1077137),
		},
		{
			Date: "2023-03-27", Code: "86970",
			EPS: floatPtr(-12.3), FwdEPS: floatPtr(0), BPS: floatPtr(590.65),
			ROE: floatPtr(-0.1), FwdROE: floatPtr(0),
			PBR: floatPtr(3.46), MktCap: floatPtr(1077137),
		},
		{Date: "2023-03-24", Code: "13060"},
		{Date: "2023-03-24", Code: "99990"},
	}
	var resp ValuationResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp.Data, want) {
		t.Errorf("data = %+v, want %+v", resp.Data, want)
	}
	if resp.PaginationKey != "next" {
		t.Errorf("pagination key = %q, want next", resp.PaginationKey)
	}
}

func TestValuationService_GetValuations(t *testing.T) {
	tests := []struct {
		name   string
		params ValuationParams
		path   string
	}{
		{"code", ValuationParams{Code: "8697"}, "/equities/valuation?code=8697"},
		{"date", ValuationParams{Date: "2026-09-14"}, "/equities/valuation?date=2026-09-14"},
		{"code and date", ValuationParams{Code: "86970", Date: "20260914"}, "/equities/valuation?code=86970&date=20260914"},
		{"range", ValuationParams{Code: "86970", From: "20260901", To: "20260914"}, "/equities/valuation?code=86970&from=20260901&to=20260914"},
		{"pagination", ValuationParams{Date: "20260914", PaginationKey: "next+part/=&"}, "/equities/valuation?date=20260914&pagination_key=next%2Bpart%2F%3D%26"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewMockClient()
			c.SetResponse("GET", tt.path, json.RawMessage(`{"data":[{"Code":"86970","EPS":"89.4"}]}`))
			jq := NewJQuantsAPI(c)
			if jq.Valuation == nil {
				t.Fatal("Valuation service is not registered")
			}
			resp, err := jq.Valuation.GetValuations(context.Background(), tt.params)
			if err != nil {
				t.Fatal(err)
			}
			if c.LastPath != tt.path || c.LastMethod != "GET" {
				t.Errorf("request = %s %s, want GET %s", c.LastMethod, c.LastPath, tt.path)
			}
			if len(resp.Data) != 1 || !compareFloat64Ptr(resp.Data[0].EPS, floatPtr(89.4)) {
				t.Errorf("data = %+v, want EPS 89.4", resp.Data)
			}
		})
	}
}

func TestValuationService_GetValuationsRequiresCodeOrDate(t *testing.T) {
	for _, params := range []ValuationParams{{}, {From: "20260901", To: "20260914"}} {
		c := client.NewMockClient()
		_, err := NewValuationService(c).GetValuations(context.Background(), params)
		if err == nil || c.RequestCount != 0 {
			t.Errorf("params %+v: error = %v, requests = %d; want validation error without request", params, err, c.RequestCount)
		}
	}
}

func TestValuationService_Pagination(t *testing.T) {
	tests := []struct {
		name string
		path string
		next string
		get  func(context.Context, *ValuationService) ([]Valuation, error)
	}{
		{
			name: "code", path: "/equities/valuation?code=86970",
			next: "/equities/valuation?code=86970&pagination_key=next",
			get: func(ctx context.Context, s *ValuationService) ([]Valuation, error) {
				return s.GetValuationsByCode(ctx, "86970")
			},
		},
		{
			name: "date", path: "/equities/valuation?date=20260914",
			next: "/equities/valuation?date=20260914&pagination_key=next",
			get: func(ctx context.Context, s *ValuationService) ([]Valuation, error) {
				return s.GetValuationsByDate(ctx, "20260914")
			},
		},
		{
			name: "code and date", path: "/equities/valuation?code=86970&date=20260914",
			next: "/equities/valuation?code=86970&date=20260914&pagination_key=next",
			get: func(ctx context.Context, s *ValuationService) ([]Valuation, error) {
				return s.GetValuationsByCodeAndDate(ctx, "86970", "20260914")
			},
		},
		{
			name: "range", path: "/equities/valuation?code=86970&from=20260901&to=20260914",
			next: "/equities/valuation?code=86970&from=20260901&pagination_key=next&to=20260914",
			get: func(ctx context.Context, s *ValuationService) ([]Valuation, error) {
				return s.GetValuationsByCodeAndDateRange(ctx, "86970", "20260901", "20260914")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewMockClient()
			c.SetResponse("GET", tt.path, json.RawMessage(`{"data":[{"Code":"86970","EPS":89.4}],"pagination_key":"next"}`))
			c.SetResponse("GET", tt.next, json.RawMessage(`{"data":[{"Code":"72030","EPS":100}],"pagination_key":""}`))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			data, err := tt.get(ctx, NewValuationService(c))
			if err != nil {
				t.Fatal(err)
			}
			if len(data) != 2 || data[0].Code != "86970" || data[1].Code != "72030" {
				t.Errorf("data = %+v, want both pages in order", data)
			}
			if c.RequestCount != 2 || c.LastPath != tt.next {
				t.Errorf("requests = %d, last path = %s", c.RequestCount, c.LastPath)
			}
		})
	}
}

func TestValuationService_Errors(t *testing.T) {
	t.Run("API error on second page", func(t *testing.T) {
		c := client.NewMockClient()
		c.SetResponse("GET", "/equities/valuation?code=86970", json.RawMessage(`{"data":[{"Code":"86970"}],"pagination_key":"next"}`))
		apiErr := &client.APIError{StatusCode: 429, Body: `{"message":"Too Many Requests"}`}
		c.SetError("GET", "/equities/valuation?code=86970&pagination_key=next", apiErr)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		data, err := NewValuationService(c).GetValuationsByCode(ctx, "86970")
		var got *client.APIError
		if data != nil || !errors.As(err, &got) || got != apiErr {
			t.Errorf("data = %+v, error = %v; want original API error and no partial data", data, err)
		}
	})
	for _, raw := range []string{`{"data":[{"EPS":"not-a-number"}]}`, `{"data":[{"MktCap":{}}]}`} {
		t.Run(raw, func(t *testing.T) {
			c := client.NewMockClient()
			c.SetResponse("GET", "/equities/valuation?code=86970", json.RawMessage(raw))
			if _, err := NewValuationService(c).GetValuations(context.Background(), ValuationParams{Code: "86970"}); err == nil {
				t.Fatal("expected a decoding error")
			}
		})
	}
}

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"backend/internal/types"
)

func TestParseParams_ValidationErrors(t *testing.T) {
	// service はバリデーション失敗パスしか通らないため nil のままで安全
	h := &GourmetHandler{}

	testCases := []struct {
		name         string
		query        url.Values
		wantContains string
	}{
		{
			name:         "all params missing",
			query:        url.Values{},
			wantContains: "either lat/lng, address, or keyword must be provided",
		},
		{
			name: "lat only",
			query: url.Values{
				"lat": []string{"35.0"},
			},
			wantContains: "lat and lng must be provided together",
		},
		{
			name: "lng only",
			query: url.Values{
				"lng": []string{"139.0"},
			},
			wantContains: "lat and lng must be provided together",
		},
		{
			name: "lat not number",
			query: url.Values{
				"lat": []string{"abc"},
				"lng": []string{"139.0"},
			},
			wantContains: "lat must be a number",
		},
		{
			name: "lat NaN",
			query: url.Values{
				"lat": []string{"NaN"},
				"lng": []string{"139.0"},
			},
			wantContains: "lat must be a finite number",
		},
		{
			name: "lat Inf",
			query: url.Values{
				"lat": []string{"Inf"},
				"lng": []string{"139.0"},
			},
			wantContains: "lat must be a finite number",
		},
		{
			name: "lng NaN",
			query: url.Values{
				"lat": []string{"35.0"},
				"lng": []string{"NaN"},
			},
			wantContains: "lng must be a finite number",
		},
		{
			name: "range out of bounds",
			query: url.Values{
				"lat":   []string{"35.0"},
				"lng":   []string{"139.0"},
				"range": []string{"6"},
			},
			wantContains: "range must be 1, 2, 3, 4, or 5",
		},
		{
			name: "range without location",
			query: url.Values{
				"keyword": []string{"寿司"},
				"range":   []string{"1"},
			},
			wantContains: "range requires lat and lng",
		},
		{
			name: "start not integer",
			query: url.Values{
				"keyword": []string{"寿司"},
				"start":   []string{"abc"},
			},
			wantContains: "start must be an integer",
		},
		{
			name: "count not integer",
			query: url.Values{
				"keyword": []string{"寿司"},
				"count":   []string{"abc"},
			},
			wantContains: "count must be an integer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := "/"
			if encoded := tc.query.Encode(); encoded != "" {
				path += "?" + encoded
			}
			r := httptest.NewRequest(http.MethodGet, path, nil)

			_, err := h.parseParams(r)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantContains) {
				t.Fatalf("error = %q, want contains %q", err.Error(), tc.wantContains)
			}
		})
	}
}

func TestHandle_BadRequestReturnsJSON(t *testing.T) {
	// service はバリデーション失敗パスしか通らないため nil のままで安全
	h := &GourmetHandler{}

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	h.Handle(rec, r)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("content-type = %q, want contains application/json", contentType)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	errorObj, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error object in response: %#v", body)
	}
	if _, ok := errorObj["message"]; !ok {
		t.Fatalf("missing error.message in response: %#v", body)
	}
}

func TestParseParams_ValidCases(t *testing.T) {
	h := &GourmetHandler{}

	testCases := []struct {
		name  string
		query url.Values
		want  types.GourmetSearchParams
	}{
		{
			name: "lat/lng only",
			query: url.Values{
				"lat": []string{"35.6895"},
				"lng": []string{"139.6917"},
			},
			want: types.GourmetSearchParams{Lat: 35.6895, Lng: 139.6917, Start: 1, Count: 20},
		},
		{
			name: "keyword only",
			query: url.Values{
				"keyword": []string{"寿司"},
			},
			want: types.GourmetSearchParams{Keyword: "寿司", Start: 1, Count: 20},
		},
		{
			name: "lat/lng with range",
			query: url.Values{
				"lat":   []string{"35.0"},
				"lng":   []string{"139.0"},
				"range": []string{"3"},
			},
			want: types.GourmetSearchParams{Lat: 35.0, Lng: 139.0, Range: 3, Start: 1, Count: 20},
		},
		{
			name: "keyword with start and count",
			query: url.Values{
				"keyword": []string{"ラーメン"},
				"start":   []string{"11"},
				"count":   []string{"5"},
			},
			want: types.GourmetSearchParams{Keyword: "ラーメン", Start: 11, Count: 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/?"+tc.query.Encode(), nil)

			got, err := h.parseParams(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("parseParams() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

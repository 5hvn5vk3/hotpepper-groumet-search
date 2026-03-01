package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/types"
)

func TestMaskAPIKeyInLog(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "? separator",
			err:  fmt.Errorf("url?key=SECRET&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
		{
			name: "& separator",
			err:  fmt.Errorf("url?lat=35&key=SECRET"),
			want: "url?lat=35&key=[REDACTED]",
		},
		{
			name: "no key",
			err:  fmt.Errorf("network timeout"),
			want: "network timeout",
		},
		{
			name: "nil",
			err:  nil,
			want: "",
		},
		{
			name: "symbol in key",
			err:  fmt.Errorf("url?key=A+B/C&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := maskAPIKeyInLog(tc.err)
			if got != tc.want {
				t.Fatalf("maskAPIKeyInLog() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHandleServiceError_MasksAPIKeyInLog(t *testing.T) {
	var buf bytes.Buffer

	original := log.Writer()
	t.Cleanup(func() {
		log.SetOutput(original)
	})

	log.SetOutput(&buf)

	rec := httptest.NewRecorder()
	handleServiceError(rec, fmt.Errorf("url?key=SECRETKEY&lat=35.0"))

	if strings.Contains(buf.String(), "SECRETKEY") {
		t.Fatalf("log leaked raw key: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "[REDACTED]") {
		t.Fatalf("masked key not found in log: %s", buf.String())
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestHandleServiceError_HotpepperCodeMapping(t *testing.T) {
	testCases := []struct {
		name        string
		code        int
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "3000 -> 400",
			code:        3000,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "検索条件が正しくありません",
		},
		{
			name:        "1000 -> 502",
			code:        1000,
			wantStatus:  http.StatusBadGateway,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "2000 -> 500",
			code:        2000,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "unknown -> 500",
			code:        9999,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handleServiceError(rec, &types.HotpepperAPIError{Code: tc.code})

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var body map[string]map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response JSON: %v", err)
			}

			gotMessage := body["error"]["message"]
			if gotMessage != tc.wantMessage {
				t.Fatalf("message = %q, want %q", gotMessage, tc.wantMessage)
			}
		})
	}
}

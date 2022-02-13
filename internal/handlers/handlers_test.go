package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(StatusHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		target string
		method string
		want   want
	}{
		{
			name:   "positive test #1",
			target: "/update/gauge/Alloc/123.000000",
			method: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:   "Invalid Metric Type",
			target: "/update/gaugage/Alloc/123.000000",
			method: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:   "Invalid Request Method",
			target: "/update/gauge/Alloc/123.000000",
			method: http.MethodGet,
			want: want{
				code: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
		})
	}
}

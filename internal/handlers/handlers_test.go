package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			name:   "positive test #2",
			target: "/update/gauge/testGauge/100",
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
				code: 501,
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
		{
			name:   "Without ID",
			target: "/update/counter/",
			method: http.MethodPost,
			want: want{
				code: 404,
			},
		},
		{
			name:   "Without ID",
			target: "/update/gauge/",
			method: http.MethodPost,
			want: want{
				code: 404,
			},
		},
		{
			name:   "Invalid Value",
			target: "/update/counter/testCounter/none",
			method: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:   "Invalid Value",
			target: "/update/counter/testGauge/none",
			method: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name:   "Update Invalid Type",
			target: "/update/unknown/testCounter/100",
			method: http.MethodPost,
			want: want{
				code: 501,
			},
		},
	}
	app := &config.Application{}
	storage := repository.NewMapStorageRepo()
	repo := NewRepo(app, storage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(repo.UpdateMetricHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
		})
	}
}

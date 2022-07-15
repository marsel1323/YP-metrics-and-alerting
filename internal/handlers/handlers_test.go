package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/repository"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	type metric struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "Gauge valid",
			metric: metric{
				metricType:  "gauge",
				metricName:  "Alloc",
				metricValue: "123.000000",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Counter valid",
			metric: metric{
				metricType:  "counter",
				metricName:  "PollCount",
				metricValue: "123",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Counter invalid metric value",
			metric: metric{
				metricType:  "counter",
				metricName:  "PollCount",
				metricValue: "123.000000",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Gauge invalid metric value",
			metric: metric{
				metricType:  "gauge",
				metricName:  "Alloc",
				metricValue: "a",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Unknown metric type",
			metric: metric{
				metricType:  "metric",
				metricName:  "Alloc",
				metricValue: "a",
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
	}

	app := &config.Application{}
	mapStorage := repository.NewMapStorageRepo()
	repo := NewRepo(app, mapStorage)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			method := http.MethodPost
			target := "/update/{metricType}/{metricName}/{metricValue}"
			metricType := test.metric.metricType
			metricName := test.metric.metricName
			metricValue := test.metric.metricValue

			request := httptest.NewRequest(method, target, nil)
			recorder := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", metricType)
			rctx.URLParams.Add("metricName", metricName)
			rctx.URLParams.Add("metricValue", metricValue)

			request = request.WithContext(
				context.WithValue(
					request.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			handler := http.HandlerFunc(repo.UpdateMetricHandler)
			handler.ServeHTTP(recorder, request)
			result := recorder.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			defer result.Body.Close()
			//resBody, err := io.ReadAll(result.Body)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//t.Log(resBody)
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	type want struct {
		statusCode  int
		metricValue string
	}
	type metric struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "Gauge valid",
			metric: metric{
				metricType:  "gauge",
				metricName:  "Alloc",
				metricValue: "123.000000",
			},
			want: want{
				statusCode:  http.StatusOK,
				metricValue: "123.000000",
			},
		},
		{
			name: "Counter valid",
			metric: metric{
				metricType:  "counter",
				metricName:  "PollCount",
				metricValue: "123",
			},
			want: want{
				statusCode:  http.StatusOK,
				metricValue: "123",
			},
		},
	}

	app := &config.Application{}
	mapStorage := repository.NewMapStorageRepo()
	repo := NewRepo(app, mapStorage)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metricType := test.metric.metricType
			metricName := test.metric.metricName
			metricValue := test.metric.metricValue

			request := httptest.NewRequest(
				http.MethodPost,
				"/update/{metricType}/{metricName}/{metricValue}",
				nil,
			)
			recorder := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", metricType)
			rctx.URLParams.Add("metricName", metricName)
			rctx.URLParams.Add("metricValue", metricValue)

			request = request.WithContext(
				context.WithValue(
					request.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			handler := http.HandlerFunc(repo.UpdateMetricHandler)
			handler.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, test.want.statusCode, result.StatusCode)

			request = httptest.NewRequest(
				http.MethodGet,
				"/value/{metricType}/{metricName}",
				nil,
			)
			rctx = chi.NewRouteContext()
			rctx.URLParams.Add("metricType", metricType)
			rctx.URLParams.Add("metricName", metricName)

			request = request.WithContext(
				context.WithValue(
					request.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			handler = http.HandlerFunc(repo.GetMetricHandler)
			handler.ServeHTTP(recorder, request)
			result = recorder.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			t.Log(result)
		})
	}
}

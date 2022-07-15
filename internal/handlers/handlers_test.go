package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/render"
	"YP-metrics-and-alerting/internal/repository"
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_UpdateMetricHandler(t *testing.T) {
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

func Test_GetMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	type metric struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "Invalid metric type",
			metric: metric{
				metricType: "metric",
				metricName: "Alloc",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Metric not found",
			metric: metric{
				metricType: "gauge",
				metricName: "Alloc",
			},
			want: want{
				statusCode: http.StatusNotFound,
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

			recorder := httptest.NewRecorder()

			request := httptest.NewRequest(
				http.MethodGet,
				"/value/{metricType}/{metricName}",
				nil,
			)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", metricType)
			rctx.URLParams.Add("metricName", metricName)

			request = request.WithContext(
				context.WithValue(
					request.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			handler := http.HandlerFunc(repo.GetMetricHandler)
			handler.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func Test_UpdateAndGetMetricHandler(t *testing.T) {
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
			name: "Invalid metric type",
			metric: metric{
				metricType: "metric",
				metricName: "Alloc",
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
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
		//{
		//	name: "Metric not found",
		//	metric: metric{
		//		metricType:  "counter",
		//		metricName:  "metric",
		//		metricValue: "123",
		//	},
		//	want: want{
		//		statusCode: http.StatusNotFound,
		//	},
		//},
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

			handler.ServeHTTP(recorder, request)
			result = recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			t.Log(result)
		})
	}
}

// TODO: dont know how to test
func Test_GetInfoPageHandler(t *testing.T) {
	t.SkipNow()
	app := &config.Application{}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	render.NewRenderer(app)

	mapStorage := repository.NewMapStorageRepo()
	repo := NewRepo(app, mapStorage)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rctx := chi.NewRouteContext()

	request = request.WithContext(
		context.WithValue(
			request.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	handler := http.HandlerFunc(repo.GetInfoPageHandler)
	handler.ServeHTTP(recorder, request)
	result := recorder.Result()
	defer result.Body.Close()
	t.Log(result)
	assert.Equal(t, 200, result.StatusCode)
}

func Test_UpdateMetricJSONHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		metric string
		want   want
	}{
		{
			name: "Invalid metric type",
			metric: `{
				"id": "test1",
				"type": "metric",
				"value": 123.000000
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Gauge valid",
			metric: `{
				"id": "test2",
				"type": "gauge",
				"value": 123.000000
			}`,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Gauge invalid metric delta",
			metric: `{
				"id": "test3",
				"type": "gauge",
				"delta": 123
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Counter valid delta",
			metric: `{
				"id": "test4",
				"type": "counter",
				"delta": 123
			}`,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Counter invalid metric delta",
			metric: `{
				"id": "test5",
				"type": "counter",
				"delta": 123.000000
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Counter invalid metric value",
			metric: `{
				"id": "test6",
				"type": "counter",
				"value": 123.000000
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Counter invalid metric value",
			metric: `{
				"id": "test7",
				"type": "counter",
				"value": "a"
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Gauge invalid metric value",
			metric: `{
				"id": "test8",
				"type": "gauge",
				"value": "a"
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	cfg := config.InitServerConfig()
	app := &config.Application{
		Config: cfg,
	}
	mapStorage := repository.NewMapStorageRepo()
	repo := NewRepo(app, mapStorage)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			method := http.MethodPost
			target := "/update"

			body := []byte(test.metric)
			//body, err := json.Marshal(test.metric)
			//if err != nil {
			//	log.Println(err)
			//	return
			//}

			request := httptest.NewRequest(method, target, bytes.NewReader(body))
			recorder := httptest.NewRecorder()

			rctx := chi.NewRouteContext()

			request = request.WithContext(
				context.WithValue(
					request.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			handler := http.HandlerFunc(repo.UpdateMetricJSONHandler)
			handler.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

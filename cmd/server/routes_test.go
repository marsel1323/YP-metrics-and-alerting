package main

import (
	"YP-metrics-and-alerting/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	app := &config.Application{}

	r := Routes(app)
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/update/gauge/Alloc/123.000000")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "", body)
	resp.Body.Close()

	resp, body = testRequest(t, ts, "GET", "/value/gauge/Alloc")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "123.000", body)
	resp.Body.Close()

	resp, body = testRequest(t, ts, "POST", "/update/counter/testCounter/123")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "", body)
	resp.Body.Close()

	resp, body = testRequest(t, ts, "GET", "/value/counter/testCounter")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "123", body)
	resp.Body.Close()

	resp, body = testRequest(t, ts, "POST", "/update/counter/testCounter/321")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "", body)
	resp.Body.Close()

	resp, body = testRequest(t, ts, "GET", "/value/counter/testCounter")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "444", body)
	resp.Body.Close()
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

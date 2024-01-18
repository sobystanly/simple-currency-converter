package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"platform-sre-interview-excercise-master/controller"
	"testing"
)

func Test_NewExchangeRate(t *testing.T) {
	t.Run("successfully initialize rates handler", func(t *testing.T) {
		c := &mockConverter{}
		mc := &mockCache{}
		actual := NewExchangeRate(c, mc)
		assert.Equal(t, &rates{converter: c, ratesCache: mc}, actual)
	})
}

func TestRates_HealthCheck(t *testing.T) {
	t.Run("successfully return health check status ok", func(t *testing.T) {
		c := &mockConverter{}
		mc := &mockCache{}
		er := NewExchangeRate(c, mc)

		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		er.HealthCheck(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRates_Convert(t *testing.T) {
	t.Run("Given from as euro and to as euro successfully convert euro to usd", func(t *testing.T) {
		r := NewExchangeRate(&mockConverter{conversionResponse: map[string]any{
			"date": "2023-12-13",
			"usd":  0.011991018,
		}}, &mockCache{})
		reqBody := []byte(`{"from": "euro", "to": "usd"}`)
		req, err := http.NewRequest(http.MethodGet, "/convert?from=euro&to=usd", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		r.Convert(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, 0.011991018, resp["usd"])

	})

	t.Run("Given from as euro and to as euro successfully convert euro to usd from cache", func(t *testing.T) {
		r := NewExchangeRate(&mockConverter{conversionResponse: map[string]any{
			"date": "2023-12-13",
			"usd":  0.011991018,
		}}, &mockCache{value: 0.011991018, found: true})
		reqBody := []byte(`{"from": "euro", "to": "usd"}`)
		req, err := http.NewRequest(http.MethodGet, "/convert?from=euro&to=usd", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		r.Convert(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Given from as euro and to as euro failed to convert euro to usd, required query param to is missing", func(t *testing.T) {
		r := NewExchangeRate(&mockConverter{conversionResponse: map[string]any{
			"date": "2023-12-13",
			"usd":  0.011991018,
		}}, &mockCache{})
		reqBody := []byte(`{"from": "euro", "to": "usd"}`)
		req, err := http.NewRequest(http.MethodGet, "/convert?from=euro&to=", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		r.Convert(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("Given from as euro and to as euro failed to convert euro to usd, error from controller Not found", func(t *testing.T) {
		r := NewExchangeRate(&mockConverter{err: controller.ErrConversionNotFound}, &mockCache{})
		reqBody := []byte(`{"from": "euro", "to": "usd"}`)
		req, err := http.NewRequest(http.MethodGet, "/convert?from=euro&to=usd", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		r.Convert(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Given from as euro and to as euro failed to convert euro to usd, error from controller, internalServerError from conversion api", func(t *testing.T) {
		r := NewExchangeRate(&mockConverter{err: controller.ErrUnexpected}, &mockCache{})
		reqBody := []byte(`{"from": "euro", "to": "usd"}`)
		req, err := http.NewRequest(http.MethodGet, "/convert?from=euro&to=usd", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal("error creating request", err)
		}

		w := httptest.NewRecorder()

		r.Convert(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})
}

// ///////////////////////// MOCKS ////////////////////////
type mockConverter struct {
	conversionResponse map[string]any
	err                error
}

func (m mockConverter) Convert(ctx context.Context, from, to string) (map[string]any, error) {
	return m.conversionResponse, m.err
}

type mockCache struct {
	value float64
	found bool
}

func (m mockCache) Get(key string) (float64, bool) {
	return m.value, m.found
}

func (m mockCache) Add(key string, value float64) {
}

package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_NewConverter(t *testing.T) {
	t.Run("successfully initialize a converter", func(t *testing.T) {
		mc := &mockClient{}
		actual := NewConverter(mc)
		assert.Equal(t, &converter{client: mc}, actual)
	})
}

func TestConverter_Convert(t *testing.T) {
	t.Run("Given from currency as euro and to currency as usd convert euro to usd", func(t *testing.T) {
		ctx := context.Background()
		// Your mock response payload
		mockResponse := map[string]interface{}{
			"date": "2023-12-13",
			"usd":  0.011991018,
		}
		responseJSON, _ := json.Marshal(mockResponse)
		c := NewConverter(mockClient{resp: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(responseJSON)))}})
		actual, err := c.Convert(ctx, "eur", "usd")
		assert.NotNil(t, actual)
		assert.Equal(t, mockResponse, actual)
		assert.Nil(t, err)
	})

	t.Run("failed to get conversion, unexpected error from api fallback to fallback url", func(t *testing.T) {
		ctx := context.Background()
		c := NewConverter(mockClient{err: errors.New("some random error from api")})
		actual, err := c.Convert(ctx, "eur", "usd")
		assert.Nil(t, actual)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("some random error from api"), err)
	})

	t.Run("failed to get conversion, not found", func(t *testing.T) {
		ctx := context.Background()
		c := NewConverter(mockClient{resp: &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader(`{}`))}})
		actual, err := c.Convert(ctx, "eur", "usd")
		assert.Nil(t, actual)
		assert.NotNil(t, err)
		assert.Equal(t, ErrConversionNotFound, err)
	})

	t.Run("failed to get conversion, 500 from conversion api", func(t *testing.T) {
		ctx := context.Background()
		c := NewConverter(mockClient{resp: &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(`{}`))}})
		actual, err := c.Convert(ctx, "eur", "usd")
		assert.Nil(t, actual)
		assert.NotNil(t, err)
		assert.Equal(t, ErrUnexpected, err)
	})

	t.Run("failed to get conversion, error decoding response from conversion API", func(t *testing.T) {
		ctx := context.Background()
		c := NewConverter(mockClient{resp: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{`))}})
		actual, err := c.Convert(ctx, "eur", "usd")
		assert.Nil(t, actual)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("unexpected EOF"), err)
	})
}

type mockClient struct {
	resp *http.Response
	err  error
}

func (m mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

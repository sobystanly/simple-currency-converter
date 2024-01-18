package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	api         = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/%s/currencies/%s/%s.json"
	fallbackApi = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/%s/currencies/%s/%s.min.json"
)

var (
	ErrConversionNotFound = errors.New("conversion rate not found")
	ErrUnexpected         = errors.New("unexpected error from conversion api")
)

type (
	converter struct {
		client httpClient
	}
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func NewConverter(client httpClient) *converter {
	return &converter{client: client}
}

func (c *converter) Convert(ctx context.Context, from, to string) (map[string]any, error) {
	logger := logrus.WithContext(ctx)
	conversion := make(map[string]any, 0)
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	url := fmt.Sprintf(api, formattedDate, from, to)
	fallbackUrl := fmt.Sprintf(fallbackApi, formattedDate, from, to)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("error creating request for api conversion: %s", err)
		return nil, err
	}
	fallbackReq, err := http.NewRequest(http.MethodGet, fallbackUrl, nil)
	if err != nil {
		logger.Errorf("error creating request for api conversion: %s", err)
		return nil, err
	}

	logger.Debugf("resolved request for conversion: %s", url)
	response, err := c.client.Do(req)
	if err != nil {
		logger.Errorf("error from conversion api: %s", err)
		logger.Debugf("received error from api: %s, trying fallback api: %s", url, fallbackUrl)
		response, err = c.client.Do(fallbackReq)
		if err != nil {
			logger.Errorf("error from conversion fallback api:%s err: %s", fallbackApi, err)
			return nil, err
		}
	}
	logger.Debugf("received response from conversion api: %#v", response)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Errorf("unexpected status code from conversion api: %d", response.StatusCode)
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrConversionNotFound
		}
		return nil, ErrUnexpected
	}

	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&conversion); err != nil {
		logger.Errorf("error converting response from conversion api: %s", err)
		return nil, err
	}

	return conversion, nil
}

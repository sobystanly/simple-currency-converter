package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"platform-sre-interview-excercise-master/controller"
	"time"
)

const (
	from = "from"
	to   = "to"
)

type (
	rates struct {
		converter  converterI
		ratesCache ratesCache
	}
	converterI interface {
		Convert(ctx context.Context, from, to string) (map[string]any, error)
	}
	ratesCache interface {
		Get(key string) (float64, bool)
		Add(key string, value float64)
	}
	dataResponse struct {
	}
)

func NewExchangeRate(converter converterI, ratesCache ratesCache) *rates {
	return &rates{converter: converter, ratesCache: ratesCache}
}

func (er *rates) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (er *rates) Convert(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithContext(r.Context())
	requestID := extractRequestID(r)
	logger = logger.WithField("requestID", requestID)
	logger.Debugf("received a request: %#v to convert currencies, requestID: %s", r, requestID)
	fromCurrency := getQueryParam(from, r)
	toCurrency := getQueryParam(to, r)

	if fromCurrency == "" || toCurrency == "" {
		logger.Errorf("required query params from: %s or to: %s is missing", fromCurrency, toCurrency)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("from and to query parameters is required, but received from: %s, to: %s", fromCurrency, toCurrency))
		return
	}

	//check if it already exists in the cache
	value, exist := er.ratesCache.Get(fmt.Sprintf("%s-%s", fromCurrency, toCurrency))
	if exist {
		currentTime := time.Now()
		formattedDate := currentTime.Format("2006-01-02")
		conversion := map[string]interface{}{
			"date":     formattedDate,
			toCurrency: value,
		}
		logger.Debugf("successfully received conversion rate: %f from: %s, to: %s from cache", conversion[toCurrency], fromCurrency, toCurrency)
		respondWithJSON(w, http.StatusOK, conversion)
	}

	logger.Debugf("cache miss reach out to conversion API")
	conversion, err := er.converter.Convert(r.Context(), fromCurrency, toCurrency)
	if err != nil {
		logger.Errorf("error converting from: %s, to: %s", fromCurrency, toCurrency)
		if errors.Is(err, controller.ErrConversionNotFound) {
			respondWithError(w, http.StatusNotFound, map[string]string{"message": "conversion not found"})
			return
		}
		respondWithError(w, http.StatusInternalServerError, map[string]string{"message": "error processing conversion"})
		return
	}

	currentRate, ok := conversion[toCurrency].(float64)
	if ok {
		//update cache since it was a cache miss
		er.ratesCache.Add(fmt.Sprintf("%s-%s", fromCurrency, toCurrency), currentRate)
	}

	logger.Debugf("successfully received cnversion rate: %f from: %s, to: %s", conversion[toCurrency], fromCurrency, toCurrency)
	respondWithJSON(w, http.StatusOK, conversion)
}

func respondWithJSON(w http.ResponseWriter, code int, payload map[string]interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getQueryParam(key string, r *http.Request) string {
	return r.URL.Query().Get(key)
}

func extractRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		// Generate a new requestID if not provided in headers
		requestID = uuid.New().String()
	}
	return requestID
}

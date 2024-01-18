package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (er *rates) GetRoutes() http.Handler {
	//create router
	router := mux.NewRouter()

	router.Use(requestIDMiddleware)

	var routes = Routes{
		{
			Name:        "Health Check",
			Method:      http.MethodGet,
			Pattern:     "/health",
			HandlerFunc: er.HealthCheck,
		},
		{
			Name:        "GET conversion rate",
			Method:      http.MethodGet,
			Pattern:     "/convert",
			HandlerFunc: er.Convert,
		},
	}
	for _, route := range routes {
		//If need we can add middlewares here
		r := router.Methods(route.Method).Name(route.Name).HandlerFunc(route.HandlerFunc)
		r.Path(route.Pattern)
	}
	return router
}

// These things normally happen upstream in a gateway(adding tenantID, requestID etc...), adding it here for simplicity,
// Middleware to generate and add requestID to each request
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique requestID using Google's UUID library
		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		// Add the requestID as a header in the response
		w.Header().Set("X-Request-ID", requestID)

		// Call the next handler with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

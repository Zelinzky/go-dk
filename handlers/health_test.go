package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"

	"go-dk/handlers"
)

type pingerMock struct {
	err error
}

func (p *pingerMock) Ping(ctx context.Context) error {
	return p.err
}

func TestHealth(t *testing.T) {
	t.Run("returns 200", func(t *testing.T) {
		is := is.New(t)

		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{})
		requestStatusCode, _, _ := makeGetRequest(mux, "/health")
		is.Equal(http.StatusOK, requestStatusCode)
	})
	t.Run("returns 502 if the database cannot be pinged", func(t *testing.T) {
		is := is.New(t)

		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{err: errors.New("oh no")})
		code, _, _ := makeGetRequest(mux, "/health")
		is.Equal(http.StatusBadGateway, code)
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	request := httptest.NewRequest(http.MethodGet, target, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	result := response.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}

package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"

	"go-dk/handlers"
	"go-dk/model"
)

type signupperMock struct {
	email model.Email
}

func (s *signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	s.email = email
	return "123", nil
}

type senderMock struct {
	m model.Message
}

func (s *senderMock) Send(ctx context.Context, m model.Message) error {
	s.m = m
	return nil
}

func TestNewsletterSignup(t *testing.T) {
	mux := chi.NewMux()
	s := &signupperMock{}
	q := &senderMock{}
	handlers.NewsletterSignup(mux, s, q)

	t.Run("signs up a valid email address and sends a message", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFromHeader(), strings.NewReader("email=me%40example.com"))
		is.Equal(http.StatusFound, code)
		is.Equal(model.Email("me@example.com"), s.email)

		is.Equal(q.m, model.Message{
			"job":   "confirmation_email",
			"email": "me@example.com",
			"token": "123",
		})
	})

	t.Run("rejects an invalid email address", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFromHeader(), strings.NewReader("email=notanemail"))
		is.Equal(http.StatusBadRequest, code)
	})
}

// makePostResquest and returns the status code, response header and body
func makePostRequest(handler http.Handler, target string, header http.Header, body io.Reader) (int, http.Header, string) {
	request := httptest.NewRequest(http.MethodPost, target, body)
	request.Header = header
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	result := response.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}

func createFromHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return header
}

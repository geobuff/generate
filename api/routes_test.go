package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/generate/utils"
	"github.com/gorilla/mux"
)

func newTestServer(service utils.IService) *Server {
	return NewServer(":8080", 1, []string{}, []string{}, []string{}, service)
}

func TestPing(t *testing.T) {
	service := new(utils.MockService)
	server := newTestServer(service)

	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	writer := httptest.NewRecorder()
	server.ping(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %v; got %v", http.StatusOK, result.StatusCode)
	}
}

func TestCreateTrivia(t *testing.T) {
	tt := []struct {
		name               string
		createTriviaResult error
		status             int
	}{
		{
			name:               "error on service.CreateTrivia",
			createTriviaResult: errors.New("test"),
			status:             http.StatusInternalServerError,
		},
		{
			name:               "happy path",
			createTriviaResult: nil,
			status:             http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := new(utils.MockService)
			service.On("CreateTrivia").Return(tc.createTriviaResult)
			server := newTestServer(service)

			request, err := http.NewRequest("POST", "", nil)
			if err != nil {
				t.Fatal(err)
			}

			writer := httptest.NewRecorder()
			server.createTrivia(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}

func TestRegenerateTrivia(t *testing.T) {
	tt := []struct {
		name                   string
		date                   string
		regenerateTriviaResult error
		status                 int
	}{
		{
			name:                   "error on service.RegenerateTrivia",
			date:                   "2022-01-01",
			regenerateTriviaResult: errors.New("test"),
			status:                 http.StatusInternalServerError,
		},
		{
			name:                   "happy path",
			date:                   "2022-01-01",
			regenerateTriviaResult: nil,
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := new(utils.MockService)
			service.On("RegenerateTrivia", tc.date).Return(tc.regenerateTriviaResult)
			server := newTestServer(service)

			request, err := http.NewRequest("PUT", "", nil)
			if err != nil {
				t.Fatal(err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"date": tc.date,
			})

			writer := httptest.NewRecorder()
			server.regenerateTrivia(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}

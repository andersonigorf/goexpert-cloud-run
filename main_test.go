package main

import (
	"github.com/andersonigorf/goexpert-cloud-run/configs"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchCepAndWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid Path And CEP - HTTP: 200",
			url:            "http://localhost:8080/?cep=72547240",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid CEP Format - HTTP: 422 (invalid zipcode)",
			url:            "http://localhost:8080/?cep=72547",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid CEP - HTTP: 404 (can not find zipcode)",
			url:            "http://localhost:8080/?cep=72547249",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty CEP",
			url:            "http://localhost:8080/?cep=",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Path",
			url:            "http://localhost:8080/invalid/?cep=12345-678",
			expectedStatus: http.StatusNotFound,
		},
	}

	config := &configs.Config{
		WeatherApiKey: "8b8a8316aa174cf681c175307240306",
	}
	handler := SearchCepAndWeatherHandler(config)

	router := mux.NewRouter()
	router.HandleFunc("/", handler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

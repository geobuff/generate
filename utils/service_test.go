package utils

import (
	"testing"

	"github.com/geobuff/generate/storage"
)

func TestCreateTrivia(t *testing.T) {
	tt := []struct {
		name     string
		expected string
	}{
		{
			name:     "happy path",
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			store := storage.NewMockStore()
			service := NewService(store)

			err := service.CreateTrivia()

			if tc.expected != "" && err != nil && err.Error() != tc.expected {
				t.Error(err)
			}
		})
	}
}

func BenchmarkCreateTrivia(b *testing.B) {
	store := storage.NewMockStore()
	service := NewService(store)

	for n := 0; n < b.N; n++ {
		service.CreateTrivia()
	}
}

func TestRegenerateTrivia(t *testing.T) {
	tt := []struct {
		name     string
		date     string
		expected string
	}{
		{
			name:     "invalid date",
			date:     "",
			expected: "parsing time \"\" as \"2006-02-01\": cannot parse \"\" as \"2006\"",
		},
		{
			name:     "happy path",
			date:     "2022-01-01",
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			store := storage.NewMockStore()
			service := NewService(store)

			err := service.RegenerateTrivia(tc.date)

			if tc.expected != "" && err != nil && err.Error() != tc.expected {
				t.Error(err)
			}
		})
	}
}

func BenchmarkRegenerateTrivia(b *testing.B) {
	store := storage.NewMockStore()
	service := NewService(store)

	for n := 0; n < b.N; n++ {
		service.RegenerateTrivia("2022-01-01")
	}
}

package utils

import (
	"testing"

	"github.com/geobuff/generate/storage"
)

func TestCreateTrivia(t *testing.T) {
	store := storage.NewMockStore()
	service := NewService(store)
	err := service.CreateTrivia()
	if err != nil {
		t.Error(err)
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
	store := storage.NewMockStore()
	service := NewService(store)
	err := service.RegenerateTrivia("2022-01-01")
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkRegenerateTrivia(b *testing.B) {
	store := storage.NewMockStore()
	service := NewService(store)

	for n := 0; n < b.N; n++ {
		service.RegenerateTrivia("2022-01-01")
	}
}

package utils

import (
	"github.com/geobuff/generate/types"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateTrivia() (*types.TriviaDto, error) {
	args := m.Called()
	return args.Get(0).(*types.TriviaDto), args.Error(1)
}

func (m *MockService) RegenerateTrivia(dateString string) (*types.TriviaDto, error) {
	args := m.Called(dateString)
	return args.Get(0).(*types.TriviaDto), args.Error(1)
}

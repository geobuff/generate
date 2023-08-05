package utils

import "github.com/stretchr/testify/mock"

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateTrivia() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) RegenerateTrivia(dateString string) error {
	args := m.Called(dateString)
	return args.Error(0)
}

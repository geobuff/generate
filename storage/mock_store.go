package storage

import (
	"time"

	"github.com/geobuff/generate/types"
)

type MockStore struct{}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s *MockStore) ClearTriviaPlayTriviaId(triviaId int) error {
	return nil
}

func (s *MockStore) DeleteTriviaAnswers(triviaQuestionId int) error {
	return nil
}

func (s *MockStore) DeleteTrivia(trivia *types.TriviaDto) error {
	return nil
}

func (s *MockStore) GetTrivia(date string) (*types.TriviaDto, error) {
	return &types.TriviaDto{}, nil
}

func (s *MockStore) GetMap(className string) (types.MapDto, error) {
	return types.MapDto{}, nil
}

func (s *MockStore) SetTriviaMaxScore(triviaID, maxScore int) error {
	return nil
}

func (s *MockStore) GetMappingEntries(key string) ([]types.MappingEntryDto, error) {
	if key == "us-states" {
		return states, nil
	}

	if key == "world-capitals" {
		return capitals, nil
	}

	return countries, nil
}

func (s *MockStore) GetTodaysManualTriviaQuestions() ([]types.ManualTriviaQuestion, error) {
	return []types.ManualTriviaQuestion{}, nil
}

func (s *MockStore) GetTriviaQuestionCategories(onlyActive bool) ([]types.TriviaQuestionCategory, error) {
	return []types.TriviaQuestionCategory{}, nil
}

func (s *MockStore) CreateTriviaQuestion(question types.TriviaQuestion) (int, error) {
	return 0, nil
}

func (s *MockStore) CreateTriviaAnswer(answer types.TriviaAnswer) error {
	return nil
}

func (s *MockStore) GetManualTriviaQuestions(typeID int, lastUsedMax string, allowedCategories []int) ([]types.ManualTriviaQuestion, error) {
	return []types.ManualTriviaQuestion{}, nil
}

func (s *MockStore) GetManualTriviaAnswers(questionID int) ([]types.ManualTriviaAnswer, error) {
	return []types.ManualTriviaAnswer{}, nil
}

func (s *MockStore) UpdateManualTriviaQuestionLastUsed(questionID int) error {
	return nil
}

func (s *MockStore) TriviaDoesNotExistForDate(date time.Time) (bool, error) {
	return true, nil
}

func (s *MockStore) CreateTrivia(name string, date time.Time) (int, error) {
	return 0, nil
}

package storage

import (
	"database/sql"

	"github.com/geobuff/generate/types"
)

type IStore interface {
	GetConnection() *sql.DB
	ClearTriviaPlayTriviaId(triviaId int) error
	DeleteTriviaAnswers(triviaQuestionId int) error
	GetTrivia(date string) (*types.TriviaDto, error)
	DeleteTrivia(trivia *types.TriviaDto) error
	SetTriviaMaxScore(triviaID, maxScore int) error
	GetMappingEntries(key string) ([]types.MappingEntryDto, error)
	GetTodaysManualTriviaQuestions() ([]types.ManualTriviaQuestion, error)
	GetTriviaQuestionCategories(onlyActive bool) ([]types.TriviaQuestionCategory, error)
	GetMap(className string) (types.MapDto, error)
	CreateTriviaQuestion(question types.TriviaQuestion) (int, error)
	CreateTriviaAnswer(answer types.TriviaAnswer) error
	GetManualTriviaQuestions(typeID int, lastUsedMax string, allowedCategories []int) ([]types.ManualTriviaQuestion, error)
	GetManualTriviaAnswers(questionID int) ([]types.ManualTriviaAnswer, error)
	UpdateManualTriviaQuestionLastUsed(questionID int) error
}

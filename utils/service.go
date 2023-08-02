package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/geobuff/generate/storage"
	"github.com/geobuff/generate/types"
)

type IService interface {
	CreateTrivia() error
	RegenerateTrivia(dateString string) error
}

type Service struct {
	store storage.IStore
}

func NewService(store storage.IStore) *Service {
	return &Service{
		store,
	}
}

func (s *Service) CreateTrivia() error {
	date := time.Now().AddDate(0, 0, 1)
	return s.createTriviaForDate(date)
}

func (s *Service) RegenerateTrivia(dateString string) error {
	trivia, err := s.store.GetTrivia(dateString)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != sql.ErrNoRows {
		if err = s.store.DeleteTrivia(trivia); err != nil {
			return err
		}
	}

	date, err := time.Parse("2006-01-02", dateString)
	return s.createTriviaForDate(date)
}

func (s *Service) createTriviaForDate(date time.Time) error {
	var id int
	if err := s.store.GetConnection().QueryRow("SELECT id FROM trivia WHERE date = $1", date).Scan(&id); err != sql.ErrNoRows {
		return errors.New("trivia for current date already created")
	}

	_, month, day := date.Date()
	weekday := date.Weekday().String()
	statement := "INSERT INTO trivia (name, date, maxscore) VALUES ($1, $2, $3) RETURNING id;"
	if err := s.store.GetConnection().QueryRow(statement, fmt.Sprintf("%s, %s %d", weekday, month, day), date, 0).Scan(&id); err != nil {
		return err
	}

	count, err := s.generateQuestions(id, 10)
	if err != nil {
		return err
	}

	return s.store.SetTriviaMaxScore(id, count)
}

func (s *Service) generateQuestions(triviaId, max int) (int, error) {
	countries, err := s.store.GetMappingEntries("world-countries")
	if err != nil {
		return 0, err
	}

	capitals, err := s.store.GetMappingEntries("world-capitals")
	if err != nil {
		return 0, err
	}

	states, err := s.store.GetMappingEntries("us-states")
	if err != nil {
		return 0, err
	}

	count := 0
	err = s.whatCountry(triviaId, countries)
	if err != nil {
		return count, err
	}

	err = s.whatCapital(triviaId, countries, capitals)
	if err != nil {
		return count, err
	}

	err = s.whatUSState(triviaId, states)
	if err != nil {
		return count, err
	}

	err = s.whatFlag(triviaId, countries)
	if err != nil {
		return count, err
	}
	count = count + 4

	questions, err := s.store.GetTodaysManualTriviaQuestions()
	if err != nil && err != sql.ErrNoRows {
		return count, err
	}

	todaysQuestionCount, err := s.createQuestionsAndAnswers(questions, triviaId, max-count)
	if err != nil {
		return count, err
	}

	count = count + todaysQuestionCount
	if count < max {
		categories, err := s.store.GetTriviaQuestionCategories(true)
		if err != nil {
			return count, err
		}

		remainder := max - count
		var allowedCategories []types.TriviaQuestionCategory
		for i := 0; i < remainder; i++ {
			index := rand.Intn(len(categories))
			allowedCategories = append(allowedCategories, categories[index])
			categories = append(categories[:index], categories[index+1:]...)
		}

		maxTextCount := remainder / 2
		var textCategories []int
		var imageCategories []int
		for i := 0; i < len(allowedCategories); i++ {
			curr := allowedCategories[i]
			if !curr.ImageOnly && len(textCategories) < maxTextCount {
				textCategories = append(textCategories, curr.ID)
			} else {
				imageCategories = append(imageCategories, curr.ID)
			}
		}

		textQuestionCount, err := s.setRandomManualTriviaQuestions(triviaId, types.QUESTION_TYPE_TEXT, maxTextCount, textCategories)
		if err != nil && err != sql.ErrNoRows {
			return count, err
		}
		count = count + textQuestionCount

		imageQuestionCount, err := s.setRandomManualTriviaQuestions(triviaId, types.QUESTION_TYPE_IMAGE, max-count, imageCategories)
		if err != nil && err != sql.ErrNoRows {
			return count, err
		}
		count = count + imageQuestionCount
	}

	return count, nil
}

func (s *Service) whatCountry(triviaId int, countries []types.MappingEntryDto) error {
	max := len(types.TopLandmass)
	index := rand.Intn(max)
	country := types.TopLandmass[index]

	question := types.TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      types.QUESTION_TYPE_MAP,
		Question:    "Which country is highlighted above?",
		Map:         "WorldCountries",
		Highlighted: country,
	}

	questionId, err := s.store.CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := types.TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country,
		IsCorrect:        true,
	}

	err = s.store.CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i, val := range countries {
		if val.SVGName == country {
			index = i
			break
		}
	}
	countries = append(countries[:index], countries[index+1:]...)
	max = len(countries)

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		country = countries[index].SVGName
		answer := types.TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             country,
			IsCorrect:        false,
		}

		err = s.store.CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func (s *Service) whatCapital(triviaId int, countries []types.MappingEntryDto, capitals []types.MappingEntryDto) error {
	max := len(types.TopLandmass)
	index := rand.Intn(max)
	country := types.TopLandmass[index]
	var code string
	for _, val := range countries {
		if val.SVGName == country {
			code = val.Code
			break
		}
	}

	var capitalName string
	for _, value := range capitals {
		if value.Code == code {
			capitalName = value.SVGName
		}
	}

	question := types.TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      types.QUESTION_TYPE_MAP,
		Question:    fmt.Sprintf("What is the capital city of %s?", country),
		Map:         "WorldCapitals",
		Highlighted: capitalName,
	}

	questionId, err := s.store.CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := types.TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             capitalName,
		IsCorrect:        true,
	}
	err = s.store.CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i, val := range capitals {
		if val.SVGName == capitalName {
			index = i
			break
		}
	}

	capitals = append(capitals[:index], capitals[index+1:]...)
	max = len(capitals)
	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		capital := capitals[index]
		answer := types.TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             capital.SVGName,
			IsCorrect:        false,
		}

		err = s.store.CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		capitals = append(capitals[:index], capitals[index+1:]...)
		max = max - 1
	}

	return nil
}

func (s *Service) whatUSState(triviaId int, states []types.MappingEntryDto) error {
	max := len(states)
	index := rand.Intn(max)
	state := states[index]
	states = append(states[:index], states[index+1:]...)
	max = max - 1

	question := types.TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      types.QUESTION_TYPE_MAP,
		Question:    "Which US state is highlighted above?",
		Map:         "UsStates",
		Highlighted: state.SVGName,
	}

	questionId, err := s.store.CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := types.TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             state.SVGName,
		IsCorrect:        true,
	}

	err = s.store.CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		state = states[index]
		answer := types.TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             state.SVGName,
			IsCorrect:        false,
		}

		err = s.store.CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		states = append(states[:index], states[index+1:]...)
		max = max - 1
	}

	return nil
}

func (s *Service) whatFlag(triviaId int, countries []types.MappingEntryDto) error {
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]
	countries = append(countries[:index], countries[index+1:]...)
	max = max - 1

	question := types.TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   types.QUESTION_TYPE_FLAG,
		Question: "Which country has this flag?",
		FlagCode: country.Code,
	}

	questionId, err := s.store.CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := types.TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country.SVGName,
		IsCorrect:        true,
	}
	err = s.store.CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		country = countries[index]
		answer := types.TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             country.SVGName,
			IsCorrect:        false,
		}

		err = s.store.CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func (s *Service) setRandomManualTriviaQuestions(triviaID, typeID, quantity int, allowedCategories []int) (int, error) {
	lastUsedMax := time.Now().AddDate(0, 0, -7)
	questions, err := s.store.GetManualTriviaQuestions(typeID, lastUsedMax.Format("2006-01-02"), allowedCategories)
	if err != nil {
		return 0, err
	}

	return s.createQuestionsAndAnswers(questions, triviaID, quantity)
}

func (s *Service) createQuestionsAndAnswers(questions []types.ManualTriviaQuestion, triviaID, quantity int) (int, error) {
	count := 0
	for i := 0; i < quantity; i++ {
		if len(questions) == 0 {
			break
		}

		max := len(questions)
		index := rand.Intn(max)
		manualQuestion := questions[index]

		question := types.TriviaQuestion{
			TriviaId:           triviaID,
			TypeID:             manualQuestion.TypeID,
			Question:           manualQuestion.Question,
			Explainer:          manualQuestion.Explainer,
			Map:                manualQuestion.Map,
			Highlighted:        manualQuestion.Highlighted,
			FlagCode:           manualQuestion.FlagCode,
			ImageURL:           manualQuestion.ImageURL,
			ImageAttributeName: manualQuestion.ImageAttributeName,
			ImageAttributeURL:  manualQuestion.ImageAttributeURL,
			ImageWidth:         manualQuestion.ImageWidth,
			ImageHeight:        manualQuestion.ImageHeight,
			ImageAlt:           manualQuestion.ImageAlt,
		}

		questionID, err := s.store.CreateTriviaQuestion(question)
		if err != nil {
			return count, err
		}

		answers, err := s.store.GetManualTriviaAnswers(manualQuestion.ID)
		if err != nil {
			return count, err
		}

		for _, answer := range answers {
			newAnswer := types.TriviaAnswer{
				TriviaQuestionID: questionID,
				Text:             answer.Text,
				IsCorrect:        answer.IsCorrect,
				FlagCode:         answer.FlagCode,
			}

			if err := s.store.CreateTriviaAnswer(newAnswer); err != nil {
				return count, err
			}
		}

		questions = append(questions[:index], questions[index+1:]...)
		if err := s.store.UpdateManualTriviaQuestionLastUsed(manualQuestion.ID); err != nil {
			return count, err
		}
		count = count + 1
	}

	return count, nil
}

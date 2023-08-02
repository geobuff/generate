package storage

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/geobuff/generate/types"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	Connection *sql.DB
}

func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	connection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		Connection: connection,
	}, err
}

func (s *PostgresStore) GetConnection() *sql.DB {
	return s.Connection
}

func (s *PostgresStore) ClearTriviaPlayTriviaId(triviaId int) error {
	var id int
	statement := "UPDATE triviaplays set triviaid = null WHERE triviaid = $1 RETURNING id;"
	return s.Connection.QueryRow(statement, triviaId).Scan(&id)
}

func (s *PostgresStore) DeleteTriviaAnswers(triviaQuestionId int) error {
	statement := "DELETE FROM triviaAnswers WHERE triviaQuestionId = $1 RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, triviaQuestionId).Scan(&id)
}

func (s *PostgresStore) DeleteTrivia(trivia *types.TriviaDto) error {
	if err := s.clearTriviaPlayTriviaId(trivia.ID); err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, question := range trivia.Questions {
		if err := s.deleteTriviaAnswers(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}

		if err := s.deleteTriviaQuestion(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	var id int
	return s.Connection.QueryRow("DELETE FROM trivia WHERE id = $1 RETURNING id;", trivia.ID).Scan(&id)
}

func (s *PostgresStore) clearTriviaPlayTriviaId(triviaId int) error {
	var id int
	statement := "UPDATE triviaplays set triviaid = null WHERE triviaid = $1 RETURNING id;"
	return s.Connection.QueryRow(statement, triviaId).Scan(&id)
}

func (s *PostgresStore) deleteTriviaAnswers(triviaQuestionId int) error {
	statement := "DELETE FROM triviaAnswers WHERE triviaQuestionId = $1 RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, triviaQuestionId).Scan(&id)
}

func (s *PostgresStore) deleteTriviaQuestion(questionId int) error {
	statement := "DELETE FROM triviaQuestions WHERE id = $1 RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, questionId).Scan(&id)
}

func (s *PostgresStore) GetTrivia(date string) (*types.TriviaDto, error) {
	var result types.TriviaDto
	err := s.Connection.QueryRow("SELECT id, name, maxscore from trivia WHERE date = $1;", date).Scan(&result.ID, &result.Name, &result.MaxScore)
	if err != nil {
		return nil, err
	}

	questions, err := s.getTriviaQuestions(result.ID)
	if err != nil {
		return nil, err
	}

	result.Questions = questions
	return &result, nil
}

func (s *PostgresStore) getTriviaQuestions(triviaId int) ([]types.QuestionDto, error) {
	rows, err := s.Connection.Query("SELECT q.id, t.name, q.question, q.map, q.highlighted, q.flagCode, f.url, q.imageUrl, q.imageAttributeName, q.imageAttributeUrl, q.imageWidth, q.imageHeight, q.imageAlt, q.explainer FROM triviaQuestions q JOIN triviaQuestionType t ON t.id = q.typeId LEFT JOIN flagEntries f ON f.code = q.flagCode WHERE q.triviaId = $1;", triviaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []types.QuestionDto{}
	for rows.Next() {
		var question types.QuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question, &question.MapName, &question.Highlighted, &question.FlagCode, &question.FlagUrl, &question.ImageURL, &question.ImageAttributeName, &question.ImageAttributeURL, &question.ImageWidth, &question.ImageHeight, &question.ImageAlt, &question.Explainer); err != nil {
			return nil, err
		}

		if question.MapName != "" {
			svgMap, err := s.GetMap(question.MapName)
			if err != nil {
				return nil, err
			}
			question.Map = svgMap
		}

		answers, err := s.getTriviaAnswers(question.ID)
		if err != nil {
			return nil, err
		}

		question.Answers = answers
		questions = append(questions, question)
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions, nil
}

func (s *PostgresStore) getTriviaAnswers(triviaQuestionId int) ([]types.AnswerDto, error) {
	rows, err := s.Connection.Query("SELECT a.text, a.isCorrect, a.flagCode, f.url FROM triviaAnswers a LEFT JOIN flagentries f ON f.code = a.flagcode WHERE triviaQuestionId = $1;", triviaQuestionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers = []types.AnswerDto{}
	for rows.Next() {
		var answer types.AnswerDto
		if err = rows.Scan(&answer.Text, &answer.IsCorrect, &answer.FlagCode, &answer.FlagUrl); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	if len(answers) > 2 {
		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})
	}

	return answers, nil
}

func (s *PostgresStore) GetMap(className string) (types.MapDto, error) {
	statement := "SELECT * from maps WHERE classname = $1;"
	var m types.MapDto
	err := s.Connection.QueryRow(statement, className).Scan(&m.ID, &m.Key, &m.ClassName, &m.Label, &m.ViewBox)
	if err != nil {
		return types.MapDto{}, err
	}

	elements, err := s.getMapElements(m.ID)
	if err != nil {
		return types.MapDto{}, err
	}

	m.Elements = elements
	return m, nil
}

func (s *PostgresStore) getMapElements(mapId int) ([]types.MapElementDto, error) {
	rows, err := s.Connection.Query("SELECT e.id, e.mapid, t.name, e.elementid, e.name, e.d, e.points, e.x, e.y, e.width, e.height, e.cx, e.cy, e.r, e.transform, e.xlinkhref, e.clippath, e.clippathid, e.x1, e.y1, e.x2, e.y2 FROM mapElements e JOIN mapElementType t ON t.id = e.typeid WHERE e.mapId = $1;", mapId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements = []types.MapElementDto{}
	for rows.Next() {
		var e types.MapElementDto
		if err = rows.Scan(&e.EntryID, &e.MapID, &e.Type, &e.ID, &e.Name, &e.D, &e.Points, &e.X, &e.Y, &e.Width, &e.Height, &e.Cx, &e.Cy, &e.R, &e.Transform, &e.XlinkHref, &e.ClipPath, &e.ClipPathId, &e.X1, &e.Y1, &e.X2, &e.Y2); err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}

	return elements, rows.Err()
}

func (s *PostgresStore) SetTriviaMaxScore(triviaID, maxScore int) error {
	statement := "UPDATE trivia SET maxScore = $1 WHERE id = $2 RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, maxScore, triviaID).Scan(&id)
}

func (s *PostgresStore) GetMappingEntries(key string) ([]types.MappingEntryDto, error) {
	rows, err := s.Connection.Query("SELECT m.id, m.groupid, m.name, m.code, COALESCE(f.url, ''), m.svgname, lower(m.alternativenames::text)::text[], lower(m.prefixes::text)::text[], m.grouping from mappingEntries m JOIN mappingGroups g ON g.id = m.groupId LEFT JOIN flagEntries f ON f.code = m.code WHERE g.key = $1;", key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []types.MappingEntryDto{}
	for rows.Next() {
		var entry types.MappingEntryDto
		if err = rows.Scan(&entry.ID, &entry.GroupID, &entry.Name, &entry.Code, &entry.FlagUrl, &entry.SVGName, &entry.AlternativeNames, &entry.Prefixes, &entry.Grouping); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *PostgresStore) GetTodaysManualTriviaQuestions() ([]types.ManualTriviaQuestion, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := s.Connection.Query("SELECT * FROM manualtriviaquestions WHERE quizDate = $1;", today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []types.ManualTriviaQuestion{}
	for rows.Next() {
		var question types.ManualTriviaQuestion
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL, &question.LastUsed, &question.QuizDate, &question.Explainer, &question.LastUpdated, &question.CategoryID, &question.ImageAttributeName, &question.ImageAttributeURL, &question.ImageWidth, &question.ImageHeight, &question.ImageAlt); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func (s *PostgresStore) GetTriviaQuestionCategories(onlyActive bool) ([]types.TriviaQuestionCategory, error) {
	statement := "SELECT * from triviaquestioncategory"
	if onlyActive {
		statement += " WHERE isactive"
	}
	statement += ";"

	rows, err := s.Connection.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories = []types.TriviaQuestionCategory{}
	for rows.Next() {
		var category types.TriviaQuestionCategory
		if err = rows.Scan(&category.ID, &category.Name, &category.IsActive, &category.ImageOnly); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (s *PostgresStore) CreateTriviaQuestion(question types.TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl, imageAttributeName, imageAttributeUrl, imageWidth, imageHeight, imageAlt, explainer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;"
	var id int
	err := s.Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.ImageAttributeName, question.ImageAttributeURL, question.ImageWidth, question.ImageHeight, question.ImageAlt, question.Explainer).Scan(&id)
	return id, err
}

func (s *PostgresStore) CreateTriviaAnswer(answer types.TriviaAnswer) error {
	statement := "INSERT INTO triviaAnswers (triviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, answer.TriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func convertCategories(categories []int) []string {
	var result []string
	for _, val := range categories {
		result = append(result, strconv.Itoa(val))
	}
	return result
}

func (s *PostgresStore) GetManualTriviaQuestions(typeID int, lastUsedMax string, allowedCategories []int) ([]types.ManualTriviaQuestion, error) {
	statement := "SELECT DISTINCT ON (categoryid) * FROM manualtriviaquestions WHERE typeid = $1 AND quizdate IS null AND (lastUsed IS null OR lastUsed < $2) AND categoryid = ANY($3);"
	rows, err := s.Connection.Query(statement, typeID, lastUsedMax, pq.Array(convertCategories(allowedCategories)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []types.ManualTriviaQuestion{}
	for rows.Next() {
		var question types.ManualTriviaQuestion
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL, &question.LastUsed, &question.QuizDate, &question.Explainer, &question.LastUpdated, &question.CategoryID, &question.ImageAttributeName, &question.ImageAttributeURL, &question.ImageWidth, &question.ImageHeight, &question.ImageAlt); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func (s *PostgresStore) GetManualTriviaAnswers(questionID int) ([]types.ManualTriviaAnswer, error) {
	rows, err := s.Connection.Query("SELECT * FROM manualtriviaanswers WHERE manualtriviaquestionid = $1;", questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers = []types.ManualTriviaAnswer{}
	for rows.Next() {
		var answer types.ManualTriviaAnswer
		if err = rows.Scan(&answer.ID, &answer.ManualTriviaQuestionID, &answer.Text, &answer.IsCorrect, &answer.FlagCode); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, rows.Err()
}

func (s *PostgresStore) UpdateManualTriviaQuestionLastUsed(questionID int) error {
	today := time.Now().Format("2006-01-02")
	statement := "UPDATE manualtriviaquestions SET lastUsed = $2 WHERE id = $1 RETURNING id;"
	var id int
	return s.Connection.QueryRow(statement, questionID, today).Scan(&id)
}

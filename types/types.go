package types

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const (
	QUESTION_TYPE_TEXT int = iota + 1
	QUESTION_TYPE_IMAGE
	QUESTION_TYPE_FLAG
	QUESTION_TYPE_MAP
)

type TriviaDto struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	MaxScore  int           `json:"maxScore"`
	Questions []QuestionDto `json:"questions"`
}

type QuestionDto struct {
	ID                 int            `json:"id"`
	Type               string         `json:"type"`
	Question           string         `json:"question"`
	MapName            string         `json:"mapName"`
	Map                MapDto         `json:"map"`
	Highlighted        string         `json:"highlighted"`
	FlagCode           string         `json:"flagCode"`
	FlagUrl            sql.NullString `json:"flagUrl"`
	ImageURL           string         `json:"imageUrl"`
	ImageAttributeName string         `json:"imageAttributeName"`
	ImageAttributeURL  string         `json:"imageAttributeUrl"`
	ImageWidth         int            `json:"imageWidth"`
	ImageHeight        int            `json:"imageHeight"`
	ImageAlt           string         `json:"imageAlt"`
	Explainer          string         `json:"explainer"`
	Answers            []AnswerDto    `json:"answers"`
}

type AnswerDto struct {
	Text      string         `json:"text"`
	IsCorrect bool           `json:"isCorrect"`
	FlagCode  string         `json:"flagCode"`
	FlagUrl   sql.NullString `json:"flagUrl"`
}

type MapDto struct {
	ID        int             `json:"id"`
	Key       string          `json:"key"`
	ClassName string          `json:"className"`
	Label     string          `json:"label"`
	ViewBox   string          `json:"viewBox"`
	Elements  []MapElementDto `json:"elements"`
}

type MapElementDto struct {
	EntryID    int    `json:"entryId"`
	MapID      int    `json:"mapId"`
	Type       string `json:"type"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	D          string `json:"d"`
	Points     string `json:"points"`
	X          string `json:"x"`
	Y          string `json:"y"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	Cx         string `json:"cx"`
	Cy         string `json:"cy"`
	R          string `json:"r"`
	Transform  string `json:"transform"`
	XlinkHref  string `json:"xlinkHref"`
	ClipPath   string `json:"clipPath"`
	ClipPathId string `json:"clipPathId"`
	X1         string `json:"x1"`
	Y1         string `json:"y1"`
	X2         string `json:"x2"`
	Y2         string `json:"y2"`
}

type MappingEntryDto struct {
	ID               int             `json:"id"`
	GroupID          int             `json:"groupId"`
	Name             string          `json:"name"`
	Code             string          `json:"code"`
	FlagUrl          string          `json:"flagUrl"`
	SVGName          string          `json:"svgName"`
	AlternativeNames *pq.StringArray `json:"alternativeNames"`
	Prefixes         *pq.StringArray `json:"prefixes"`
	Grouping         string          `json:"grouping"`
}

type ManualTriviaQuestion struct {
	ID                 int          `json:"id"`
	TypeID             int          `json:"typeId"`
	CategoryID         int          `json:"categoryId"`
	Question           string       `json:"question"`
	Map                string       `json:"map"`
	Highlighted        string       `json:"highlighted"`
	FlagCode           string       `json:"flagCode"`
	ImageURL           string       `json:"imageUrl"`
	ImageAttributeName string       `json:"imageAttributeName"`
	ImageAttributeURL  string       `json:"imageAttributeUrl"`
	ImageWidth         int          `json:"imageWidth"`
	ImageHeight        int          `json:"imageHeight"`
	ImageAlt           string       `json:"imageAlt"`
	Explainer          string       `json:"explainer"`
	LastUsed           sql.NullTime `json:"lastUsed"`
	QuizDate           sql.NullTime `json:"quizDate"`
	LastUpdated        time.Time    `json:"lastUpdated"`
}

type TriviaQuestionCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	ImageOnly bool   `json:"imageOnly"`
}

type TriviaQuestion struct {
	ID                 int    `json:"id"`
	TriviaId           int    `json:"triviaId"`
	TypeID             int    `json:"typeId"`
	Question           string `json:"question"`
	Map                string `json:"map"`
	Highlighted        string `json:"highlighted"`
	FlagCode           string `json:"flagCode"`
	ImageURL           string `json:"imageUrl"`
	ImageAttributeName string `json:"imageAttributeName"`
	ImageAttributeURL  string `json:"ImageAttributeUrl"`
	ImageWidth         int    `json:"imageWidth"`
	ImageHeight        int    `json:"imageHeight"`
	ImageAlt           string `json:"imageAlt"`
	Explainer          string `json:"explainer"`
}

type TriviaAnswer struct {
	ID               int    `json:"id"`
	TriviaQuestionID int    `json:"triviaQuestionId"`
	Text             string `json:"text"`
	IsCorrect        bool   `json:"isCorrect"`
	FlagCode         string `json:"flagCode"`
}

type ManualTriviaAnswer struct {
	ID                     int    `json:"id"`
	ManualTriviaQuestionID int    `json:"manualTriviaQuestionId"`
	Text                   string `json:"text"`
	IsCorrect              bool   `json:"isCorrect"`
	FlagCode               string `json:"flagCode"`
}

var TopLandmass = []string{
	"Russia",
	"Canada",
	"China",
	"United States",
	"Brazil",
	"Australia",
	"India",
	"Argentina",
	"Kazakhstan",
	"Algeria",
	"Democratic Republic of the Congo",
	"Denmark",
	"Saudi Arabia",
	"Mexico",
	"Indonesia",
	"Sudan",
	"Libya",
	"Iran",
	"Mongolia",
	"Peru",
	"Chad",
	"Niger",
	"Angola",
	"Mali",
	"South Africa",
	"Colombia",
	"Ethiopia",
	"Bolivia",
	"Mauritania",
	"Egypt",
	"Tanzania",
	"Nigeria",
	"Venezuela",
	"Pakistan",
	"Namibia",
	"Mozambique",
	"Turkey",
	"Chile",
	"Zambia",
	"Myanmar",
	"Afghanistan",
	"Somalia",
	"Central African Republic",
	"South Sudan",
	"Ukraine",
	"Madagascar",
	"Botswana",
	"Kenya",
	"France",
	"Yemen",
	"New Zealand",
}

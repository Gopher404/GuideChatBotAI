package usecase

import (
	"fmt"
	"log"
	"main/domain"
	"math/rand"
	"strings"
)

type AttractionsFinder interface {
	Find(words []string) ([]domain.Attraction, error)
}

type AIRequester interface {
	Request(message []domain.AIMessage, temperature float32) (string, error)
}

type TourMaker struct {
	attractions AttractionsFinder
	AI          AIRequester
}

func NewTourMaker(attractions AttractionsFinder, AI AIRequester) *TourMaker {
	return &TourMaker{
		attractions: attractions,
		AI:          AI,
	}
}

type MinAttraction struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

func (tm *TourMaker) NewTour(userMessage string) (string, error) {
	messages := []domain.AIMessage{
		{"user", fmt.Sprintf(promptGetPreferences, userMessage)},
	}

	preferencesString, err := tm.AI.Request(messages, 0.3)
	if err != nil {
		return "", fmt.Errorf("ai: cannot get preferences: %w", err)
	}

	messages = append(messages, domain.AIMessage{Role: "assistant", Content: preferencesString})

	preferences := tm.parseAIPreferences(preferencesString)
	if len(preferences) == 0 {
		return "", fmt.Errorf("parseAIPreferences: cannot get preferences from: \"%s\"", preferencesString)
	}

	log.Println("preferences: ", preferences)

	messages = []domain.AIMessage{
		{"user", fmt.Sprint(promptGetDays, userMessage)},
	}

	daysString, err := tm.AI.Request(messages, 0.1)
	if err != nil {
		return "", fmt.Errorf("ai: cannot get days: %w", err)
	}

	log.Println("daysString: ", daysString)

	days := tm.findNum(daysString)
	if days == 0 {
		days = 7
	}

	messages = []domain.AIMessage{
		{"user", fmt.Sprint(promptGetCountPlacesOnDay, userMessage)},
	}

	countPlacesByDayString, err := tm.AI.Request(messages, 0.1)
	if err != nil {
		return "", fmt.Errorf("ai: cannot get days: %w", err)
	}

	log.Println("countPlacesByDayString: ", countPlacesByDayString)

	countPlacesByDay := tm.findNum(countPlacesByDayString)
	if countPlacesByDay == 0 {
		countPlacesByDay = 3
	}

	attractions, err := tm.attractions.Find(preferences)
	if err != nil {
		return "", fmt.Errorf("db: cannot get attractions: %w", err)
	}

	if len(attractions) < 10 {
		return "", fmt.Errorf("count of attractions too small (%d)", len(attractions))
	}

	minAttractions := make([]MinAttraction, len(attractions))
	for i, attraction := range attractions {
		minAttractions[i] = MinAttraction{
			Id:          attraction.Id,
			Name:        attraction.Name,
			Description: attraction.Description,
			Location:    attraction.Location,
		}
	}

	messages = []domain.AIMessage{}

	for day := 0; day < days; day++ {
		rand.Shuffle(len(minAttractions), func(i, j int) {
			minAttractions[i], minAttractions[j] = minAttractions[j], minAttractions[i]
		})

		msg := domain.AIMessage{
			Role:    "user",
			Content: createPromptGenerateTourOnDay(userMessage, days, minAttractions[:10]),
		}

		messages = append(messages, msg)

		resp, err := tm.AI.Request(messages, 0.5)
		if err != nil {
			return "", fmt.Errorf("ai: cannot create tour: %w", err)
		}

		messages = append(messages, domain.AIMessage{Role: "assistant", Content: resp})

		log.Printf("resp: '%s'\n", resp)
	}

	return fmt.Sprint("preferences: ", preferences, "days: ", days, "countPlacesByDay: ", countPlacesByDay), nil
}

func (tm *TourMaker) parseAIPreferences(msg string) []string {
	var (
		idx1 int
		idx2 int
	)

	for i := 0; i < len(msg); i++ {
		if msg[i] == '[' || msg[i] == '{' {
			idx1 = i
			break
		}
	}
	for i := len(msg) - 1; i >= 0; i-- {
		if msg[i] == ']' || msg[i] == '}' {
			idx2 = i
			break
		}
	}

	if idx1 >= idx2 {
		return []string{}
	}

	var words []string

	msg = msg[idx1:idx2+1] + " "

	for {
		idx1 = strings.Index(msg, "\"")
		if idx1 == -1 {
			idx1 = strings.Index(msg, "'")
			if idx1 == -1 {
				break
			}
		}

		idx2 = Index(msg, '"', idx1+1)

		if idx2 == -1 {
			idx2 = Index(msg, '\'', idx1+1)
			if idx2 == -1 {
				break
			}
		}

		words = append(words, msg[idx1+1:idx2])

		msg = msg[idx2+1:]
	}

	return words
}

func (tm *TourMaker) findNum(s string) int {
	var days int
	findNum := false

	for _, a := range s {
		d, ok := digits[a]
		if ok && !findNum {
			findNum = true
			days = d

		} else if findNum {

			if ok {
				days = days*10 + d
			} else {
				break
			}
		}
	}

	return days
}

var digits = map[rune]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
}

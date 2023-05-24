package util

import (
  "encoding/json"
  "strings"
)

type Choice struct {
  Text  string `json:"text"`
  Valid bool   `json:"valid"`
}

type Quiz struct {
  Question string   `json:"question"`
  Choices  []Choice `json:"choices"`
}

type Choices []Choice
type Quizzes []Quiz

func ParseQuizzes(data []byte) (Quizzes, error) {

  var err error

  quizzes := make(Quizzes, 0)

  if err = json.Unmarshal(data, &quizzes); err != nil {

    return quizzes, err
  }

  return quizzes, nil
}

func FindValidChoiceFromChoices(source Choices) Choices {

  choices := make(Choices, 0)

  for _, choice := range source {

    if choice.Valid {

      choices = append(choices, choice)
    }
  }

  return choices
}

func Valid(source Choices, target Choices) bool {

  var found bool
  data := FindValidChoiceFromChoices(source)
  for _, attempt := range target {

    found = false
    for _, choice := range data {

      // check available answer in data source
      if attempt.Text == choice.Text {

        found = true
        break
      }
    }

    if !found {

      return false
    }
  }

  return true
}

func QuizScore(source Quizzes, target Quizzes) int {

  score := 0
  n := len(source)

  // check the same question
  for _, attempt := range target {

    for _, quiz := range source {

      // check the same question, case-insensitive
      if strings.ToUpper(quiz.Question) == strings.ToUpper(attempt.Question) {

        if Valid(quiz.Choices, attempt.Choices) {

          score++
        }
      }
    }
  }

  return int(float64(score/n) * 100)
}

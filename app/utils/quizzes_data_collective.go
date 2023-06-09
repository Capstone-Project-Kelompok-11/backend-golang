package util

import (
  "lms/app/models"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/pp"
)

func QuizzesDataCollective(ctx *swag.SwagContext, quizzes []models.Quizzes) []Quizzes {

  pp.Void(ctx)

  // []quizzes -> quiz (quizzes) -> []quizzes

  res := make([]Quizzes, 0)

  for _, quiz := range quizzes {

    var dataQuizzes Quizzes
    if dataQuizzes, _ = ParseQuizzes([]byte(quiz.Data)); dataQuizzes != nil {

      // randomize quizzes without showing valid answer
      data := QuizRandHideValid(dataQuizzes)

      res = append(res, data)
    }
  }

  return res
}

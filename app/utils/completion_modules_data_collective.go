package util

import (
  "lms/app/models"
  "skfw/papaya/bunny/swag"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
)

func CompletionModulesDataCollective(ctx *swag.SwagContext, data []models.CompletionModules) []m.KMapImpl {

  pp.Void(ctx)

  var err error
  res := make([]m.KMapImpl, 0)

  pp.Void(err)

  for _, completionModule := range data {

    mm := &m.KMap{
      "score": completionModule.Score,
    }

    res = append(res, mm)
  }

  return res
}

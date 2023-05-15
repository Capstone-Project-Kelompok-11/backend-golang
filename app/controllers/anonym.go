package controllers

import (
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/koala/mapping"
)

func AnonymController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  router.Get("/ping", &mapping.KMap{
    "request":     nil,
    "description": "Testing Response",
    "responses":   swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })
}

package controllers

import (
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
)

func TestController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  router.Get("/auth/strict/ping", &m.KMap{
    "AuthToken": true,
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })

  router.Get("/auth/ping", &m.KMap{
    "AuthToken":          true,
    "DisableDeviceCheck": true,
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })

  router.Get("/ping", &m.KMap{
    "request":  nil,
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })
}

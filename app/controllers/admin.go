package controllers

import (
  "fmt"
  util "lms/app/utils"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
)

func AdminController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  router.Post("/course/image/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Product Image",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.CreatedJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()

    ids := m.KValueToString(kReq.Query.Get("id"))
    name := "example"

    fmt.Println(ids)

    return util.SwagSaveImageX256(ctx, name)
  })
}

package controllers

import (
  util "lms/app/utils"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func AnonymController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  router.Get("/ping", &m.KMap{
    "description": "Testing Response",
    "request":     nil,
    "responses":   swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })

  router.Get("/image/:img", &m.KMap{
    "description": "Get Public Image",
    "request": &m.KMap{
      "params": &m.KMap{
        "#img": "string",
      },
    },
    "responses": []byte{},
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()

    img := util.SafePathName(m.KValueToString(kReq.Path.Get("img")))

    if img != "" {

      img = posix.KPathNew("assets/public/images").JoinStr(img)

      file := kio.KFileNew(img)

      if file.IsExist() {

        if file.IsFile() {

          return ctx.SendFile(img, true)
        }
      }

      return ctx.BadRequest(kornet.Msg("image not found", true))
    }

    return ctx.BadRequest(kornet.Msg("invalid path", true))
  })
}

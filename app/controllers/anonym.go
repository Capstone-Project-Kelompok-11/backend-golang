package controllers

import (
  "lms/app/models"
  "lms/app/repository"
  util "lms/app/utils"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func AnonymController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  courseRepo, _ := repository.CourseRepositoryNew(DB)

  router.Get("/ping", &m.KMap{
    "description": "Testing Response",
    "request":     nil,
    "responses":   swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    return ctx.Message("pong")
  })

  router.Get("/image/:src", &m.KMap{
    "description": "Get Public Image",
    "request": &m.KMap{
      "params": &m.KMap{
        "#src": "string",
      },
    },
    "responses": nil,
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()

    src := util.SafePathName(m.KValueToString(kReq.Path.Get("src")))

    if src != "" {

      src = posix.KPathNew("assets/public/images").JoinStr(src)

      file := kio.KFileNew(src)

      if file.IsExist() {

        if file.IsFile() {

          return ctx.SendFile(src, true)
        }
      }

      return ctx.BadRequest(kornet.Msg("image not found", true))
    }

    return ctx.BadRequest(kornet.Msg("invalid path", true))
  })

  router.Get("/document/:src", &m.KMap{
    "description": "Get Public Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "#src": "string",
      },
    },
    "responses": nil,
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()

    src := util.SafePathName(m.KValueToString(kReq.Path.Get("src")))

    if src != "" {

      src = posix.KPathNew("assets/public/documents").JoinStr(src)

      file := kio.KFileNew(src)

      if file.IsExist() {

        if file.IsFile() {

          return ctx.SendFile(src, true)
        }
      }

      return ctx.BadRequest(kornet.Msg("document not found", true))
    }

    return ctx.BadRequest(kornet.Msg("invalid path", true))
  })

  router.Get("/video/:src", &m.KMap{
    "description": "Get Public Video",
    "request": &m.KMap{
      "params": &m.KMap{
        "#src": "string",
      },
    },
    "responses": nil,
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()

    src := util.SafePathName(m.KValueToString(kReq.Path.Get("src")))

    if src != "" {

      src = posix.KPathNew("assets/public/videos").JoinStr(src)

      file := kio.KFileNew(src)

      if file.IsExist() {

        if file.IsFile() {

          return ctx.SendFile(src, true)
        }
      }

      return ctx.BadRequest(kornet.Msg("image not found", true))
    }

    return ctx.BadRequest(kornet.Msg("invalid path", true))
  })

  router.Get("/courses", &m.KMap{
    "description": "Catch All Courses",
    "request": &m.KMap{
      "params": &m.KMap{
        "page": "number",
        "size": "number",
      },
    },
    "responses": swag.OkJSON([]m.KMapImpl{
      &m.KMap{
        "id":          "string",
        "name":        "string",
        "description": "string",
        "thumbnail":   "string",
        "video":       "string",
        "document":    "string",
        "price":       "number",
        "level":       "string",
        "rating":      "number",
        "finished":    "number",
        "members":     "number",
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var data []models.Courses

    kReq, _ := ctx.Kornet()

    page := util.ValueToInt(kReq.Query.Get("page"))
    size := util.ValueToInt(kReq.Query.Get("size"))

    if data, err = courseRepo.CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all courses", false), util.CourseDataCollective(data)))
  })
}

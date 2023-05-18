package controllers

import (
  "encoding/json"
  "github.com/shopspring/decimal"
  "lms/app/models"
  "lms/app/repository"
  util "lms/app/utils"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon/easy"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

func AdminController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  moduleRepo, _ := repository.ModuleRepositoryNew(DB)

  pp.Void(userRepo)

  router.Post("/course/thumbnail/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Course Thumbnail",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if check, _ := courseRepo.Find("id = ?", courseId); check != nil {

          if check.Thumbnail != "" {

            return util.SwagSaveImage(ctx, check.Thumbnail, func(name string) error {

              check.Thumbnail = name

              return courseRepo.Update(check, "id = ?", check.ID)
            })
          }

          return ctx.BadRequest(kornet.Msg("unknown file name to save", true))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/course/document/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Course Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if check, _ := courseRepo.Find("id = ?", courseId); check != nil {

          if check.Document != "" {

            return util.SwagSaveDocument(ctx, check.Document, func(name string) error {

              check.Document = name

              return courseRepo.Update(check, "id = ?", check.ID)
            })
          }

          return ctx.BadRequest(kornet.Msg("unknown file name to save", true))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/module/thumbnail/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Module Thumbnail",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          if check.Thumbnail != "" {

            return util.SwagSaveImage(ctx, check.Thumbnail, func(name string) error {

              check.Thumbnail = name

              return moduleRepo.Update(check, "id = ?", check.ID)
            })
          }

          return ctx.BadRequest(kornet.Msg("unknown file name to save", true))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/module/document/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Course Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          if check.Document != "" {

            return util.SwagSaveDocument(ctx, check.Document, func(name string) error {

              check.Document = name

              return moduleRepo.Update(check, "id = ?", check.ID)
            })
          }

          return ctx.BadRequest(kornet.Msg("unknown file name to save", true))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/course", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Course",
    "request": &m.KMap{
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "thumbnail?":  "string",
        "video?":      "string",
        "document?":   "string",
        "price":       "number",
        "level":       "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var body m.KMapImpl

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing body data into json format", true))
        }

        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))
        thumbnail := m.KValueToString(body.Get("thumbnail"))
        video := m.KValueToString(body.Get("video"))
        document := m.KValueToString(body.Get("document"))
        level := m.KValueToString(body.Get("level"))
        price := decimal.NewFromInt(util.ValueToInt64(body.Get("price")))

        if check, _ := courseRepo.Find("name = ?", name); check != nil {

          return ctx.BadRequest(kornet.Msg("course already exists", true))
        }

        if _, err = courseRepo.Create(&models.Courses{
          Model:       &easy.Model{},
          UserID:      user.ID,
          Name:        name,
          Description: description,
          Thumbnail:   thumbnail,
          Video:       video,
          Document:    document,
          Level:       level,
          Price:       price,
        }); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.Created(kornet.Msg("successful create new course", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/course", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "thumbnail?":  "string",
        "video?":      "string",
        "document?":   "string",
        "price":       "number",
        "level":       "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var body m.KMapImpl

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing body data into json format", true))
        }

        courseId := m.KValueToString(kReq.Query.Get("id"))

        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))
        thumbnail := m.KValueToString(body.Get("thumbnail"))
        video := m.KValueToString(body.Get("video"))
        document := m.KValueToString(body.Get("document"))
        level := m.KValueToString(body.Get("level"))
        price := decimal.NewFromInt(util.ValueToInt64(body.Get("price")))

        if check, _ := courseRepo.Find("id = ?", courseId); check != nil {

          check.Name = name
          check.Description = description
          check.Thumbnail = thumbnail
          check.Video = video
          check.Document = document
          check.Level = level
          check.Price = price

          if err = courseRepo.Update(check, "id = ?", check.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful update course", false))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/course", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Remove Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if err = courseRepo.Remove("id = ?", courseId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.Msg("successful delete course", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/module", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Module",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "thumbnail?":  "string",
        "video?":      "string",
        "document?":   "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var body m.KMapImpl

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing body data into json format", true))
        }

        courseId := m.KValueToString(kReq.Query.Get("id"))

        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))
        thumbnail := m.KValueToString(body.Get("thumbnail"))
        video := m.KValueToString(body.Get("video"))
        document := m.KValueToString(body.Get("document"))

        if _, err = courseRepo.Find("id = ?", courseId); err != nil {

          return ctx.BadRequest(kornet.Msg("course not found", true))
        }

        if check, _ := moduleRepo.Find("name = ?", name); check != nil {

          return ctx.BadRequest(kornet.Msg("module already exists", true))
        }

        if _, err = moduleRepo.Create(&models.Modules{
          Model:       &easy.Model{},
          CourseID:    courseId,
          Name:        name,
          Description: description,
          Thumbnail:   thumbnail,
          Video:       video,
          Document:    document,
        }); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.Created(kornet.Msg("successful create new module", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/module", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Module",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "thumbnail?":  "string",
        "video?":      "string",
        "document?":   "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var body m.KMapImpl

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing body data into json format", true))
        }

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))
        thumbnail := m.KValueToString(body.Get("thumbnail"))
        video := m.KValueToString(body.Get("video"))
        document := m.KValueToString(body.Get("document"))

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          check.Name = name
          check.Description = description
          check.Thumbnail = thumbnail
          check.Video = video
          check.Document = document

          if err = moduleRepo.Update(check, "id = ?", check.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful update module", false))
        }

        return ctx.BadRequest(kornet.Msg("module not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/module", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Remove Module",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    if ctx.Event() {

      if user, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(user)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        if err = moduleRepo.Remove("id = ?", moduleId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.Msg("successful delete module", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

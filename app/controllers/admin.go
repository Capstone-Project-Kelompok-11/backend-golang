package controllers

import (
  "encoding/json"
  "github.com/shopspring/decimal"
  "lms/app/models"
  "lms/app/repository"
  util "lms/app/utils"
  "net/url"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "skfw/papaya/koala/tools/posix"
  "skfw/papaya/pigeon/easy"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
  "time"
)

func AdminController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  categoryRepo, _ := repository.CategoryRepositoryNew(DB)
  categoryCourseRepo, _ := repository.CategoryCourseRepositoryNew(DB)
  checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)
  moduleRepo, _ := repository.ModuleRepositoryNew(DB)
  quizRepo, _ := repository.QuizzesRepositoryNew(DB)
  bannerRepo, _ := repository.BannerRepositoryNew(DB)
  eventRepo, _ := repository.EventRepositoryNew(DB)
  assignRepo, _ := repository.AssignmentRepositoryNew(DB)
  completionCourseRepo, _ := repository.CompletionCourseRepositoryNew(DB)

  router.Post("/course/thumbnail/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Course Thumbnail",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

          return util.SwagSaveImage(ctx, check.Thumbnail, func(name string) error {

            check.Thumbnail = name

            return courseRepo.Update(check, "id = ?", check.ID)
          })
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/course/thumbnail", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Course Thumbnail",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

        if check, _ := courseRepo.Find("id = ?", courseId); check != nil {

          if check.Thumbnail != "" {

            name := check.Thumbnail
            check.Thumbnail = ""

            if err = courseRepo.Update(check, "id = ?", check.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return util.SwagRemoveImage(ctx, name)
          }

          return ctx.BadRequest(kornet.Msg("thumbnail already removed", true))
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
      "headers": &m.KMap{
        "Authorization": "string",
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

          return util.SwagSaveDocument(ctx, check.Document, func(name string) error {

            check.Document = name

            return courseRepo.Update(check, "id = ?", check.ID)
          })
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/course/document", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Course Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

        if check, _ := courseRepo.Find("id = ?", courseId); check != nil {

          if check.Document != "" {

            name := check.Document
            check.Document = ""

            if err = courseRepo.Update(check, "id = ?", check.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return util.SwagRemoveDocument(ctx, name)
          }

          return ctx.BadRequest(kornet.Msg("document already removed", true))
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
      "headers": &m.KMap{
        "Authorization": "string",
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

          return util.SwagSaveImage(ctx, check.Thumbnail, func(name string) error {

            check.Thumbnail = name

            return moduleRepo.Update(check, "id = ?", check.ID)
          })
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/module/thumbnail", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Module Thumbnail",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          if check.Thumbnail != "" {

            name := check.Thumbnail
            check.Thumbnail = ""

            if err = moduleRepo.Update(check, "id = ?", check.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return util.SwagRemoveImage(ctx, name)
          }

          return ctx.BadRequest(kornet.Msg("thumbnail already removed", true))
        }

        return ctx.BadRequest(kornet.Msg("module not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/module/document/upload", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Upload Module Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

          return util.SwagSaveDocument(ctx, check.Document, func(name string) error {

            check.Document = name

            return moduleRepo.Update(check, "id = ?", check.ID)
          })
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/module/document", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Module Document",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          if check.Document != "" {

            name := check.Document
            check.Document = ""

            if err = moduleRepo.Update(check, "id = ?", check.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return util.SwagRemoveDocument(ctx, name)
          }

          return ctx.BadRequest(kornet.Msg("document already removed", true))
        }

        return ctx.BadRequest(kornet.Msg("module not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/course", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Course",
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "video?":      "string",
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
        video := m.KValueToString(body.Get("video"))
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
          Video:       video,
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
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "video?":      "string",
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
      "headers": &m.KMap{
        "Authorization": "string",
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
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "video?":      "string",
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
        video := m.KValueToString(body.Get("video"))

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
          Video:       video,
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
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
        "video?":      "string",
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
        video := m.KValueToString(body.Get("video"))

        if check, _ := moduleRepo.Find("id = ?", moduleId); check != nil {

          check.Name = name
          check.Description = description
          check.Video = video

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
      "headers": &m.KMap{
        "Authorization": "string",
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

  router.Get("/module/quiz", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Module Quiz",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.CreatedJSON(&kornet.Result{
      Data: util.Quizzes{
        util.Quiz{
          Choices: util.Choices{
            util.Choice{},
          },
        },
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var quizzes *models.Quizzes

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        data := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        if quizzes, err = quizRepo.Find("module_id = ?", moduleId); quizzes != nil {

          var dataQuizzes util.Quizzes
          if dataQuizzes, err = util.ParseQuizzes([]byte(quizzes.Data)); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full quizzes", false), dataQuizzes))
        }

        return ctx.BadRequest(kornet.Msg("quizzes not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/module/quiz", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Module Quiz",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "quizzes": util.Quizzes{
          util.Quiz{
            Choices: util.Choices{
              util.Choice{},
            },
          },
        },
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var quizzes util.Quizzes

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        data := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        // re-parsing quizzes data, double-check
        if quizzes, err = util.ParseQuizzes([]byte(m.KMapEncodeJSON(data.Get("quizzes")))); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse quizzes", true))
        }

        dataQuizzes := m.KMapEncodeJSON(quizzes)

        if _, err = quizRepo.Find("module_id = ?", moduleId); err != nil {

          if _, err = quizRepo.Create(&models.Quizzes{
            Model:    &easy.Model{},
            ModuleID: moduleId,
            Data:     dataQuizzes,
          }); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.Created(kornet.Msg("successful create new quizzes", false))
        }

        return ctx.BadRequest(kornet.Msg("quizzes already exists", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/module/quiz", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Module Quiz",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "quizzes": util.Quizzes{
          util.Quiz{
            Choices: util.Choices{
              util.Choice{},
            },
          },
        },
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var check *models.Quizzes
    var quizzes util.Quizzes

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        data := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        // re-parsing quizzes data, double-check
        if quizzes, err = util.ParseQuizzes([]byte(m.KMapEncodeJSON(data.Get("quizzes")))); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse quizzes", true))
        }

        dataQuizzes := m.KMapEncodeJSON(quizzes)

        if check, err = quizRepo.Find("module_id = ?", moduleId); check != nil {

          check.Data = dataQuizzes

          if err = quizRepo.Update(check, "module_id = ?", moduleId); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful update quizzes", false))
        }

        return ctx.BadRequest(kornet.Msg("quizzes not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/module/quiz", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Module Quiz",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        if err = quizRepo.Remove("module_id = ?", moduleId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.Msg("successful delete quizzes", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Catch All Course",
    "request": m.KMap{
      "params": &m.KMap{
        "size":    "number",
        "page":    "number",
        "sort?":   "string",
        "search?": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
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
        "modules":     []models.Modules{},
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var userModel *mo.UserModel
    var data []models.Courses
    var ok bool

    if ctx.Event() {

      if userModel, ok = ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        size := util.ValueToInt(kReq.Query.Get("size"))
        page := util.ValueToInt(kReq.Query.Get("page"))
        sort := m.KValueToString(kReq.Query.Get("sort"))
        search := m.KValueToString(kReq.Query.Get("search"))

        if search, sort, err = util.SafeParseSearchAndSortOrder(search, sort); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        if data, err = courseRepo.PreFindAllAndOrder(size, page, "name "+sort, "name LIKE ?", search); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        collective := util.CourseDataCollective(data)

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full course", false), collective))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/checkout/history", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Checkout History",
    "request": m.KMap{
      "params": &m.KMap{
        "size": "number",
        "page": "number",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{
      Data: []m.KMapImpl{
        &m.KMap{
          "data": nil,
          "user": &m.KMap{
            "name":     "string",
            "username": "string",
            "email":    "string",
          },
          "course": &m.KMap{
            "name": "string",
          },
          "paid":   "boolean",
          "cancel": "boolean",
        },
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    kReq, _ := ctx.Kornet()
    size := util.ValueToInt(kReq.Query.Get("size"))
    page := util.ValueToInt(kReq.Query.Get("page"))

    var err error
    var data []models.Checkout

    if data, err = checkoutRepo.Unscoped().CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    exposed := make([]m.KMapImpl, len(data))

    for i, checkout := range data {

      user := &models.Users{}
      course := &models.Courses{}

      if user, err = userRepo.Find("id = ?", checkout.UserID); err != nil {

        continue
      }

      if course, err = courseRepo.Find("id = ?", checkout.CourseID); err != nil {

        continue
      }

      exposed[i] = &m.KMap{
        "data": checkout,
        "user": &m.KMap{
          "name":     user.Name.String,
          "username": user.Username,
          "email":    user.Email,
        },
        "course": &m.KMap{
          "name": course.Name,
        },
        "paid":   checkout.Verify,
        "cancel": checkout.DeletedAt.Valid,
      }
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("checkout history", false), exposed))
  })

  router.Post("/category", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Category",
    "request": &m.KMap{
      "params": &m.KMap{
        "id?": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        courseId := m.KValueToString(body.Get("id"))

        pp.Void(courseId)

        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))

        var category *models.Categories

        // check if exists
        if category, err = categoryRepo.Find("name = ?", name); err != nil {

          if category, err = categoryRepo.Create(&models.Categories{
            Model:       &easy.Model{}, // not easy anymore ...
            Name:        name,
            Description: description,
          }); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }
        }

        if courseId != "" {

          if _, err = categoryCourseRepo.Create(&models.CategoryCourses{
            Model:      &easy.Model{},
            CourseID:   courseId,
            CategoryID: category.ID,
          }); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }
        }

        return ctx.OK(kornet.Msg("successful create category", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))

  })

  router.Post("/course/category", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Course Category",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "category_id": "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()
        courseId := m.KValueToString(kReq.Query.Get("id"))

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        categoryId := m.KValueToString(body.Get("category_id"))

        if _, err = categoryCourseRepo.Create(&models.CategoryCourses{
          Model:      &easy.Model{},
          CourseID:   courseId,
          CategoryID: categoryId,
        }); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.Msg("successful link course category", false))

      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/category", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Category",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // category id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name":        "string",
        "description": "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        categoryId := m.KValueToString(kReq.Query.Get("id"))
        name := m.KValueToString(body.Get("name"))
        description := m.KValueToString(body.Get("description"))

        var check *models.Categories

        if check, err = categoryRepo.Find("id = ?", categoryId); check != nil {

          check.Name = name
          check.Description = description

          if err = categoryRepo.Update(check, "id = ?", check.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful update category", false))
        }

        return ctx.BadRequest(kornet.Msg("category not found", true))

      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Delete("/category", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Delete Category",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // category id
      },
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        categoryId := m.KValueToString(kReq.Query.Get("id"))

        var check *models.Categories

        if check, err = categoryRepo.Find("id = ?", categoryId); check != nil {

          if err = categoryRepo.Delete(check); err != nil { // hard delete

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful delete category", false))
        }
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/categories", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Get Categories",
    "request": &m.KMap{
      "params": &m.KMap{
        "size": "number",
        "page": "number",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        size := util.ValueToInt(kReq.Query.Get("size"))
        page := util.ValueToInt(kReq.Query.Get("page"))

        var categories []models.Categories

        if categories, err = categoryRepo.CatchAll(size, page); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful get categories", false), categories))

      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/banner", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Create Banner",
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "alt": "string",
      }),
    },
    "responses": swag.CreatedJSON(&kornet.Result{}),
  },
    func(ctx *swag.SwagContext) error {

      var err error

      pp.Void(err)

      if ctx.Event() {

        if userModel, ok := ctx.Target().(*mo.UserModel); ok {

          pp.Void(userModel)

          kReq, _ := ctx.Kornet()

          body := &m.KMap{}

          if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          alt := m.KValueToString(body.Get("alt"))

          if check, _ := bannerRepo.Create(&models.Banners{
            Model: &easy.Model{},
            Alt:   alt,
          }); check != nil {

            if check.Src == "" {

              check.Src, _ = util.GenUniqFileNameOutput("assets/public/images", "banner.png")

              if err = bannerRepo.Update(check, "id = ?", check.ID); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }
            }

            return util.SwagSaveImage(ctx, check.Src, func(filename string) error {

              check.Src = filename

              return bannerRepo.Update(check, "id = ?", check.ID)
            })
          }

          return ctx.InternalServerError(kornet.Msg("unable to create new banner", true))
        }
      }

      return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
    })

  router.Put("/banner", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Banner",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "alt": "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  },
    func(ctx *swag.SwagContext) error {

      var err error

      pp.Void(err)

      if ctx.Event() {

        if userModel, ok := ctx.Target().(*mo.UserModel); ok {

          pp.Void(userModel)

          kReq, _ := ctx.Kornet()

          body := &m.KMap{}

          if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          bannerId := m.KValueToString(kReq.Query.Get("id"))
          alt := m.KValueToString(body.Get("alt"))

          if check, _ := bannerRepo.Find("id = ?", bannerId); check != nil {

            if check.Src != "" {

              if err = util.SwagRemoveImage(ctx, check.Src); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              check.Src = "" // make empty string
            }

            if check.Src == "" {

              check.Src, _ = util.GenUniqFileNameOutput("assets/public/images", "banner.png")

              if err = bannerRepo.Update(check, "id = ?", check.ID); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }
            }

            if err = util.SwagSaveImage(ctx, check.Src, func(filename string) error {

              check.Alt = alt
              check.Src = filename

              return bannerRepo.Update(check, "id = ?", check.ID)
            }); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return ctx.OK(kornet.Msg("successful update banner", false))
          }

          return ctx.InternalServerError(kornet.Msg("unable to create new banner", true))
        }
      }

      return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
    })

  router.Delete("/banner", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Update Banner",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  },
    func(ctx *swag.SwagContext) error {

      var err error

      pp.Void(err)

      if ctx.Event() {

        if userModel, ok := ctx.Target().(*mo.UserModel); ok {

          pp.Void(userModel)

          kReq, _ := ctx.Kornet()

          bannerId := m.KValueToString(kReq.Query.Get("id"))

          if check, _ := bannerRepo.Find("id = ?", bannerId); check != nil {

            if check.Src != "" {

              if err = util.SwagRemoveImage(ctx, check.Src); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }
            }

            if err = bannerRepo.Delete(check); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful delete banner", false), nil))
          }

          return ctx.InternalServerError(kornet.Msg("unable to create new banner", true))
        }
      }

      return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
    })

  router.Get("/events", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Catch All Notifications",
    "request": &m.KMap{
      "params": &m.KMap{
        "size":       "number",
        "page":       "number",
        "after_at?":  "number", // timestamp
        "before_at?": "number", // timestamp
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        var afterAt, beforeAt int64

        size := util.ValueToInt(kReq.Query.Get("size"))
        page := util.ValueToInt(kReq.Query.Get("page"))
        afterAt = int64(util.ValueToInt(kReq.Query.Get("after_at")))
        beforeAt = int64(util.ValueToInt(kReq.Query.Get("before_at")))

        pp.Void(size, page, afterAt, beforeAt)

        if afterAt == 0 {

          afterAt = time.Now().UTC().AddDate(0, 0, -7).UnixMilli()
        }

        if beforeAt == 0 {

          beforeAt = time.Now().UTC().UnixMilli()
        }

        var events []models.Events

        if events, err = eventRepo.FindAll(size, page, "created_at BETWEEN ? AND ?", time.UnixMilli(afterAt), time.UnixMilli(beforeAt)); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        data := make([]m.KMapImpl, 0)

        for _, event := range events {

          exposed := &m.KMap{}
          userId := event.UserID

          if user, _ := userRepo.Find("id = ?", userId); user != nil {

            exposed.Put("user", &m.KMap{
              "name":     user.Name.String,
              "username": user.Username,
              "image":    user.Image,
            })
          }

          exposed.Put("event", &m.KMap{
            "name":        event.Name,
            "description": event.Description,
          })

          data = append(data, exposed)
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful get events", false), data))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course/resumes", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Catch All Course Resume",
    "request": &m.KMap{
      "params": &m.KMap{
        "size": "number",
        "page": "number",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var assigns []models.Assignments

    pp.Void(err, assigns)

    kReq, _ := ctx.Kornet()
    size := util.ValueToInt(kReq.Query.Get("size"))
    page := util.ValueToInt(kReq.Query.Get("page"))

    var URL *url.URL

    if URL, err = url.Parse(ctx.BaseURL()); err != nil {

      return ctx.InternalServerError(kornet.Msg("unable to parse base url", true))
    }

    if assigns, err = assignRepo.CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    data := make([]m.KMapImpl, 0)

    for _, assign := range assigns {

      URL.Path = "/api/v1/public/documents"
      URL.RawPath = URL.Path

      exposed := &m.KMap{}
      userId := assign.UserID

      if user, _ := userRepo.Find("id = ?", userId); user != nil {

        exposed.Put("user", &m.KMap{
          "name":     user.Name.String,
          "username": user.Username,
          "image":    user.Image,
        })
      }

      URL.Path = posix.KPathNew(URL.Path).JoinStr(assign.Document)
      URL.RawPath = URL.Path

      exposed.Put("data", &m.KMap{
        "id":           assign.ID,
        "document":     assign.Document,
        "document_url": URL.String(),
        "video":        assign.Video,
      })

      data = append(data, exposed)
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful get course resumes", false), data))
  })

  router.Post("/course/resume/grade", &m.KMap{
    "AuthToken":   true,
    "Admin":       true,
    "description": "Grade Course Resume",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // resume id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "grade": "number",
      }),
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var assign *models.Assignments

    pp.Void(err, assign)

    kReq, _ := ctx.Kornet()
    assignId := m.KValueToString(kReq.Query.Get("id"))

    body := &m.KMap{}

    if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    var completionCourse *models.CompletionCourses

    pp.Void(completionCourse)

    if assign, err = assignRepo.Find("id = ?", assignId); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", assign.UserID, assign.CourseID); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    grade := util.ValueToInt(body.Get("grade"))

    score := assign.Score

    if assign, err = assignRepo.Grade(assignId, grade); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    // back scoring by worst math formula
    if score > 0 {

      // z = (a + b) / c
      // a = (z * c) - b
      completionCourse.Score = (completionCourse.Score * 2) - score
    }

    // :)
    completionCourse.Score += grade
    completionCourse.Score /= 2

    if err = completionCourseRepo.Update(completionCourse, "id = ?", completionCourse.ID); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful grade course resume", false), assign))
  })
}

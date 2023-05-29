package controllers

import (
  "encoding/json"
  "fmt"
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

var minScoreRequired = 60

func CourseController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  moduleRepo, _ := repository.ModuleRepositoryNew(DB)
  quizRepo, _ := repository.QuizzesRepositoryNew(DB)

  completionModuleRepo, _ := repository.CompletionModuleRepositoryNew(DB)

  pp.Void(courseRepo)

  router.Get("/course", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Course",
    "request": m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&m.KMap{
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
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var userModel *mo.UserModel
    var data *models.Courses
    var ok bool

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok = ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if data, err = courseRepo.PreFindByCheckUserAndCourseId(userModel.ID, courseId); data != nil {

          collective := util.CourseDataCollective([]models.Courses{*data})

          if len(collective) > 0 {

            exposed := collective[0]

            moduleExposed := make([]m.KMapImpl, len(data.Modules))

            for i, module := range data.Modules {

              completion := true

              if _, err = completionModuleRepo.Find("user_id = ? AND module_id = ?", userModel.ID, module.ID); err != nil {

                completion = false
              }

              moduleExposed[i] = &m.KMap{
                "data":       module,
                "completion": completion,
              }
            }

            exposed.Set("modules", moduleExposed)

            return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full course", false), exposed))
          }

          return ctx.InternalServerError(kornet.Msg("something went wrong", true))
        }

        return ctx.BadRequest(kornet.Msg("user not enrolled course", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/quiz", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Course",
    "request": m.KMap{
      "params": &m.KMap{
        "id": "string", // module id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(kornet.Result{
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
    var checkout *models.Checkout
    var module *models.Modules
    var quizzes *models.Quizzes

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        // get course id from module id
        if module, err = moduleRepo.Find("id = ?", moduleId); module != nil {

          if checkout, err = checkoutRepo.Find("course_id = ? AND user_id = ?", module.CourseID, userModel.ID); checkout != nil {

            if quizzes, err = quizRepo.Find("module_id = ?", moduleId); quizzes != nil {

              var dataQuizzes util.Quizzes
              if dataQuizzes, err = util.ParseQuizzes([]byte(quizzes.Data)); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              // randomize quizzes without showing valid answer
              dataQuizzes = util.QuizRandHideValid(dataQuizzes)

              return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full quizzes", false), dataQuizzes))
            }

            return ctx.BadRequest(kornet.Msg("quizzes not found", true))
          }

          return ctx.BadRequest(kornet.Msg("user not enrolled course", true))
        }

        return ctx.BadRequest(kornet.Msg("module not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/quiz", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Course",
    "request": m.KMap{
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
    "responses": swag.OkJSON(kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkout *models.Checkout
    var module *models.Modules
    var quizzes *models.Quizzes

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        moduleId := m.KValueToString(kReq.Query.Get("id"))

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse request body", true))
        }

        dataQuizzesJSON := m.KMapEncodeJSON(body.Get("quizzes"))

        var dataQuizzes util.Quizzes
        var attemptQuizzes util.Quizzes

        if attemptQuizzes, err = util.ParseQuizzes([]byte(dataQuizzesJSON)); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        // get course id from module id
        if module, err = moduleRepo.Find("id = ?", moduleId); module != nil {

          if checkout, err = checkoutRepo.Find("course_id = ? AND user_id = ?", module.CourseID, userModel.ID); checkout != nil {

            if quizzes, err = quizRepo.Find("module_id = ?", moduleId); quizzes != nil {

              if dataQuizzes, err = util.ParseQuizzes([]byte(quizzes.Data)); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              // get quizzes score
              score := util.QuizScore(dataQuizzes, attemptQuizzes)

              if minScoreRequired <= score {

                // TODO
                // check all modules in data course
                // insert into completion_courses
                // make event notify to admin

                once := false

                pp.Void(&once)

                var completionModule *models.CompletionModules

                if completionModule, err = completionModuleRepo.Find("user_id = ? AND module_id = ?", userModel.ID, module.ID); completionModule != nil {

                  if !once {

                    completionModule.Score = score

                    if err = completionModuleRepo.Update(completionModule, "id = ?", completionModule.ID); err != nil {

                      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                    }

                    return ctx.OK(kornet.Msg(fmt.Sprintf("your score is %d passed", score), false))
                  }

                  return ctx.BadRequest(kornet.Msg("take quiz only once", true))
                }

                if _, err = completionModuleRepo.Create(&models.CompletionModules{
                  Model:    &easy.Model{},
                  UserID:   userModel.ID,
                  CourseID: checkout.CourseID,
                  ModuleID: module.ID,
                  Score:    score,
                }); err != nil {

                  return ctx.InternalServerError(kornet.Msg(err.Error(), true))

                }

                // passed by score
                return ctx.Message(fmt.Sprintf("your score is %d passed", score))
              }

              // not passed
              return ctx.OK(kornet.Msg(fmt.Sprintf("your score is %d not enough to passed", score), false))
            }

            return ctx.BadRequest(kornet.Msg("quizzes not found", true))
          }

          return ctx.BadRequest(kornet.Msg("user not enrolled course", true))
        }

        return ctx.BadRequest(kornet.Msg("module not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

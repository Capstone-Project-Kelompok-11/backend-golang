package controllers

import (
  "lms/app/models"
  "lms/app/repository"
  util "lms/app/utils"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

func CourseController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  moduleRepo, _ := repository.ModuleRepositoryNew(DB)
  quizRepo, _ := repository.QuizzesRepositoryNew(DB)

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

        if data, err = courseRepo.PreloadFindByCheckUserAndCourseId(userModel.ID, courseId); data != nil {

          collective := util.CourseDataCollective([]models.Courses{*data})

          if len(collective) > 0 {

            return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full course", false), collective[0]))
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
}

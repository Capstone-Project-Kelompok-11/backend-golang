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

  courseRepo, _ := repository.CourseRepositoryNew(DB)

  pp.Void(courseRepo)

  router.Get("/course", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Course",
    "request": m.KMap{
      "params": &m.KMap{
        "id": "string",
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
    var data *models.Courses
    var ok bool

    if ctx.Event() {

      if userModel, ok = ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if data, err = courseRepo.PreloadFind("checkout.verify = ? AND checkout.user_id = ? AND courses.id = ?", true, userModel.ID, courseId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        collective := util.CourseDataCollective([]models.Courses{*data})

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full course", false), collective))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

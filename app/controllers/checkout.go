package controllers

import (
  "lms/app/models"
  "lms/app/repository"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon/easy"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

func CheckoutController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)

  pp.Void(userRepo)

  // paid, checkout verify = true
  // unpaid, checkout verify = false
  // cancel, checkout verify = false & remove = true

  router.Get("/checkout/paid", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Checkout Verified",
    "request": m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkouts []models.Checkout

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        if checkouts, err = userRepo.CatchAllCheckoutVerified("id = ?", userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), checkouts))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/checkout/unpaid", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Checkout Non Verified",
    "request": m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkouts []models.Checkout

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        if checkouts, err = userRepo.CatchAllCheckoutNonVerified("id = ?", userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), checkouts))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/checkout/cancel", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Cancelling Checkout",
    "request": m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkouts []models.Checkout

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        if checkouts, err = userRepo.CatchAllCheckoutCancelled("id = ?", userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), checkouts))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  // unpaid, create new checkout

  router.Post("/checkout", &m.KMap{
    "AuthToken":   true,
    "description": "Checkout Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkout *models.Checkout
    var course *models.Courses

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        // check course if not exists
        if course, err = courseRepo.Find("id = ?", courseId); err != nil {

          return ctx.BadRequest(kornet.Msg("course not found", true))
        }

        // check checkout if not exists
        if checkout, err = checkoutRepo.Find("user_id = ?", userModel.ID); err != nil {

          // create new checkout
          if checkout, err = checkoutRepo.Create(&models.Checkout{
            Model:    &easy.Model{},
            UserID:   userModel.ID,
            CourseID: course.ID,
          }); err != nil {

            pp.Void(checkout)

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.Created(kornet.Msg("create new checkout", false))
        }

        return ctx.BadRequest(kornet.Msg("course already checkout", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  // canceling checkout course, remove checkout

  router.Delete("/checkout", &m.KMap{
    "AuthToken":   true,
    "description": "Checkout Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
  }, func(ctx *swag.SwagContext) error {

    var err error
    var course *models.Courses
    var checkout *models.Checkout

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        // check course if not exists
        if course, err = courseRepo.Find("id = ?", courseId); err != nil {

          pp.Void(course)

          return ctx.BadRequest(kornet.Msg("course not found", true))
        }

        // check checkout if exists
        if checkout, err = checkoutRepo.Find("user_id = ?", userModel.ID); checkout != nil {

          if err = checkoutRepo.Remove(checkout, "id = ?", checkout.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("checkout canceled", false))
        }

        return ctx.BadRequest(kornet.Msg("checkout not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  // paid, verifying checkout course

  router.Post("/checkout/verify", &m.KMap{
    "AuthToken":   true,
    "description": "Checkout Course",
    "request": &m.KMap{
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

        // verified checkout
        if err = checkoutRepo.PreloadVerifyByUserId(userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        // update member count in course member active
        if err = courseRepo.UpdateMemberCountByUserId(userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.Msg("checkout verified", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

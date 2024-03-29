package controllers

import (
  "encoding/json"
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

  router.Get("/checkout/history", &m.KMap{
    "AuthToken":   true,
    "description": "Catch All Checkout",
    "request": m.KMap{
      "params": &m.KMap{
        "size?": "number",
        "page?": "number",
      },
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

        kReq, _ := ctx.Kornet()
        size := util.ValueToInt(kReq.Query.Get("size"))
        page := util.ValueToInt(kReq.Query.Get("page"))

        if size == 0 {
          size = -1
        }
        if page == 0 {
          page = -1
        }

        if checkouts, err = checkoutRepo.Unscoped().FindAll(size, page, "user_id = ?", userModel.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        exposed := util.CheckoutDataCollective(ctx, userRepo, courseRepo, checkouts)

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), exposed))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

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

        exposed := util.CheckoutDataCollective(ctx, userRepo, courseRepo, checkouts)

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), exposed))
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

        exposed := util.CheckoutDataCollective(ctx, userRepo, courseRepo, checkouts)

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), exposed))
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

        exposed := util.CheckoutDataCollective(ctx, userRepo, courseRepo, checkouts)

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all checkouts", false), exposed))
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
        if checkout, err = checkoutRepo.Find("user_id = ? AND course_id = ?", userModel.ID, course.ID); err != nil {

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
        "id": "string", // checkout id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
  }, func(ctx *swag.SwagContext) error {

    var err error
    var checkout *models.Checkout

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        checkoutId := m.KValueToString(kReq.Query.Get("id"))

        // check checkout if exists
        if checkout, err = checkoutRepo.Find("id = ?", checkoutId); checkout != nil {

          if err = checkoutRepo.Remove("id = ?", checkout.ID); err != nil {

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
      "params": &m.KMap{
        "id?": "string", // checkout id (optional)
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "payment_method": "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        checkoutId := m.KValueToString(kReq.Query.Get("id"))

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        paymentMethod := m.KValueToString(body.Get("payment_method"))

        // unverified checkout
        var checkouts []models.Checkout

        if checkoutId != "" {

          if _, err = checkoutRepo.Find("id = ? AND user_id = ?", checkoutId, userModel.ID); err != nil {

            return ctx.BadRequest(kornet.Msg("checkout not found", true))
          }

          if checkouts, err = checkoutRepo.FindAll(-1, -1, "id = ? AND user_id = ? AND verify = ?", checkoutId, userModel.ID, false); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          // verified checkout
          if err = checkoutRepo.PreloadVerifyByIdAndUserId(checkoutId, userModel.ID, paymentMethod); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

        } else {

          if checkouts, err = checkoutRepo.FindAll(-1, -1, "user_id = ? AND verify = ?", userModel.ID, false); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          // verified checkout
          if err = checkoutRepo.PreloadVerifyByUserId(userModel.ID, paymentMethod); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }
        }

        for _, checkout := range checkouts {

          // update member count in course member active
          if err = courseRepo.UpdateMemberCountById(checkout.CourseID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }
        }

        return ctx.OK(kornet.Msg("checkout verified", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

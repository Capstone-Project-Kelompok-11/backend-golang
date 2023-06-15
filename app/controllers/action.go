package controllers

import (
  "database/sql"
  "encoding/json"
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
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
  bacx "skfw/papaya/pigeon/templates/basicAuth/util"
  "time"
)

func ActionController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)
  completionCourseRepo, _ := repository.CompletionCourseRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)

  pp.Void(completionCourseRepo)
  pp.Void(courseRepo)

  router.Get("/info", &m.KMap{
    "AuthToken":   true,
    "description": "Catch User Information",
    "request": m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{
      Data: &m.KMap{
        "name":         "string",
        "username":     "string",
        "image":        "string",
        "image_url":    "string",
        "email":        "string",
        "gender":       "string",
        "phone":        "string",
        "dob":          "string",
        "address":      "string",
        "country_code": "string",
        "city":         "string",
        "postal_code":  "string",
        "verify":       "boolean",
        "admin":        "boolean",
        "balance":      "number",
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var user *models.Users
    var URL *url.URL

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        if URL, err = url.Parse(ctx.BaseURL()); err != nil {

          URL = &url.URL{}
        }

        imagePub := posix.KPathNew("/api/v1/public/image")

        // get full user information
        if user, err = userRepo.Find("id = ?", userModel.ID); user != nil {

          if user.Image != "" {

            URL.Path = imagePub.Copy().JoinStr(user.Image)
            URL.RawPath = URL.Path

            user.Image = URL.String()
          }

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful get user information", false), &m.KMap{
            "name":         user.Name.String,
            "username":     user.Username,
            "image":        user.Image,
            "email":        user.Email,
            "gender":       user.Gender.String,
            "phone":        user.Phone.String,
            "dob":          user.DOB.Time,
            "address":      user.Address.String,
            "country_code": user.CountryCode.String,
            "city":         user.City.String,
            "postal_code":  user.PostalCode.String,
            "verify":       user.Verify,
            "admin":        user.Admin,
            "balance":      user.Balance.BigInt(),
          }))
        }

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/info", &m.KMap{
    "AuthToken":   true,
    "description": "Update User Information",
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "name": "string",
        //"username":          "string",
        "gender":            "string",
        "phone":             "string",
        "dob":               "string",
        "address":           "string",
        "country_code":      "string",
        "city":              "string",
        "postal_code":       "string",
        "confirm_password?": "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var user *models.Users
    var body m.KMapImpl
    var dobT time.Time

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }

        name := m.KValueToString(body.Get("name"))
        //username := m.KValueToString(body.Get("username"))
        gender := m.KValueToString(body.Get("gender"))
        phone := m.KValueToString(body.Get("phone"))
        dob := m.KValueToString(body.Get("dob"))
        address := m.KValueToString(body.Get("address"))
        countryCode := m.KValueToString(body.Get("country_code"))
        city := m.KValueToString(body.Get("city"))
        postalCode := m.KValueToString(body.Get("postal_code"))
        confirmPassword := m.KValueToString(body.Get("confirm_password"))

        //time.Parsing
        if dobT, err = time.Parse(time.RFC3339, dob); err != nil {

          return ctx.BadRequest(kornet.Msg("unable to parse date of birthday", true))
        }

        // get full user information
        if user, err = userRepo.Find("id = ?", userModel.ID); user != nil {

          if confirmPassword != "" {

            if !bacx.CheckPasswordHash(confirmPassword, user.Password) {

              return ctx.BadRequest(kornet.Msg("confirm password does not match", true))
            }
          }

          user.Name = sql.NullString{String: name, Valid: true}
          //user.Username = username
          user.Gender = sql.NullString{String: gender, Valid: true}
          user.Phone = sql.NullString{String: phone, Valid: true}
          user.DOB = sql.NullTime{Time: dobT, Valid: true}
          user.Address = sql.NullString{String: address, Valid: true}
          user.CountryCode = sql.NullString{String: countryCode, Valid: true}
          user.City = sql.NullString{String: city, Valid: true}
          user.PostalCode = sql.NullString{String: postalCode, Valid: true}

          if err = userRepo.Update(user, "id = ?", user.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful update user information", false))
        }

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/change/password", &m.KMap{
    "AuthToken":   true,
    "description": "Change User Password",
    "request": &m.KMap{
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "password": "string",
      }),
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var pass string
    var user *models.Users
    var body m.KMapImpl

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        body = &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }

        password := m.KValueToString(body.Get("password"))

        // get full user information
        if user, err = userRepo.Find("id = ?", userModel.ID); user != nil {

          if pass, err = bacx.HashPassword(password); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          user.Password = pass

          if err = userRepo.Update(user, "id = ?", user.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.OK(kornet.Msg("successful change user password", false))
        }

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }

    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))

  })

  router.Post("/profile/image/upload", &m.KMap{
    "AuthToken":   true,
    "description": "Upload User Profile Image",
    "request": m.KMap{
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

        if check, _ := userRepo.Find("id = ?", user.ID); check != nil {

          if check.Image != "" {

            if err = util.SwagRemoveImage(ctx, check.Image); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            statusCode := ctx.Response().StatusCode()
            if !(200 <= statusCode && statusCode < 300) {

              return ctx.InternalServerError(kornet.Msg("unable to remove image", true))
            }

          } else {

            check.Image, _ = util.GenUniqFileNameOutput("assets/public/images", "profile.png")

            if err = userRepo.Update(check, "id = ?", check.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }
          }

          return util.SwagSaveImageX256(ctx, check.Image, func(filename string) error {

            check.Image = filename

            return userRepo.Update(check, "id = ?", check.ID)
          })
        }

        return ctx.BadRequest(kornet.Msg("unable to get user information", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course/certificate", &m.KMap{
    "AuthToken":   true,
    "description": "User Course Certificate",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Message{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        var completionCourse *models.CompletionCourses
        var course *models.Courses

        if course, err = courseRepo.Find("id = ?", courseId); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to get course information", true))
        }

        if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); err != nil {

          return ctx.BadRequest(kornet.Msg("you doesn't completion this course", true))
        }

        var URL *url.URL

        if URL, err = url.Parse(ctx.BaseURL()); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse base url", true))
        }

        URL.Path = "/api/v1/public/cert"
        URL.RawPath = "/api/v1/public/cert"

        docURL := posix.KPathNew(URL.String())

        // assign not fulfilled certificate, or maybe not anymore

        if completionCourse.Certificate != "" {

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("certificate already exist", false), &m.KMap{

            "certificate": completionCourse.Certificate,
            "image":       docURL.Copy().JoinStr(completionCourse.Certificate + ".jpg"),
            "url":         docURL.Copy().JoinStr(completionCourse.Certificate + ".pdf"),
          }))
        }

        var certName string

        if certName, err = util.GenerateCertificateInCaches(course.Name, pp.Qstr(userModel.Name.String, userModel.Username)); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        certName = posix.KPathNew(certName).BaseStr()

        completionCourse.Certificate = certName

        if err = completionCourseRepo.Update(completionCourse, "id = ?", completionCourse.ID); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("create certificate", false), &m.KMap{

          "certificate": completionCourse.Certificate,
          "image":       docURL.Copy().JoinStr(completionCourse.Certificate + ".jpg"),
          "url":         docURL.Copy().JoinStr(completionCourse.Certificate + ".pdf"),
        }))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

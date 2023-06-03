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
  "skfw/papaya/koala/pp"
  "skfw/papaya/koala/tools/posix"
)

func AnonymController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  courseRepo, _ := repository.CourseRepositoryNew(DB)
  categoryRepo, _ := repository.CategoryRepositoryNew(DB)
  bannerRepo, _ := repository.BannerRepositoryNew(DB)
  reviewRepo, _ := repository.ReviewCourseRepositoryNew(DB)
  userRepo, _ := repository.UserRepositoryNew(DB)

  pp.Void(reviewRepo)
  pp.Void(userRepo)

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
        "page":      "number",
        "size":      "number",
        "search?":   "string",
        "sort?":     "string",
        "category?": "string", // csv by comma
      },
    },
    "responses": swag.OkJSON(&kornet.Result{
      Data: []m.KMapImpl{
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
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var data []models.Courses

    kReq, _ := ctx.Kornet()

    page := util.ValueToInt(kReq.Query.Get("page"))
    size := util.ValueToInt(kReq.Query.Get("size"))
    search := m.KValueToString(kReq.Query.Get("search"))
    sort := m.KValueToString(kReq.Query.Get("sort"))
    category := m.KValueToString(kReq.Query.Get("category"))

    // TODO
    pp.Void(category)

    if search != "" {

      if search, sort, err = util.SafeParseSearchAndSortOrder(search, sort); err != nil {

        return ctx.BadRequest(kornet.Msg("unable to parse search", true))
      }

      if data, err = courseRepo.FindAllAndOrder(size, page, "name "+sort, "name LIKE ?", search); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }

    } else {

      if data, err = courseRepo.PreCatchAll(size, page); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    exposed := util.CourseDataCollective(data)

    for _, course := range exposed {

      categories := make([]string, 0)

      if categoryModels, ok := course.Get("category_courses").([]models.CategoryCourses); ok {

        for _, categoryModel := range categoryModels {

          if category, _ := categoryRepo.Find("id", categoryModel.CategoryID); category != nil {

            categories = append(categories, category.Name)
          }
        }
      }

      course.Put("categories", categories)
      course.Del("category_courses")
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all courses", false), exposed))
  })

  router.Get("/cert/:src", &m.KMap{
    "description": "Get Public Certificate",
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

      src = posix.KPathNew("assets/public/caches").JoinStr(src)

      file := kio.KFileNew(src)

      if file.IsExist() {

        if file.IsFile() {

          return ctx.SendFile(src, true)
        }
      }

      return ctx.BadRequest(kornet.Msg("certificate not found", true))
    }

    return ctx.BadRequest(kornet.Msg("invalid path", true))
  })

  router.Get("/banners", &m.KMap{
    "description": "Get Public Banners",
    "request": &m.KMap{
      "params": &m.KMap{
        "page": "number",
        "size": "number",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{
      Data: []models.Banners{},
    }),
  },
    func(ctx *swag.SwagContext) error {

      var err error

      pp.Void(err)

      kReq, _ := ctx.Kornet()

      size := util.ValueToInt(kReq.Query.Get("size"))
      page := util.ValueToInt(kReq.Query.Get("page"))

      var banners []models.Banners

      if banners, err = bannerRepo.CatchAll(size, page); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }

      return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all banners", false), banners))
    })

  router.Get("/course/reviews", &m.KMap{
    "description": "Get Public Course Reviews",
    "request": &m.KMap{
      "params": &m.KMap{
        "page": "number",
        "size": "number",
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var reviews []models.ReviewCourses

    kReq, _ := ctx.Kornet()
    size := util.ValueToInt(kReq.Query.Get("size"))
    page := util.ValueToInt(kReq.Query.Get("page"))

    if reviews, err = reviewRepo.CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    exposed := make([]m.KMapImpl, 0)

    for _, review := range reviews {

      var user *models.Users

      if user, err = userRepo.Find("id", review.UserID); user != nil {

        exposed = append(exposed, &m.KMap{
          "id": review.ID,
          "user": &m.KMap{
            "name":     user.Name.String,
            "username": user.Username,
            "image":    user.Image,
          },
          "rating":  review.Rating,
          "comment": review.Description,
        })

        continue
      }
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all reviews", false), exposed))
  })
}

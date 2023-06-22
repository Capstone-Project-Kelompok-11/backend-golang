package controllers

import (
  "lms/app/models"
  "lms/app/repository"
  util "lms/app/utils"
  "net/url"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "skfw/papaya/koala/tools/posix"
  "strings"
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

  router.Get("/image/:src", &m.KMap{
    "description": "Get Public Image",
    "request": &m.KMap{
      "params": &m.KMap{
        "#src":    "string",
        "width?":  "number",
        "height?": "number",
        "scale?":  "number",
      },
    },
    "responses": nil,
  }, func(ctx *swag.SwagContext) error {

    var err error

    kReq, _ := ctx.Kornet()

    src := util.SafePathName(m.KValueToString(kReq.Path.Get("src")))
    width := util.ValueToInt(kReq.Query.Get("width"))
    height := util.ValueToInt(kReq.Query.Get("height"))
    scale := util.ValueToInt(kReq.Query.Get("scale"))

    var data []byte

    if src != "" {

      src = posix.KPathNew("assets/public/images").JoinStr(src)

      file := kio.KFileNew(src)

      if file.IsExist() {

        if file.IsFile() {

          if data, err = util.ImagePreview(ctx, src, width, height, scale); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          return ctx.Send(data)
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

    category = strings.TrimSpace(category)
    var categories []string

    if category != "" {

      categories = strings.Split(category, ",")

    } else {

      categories = []string{"all"}
    }

    for i, context := range categories {

      categories[i] = strings.TrimSpace(context)
    }

    if search, sort, err = util.SafeParseSearchAndSortOrder(search, sort); err != nil {

      return ctx.BadRequest(kornet.Msg(err.Error(), true))
    }

    if search != "%%" {

      if data, err = courseRepo.FindAllAndOrder(size, page, "name "+sort, "LOWER(name) LIKE LOWER(?)", search); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }

    } else {

      if data, err = courseRepo.PreCatchAll(size, page); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    exposed := util.CourseDataCollective(ctx, userRepo, data)
    reduced := util.CategoryDataCollective(categoryRepo, categories, exposed)

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all courses", false), reduced))
  })

  router.Get("/course/populars", &m.KMap{
    "description": "Catch All Courses Populars",
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

    category = strings.TrimSpace(category)
    var categories []string

    if category != "" {

      categories = strings.Split(category, ",")

    } else {

      categories = []string{"all"}
    }

    for i, context := range categories {

      categories[i] = strings.TrimSpace(context)
    }

    if search, sort, err = util.SafeParseSearchAndSortOrder(search, sort); err != nil {

      return ctx.BadRequest(kornet.Msg(err.Error(), true))
    }

    if search != "%%" {

      if data, err = courseRepo.FindAllAndOrderPopular(size, page, "name "+sort, "LOWER(name) LIKE LOWER(?)", search); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }

    } else {

      if data, err = courseRepo.PreCatchAllPopular(size, page); err != nil {

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    exposed := util.CourseDataCollective(ctx, userRepo, data)
    reduced := util.CategoryDataCollective(categoryRepo, categories, exposed)

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all courses", false), reduced))
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

      data := util.BannersDataCollective(ctx, banners)

      return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all banners", false), data))
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

    var URL *url.URL

    if URL, err = url.Parse(ctx.BaseURL()); err != nil {

      URL = &url.URL{}
    }

    imagePub := posix.KPathNew("/api/v1/public/image")

    if reviews, err = reviewRepo.CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    exposed := make([]m.KMapImpl, 0)

    for _, review := range reviews {

      var user *models.Users
      var course *models.Courses

      if user, err = userRepo.Find("id", review.UserID); user != nil {

        if user.Image != "" {

          URL.Path = imagePub.Copy().JoinStr(user.Image)
          URL.RawPath = URL.Path

          user.Image = URL.String()
        }

        if course, err = courseRepo.Find("id", review.CourseID); course != nil {

          exposed = append(exposed, &m.KMap{
            "id": review.ID,
            "user": &m.KMap{
              "name":     user.Name.String,
              "username": user.Username,
              "image":    user.Image,
            },
            "course": &m.KMap{
              "id":   course.ID,
              "name": course.Name,
            },
            "rating":  review.Rating,
            "comment": review.Description,
          })

          continue
        }
      }
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all reviews", false), exposed))
  })

  router.Get("/categories", &m.KMap{
    "description": "Get Public Categories",
    "request": &m.KMap{
      "params": &m.KMap{
        "page?": "number",
        "size?": "number",
      },
    },
    "responses": swag.OkJSON(&kornet.Result{
      Data: []string{},
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    kReq, _ := ctx.Kornet()

    size := util.ValueToInt(kReq.Query.Get("size"))
    page := util.ValueToInt(kReq.Query.Get("page"))

    // catch all data
    if page == 0 {
      page = -1
    }

    // catch all data
    if size == 0 {
      size = -1
    }

    var categories []models.Categories

    if categories, err = categoryRepo.CatchAll(size, page); err != nil {

      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
    }

    data := make([]string, len(categories))

    for i, category := range categories {

      data[i] = category.Name
    }

    return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch all categories", false), data))

  })
}

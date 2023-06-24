package controllers

import (
  "encoding/json"
  "fmt"
  "gorm.io/gorm"
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

var minScoreRequired = 60

func CourseController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)
  checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)
  courseRepo, _ := repository.CourseRepositoryNew(DB)
  moduleRepo, _ := repository.ModuleRepositoryNew(DB)
  quizRepo, _ := repository.QuizzesRepositoryNew(DB)
  eventRepo, _ := repository.EventRepositoryNew(DB)
  reviewRepo, _ := repository.ReviewCourseRepositoryNew(DB)

  completionModuleRepo, _ := repository.CompletionModuleRepositoryNew(DB)
  completionCourseRepo, _ := repository.CompletionCourseRepositoryNew(DB)
  assignRepo, _ := repository.AssignmentRepositoryNew(DB)
  categoryRepo, _ := repository.CategoryRepositoryNew(DB)

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
      "id":            "string",
      "name":          "string",
      "description":   "string",
      "thumbnail":     "string",
      "thumbnail_url": "string",
      "video":         "string",
      "document":      "string",
      "price":         "number",
      "level":         "string",
      "rating":        "number",
      "finished":      "number",
      "members":       "number",
      "modules":       []models.Modules{},
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var userModel *mo.UserModel
    var data *models.Courses
    var URL *url.URL
    var ok bool

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok = ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()

        courseId := m.KValueToString(kReq.Query.Get("id"))

        if URL, err = url.Parse(ctx.BaseURL()); err != nil {

          URL = &url.URL{}
        }

        imagePub := posix.KPathNew("/api/v1/public/image")

        pp.Void(URL, imagePub)

        if data, err = courseRepo.PreFindByCheckUserAndCourseId(userModel.ID, courseId); data != nil {

          collective := util.CourseDataCollective(ctx, userRepo, []models.Courses{*data})
          collective = util.CategoryDataCollective(categoryRepo, []string{}, collective)

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

            // check completion course
            exposed.Put("completion", false)

            var completionCourse *models.CompletionCourses

            if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); completionCourse != nil {

              pp.Void(completionCourse)
              exposed.Set("completion", true)
            }

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

              dataQuizzes := util.QuizzesDataCollective(ctx, []models.Quizzes{*quizzes})

              if len(dataQuizzes) > 0 {

                return ctx.OK(kornet.ResultNew(kornet.MessageNew("catch full quizzes", false), dataQuizzes[0]))
              }

              return ctx.OK(kornet.Msg("quizzes empty", true))
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
                courseId := module.CourseID

                var RecoverCompletionCourse = func(ctx *swag.SwagContext) error {

                  var completionModules []models.CompletionModules
                  var modules []models.Modules
                  var course *models.Courses

                  if completionModules, err = completionModuleRepo.FindAll(-1, -1, "course_id = ?", courseId); err != nil {

                    return err
                  }
                  if modules, err = moduleRepo.FindAll(-1, -1, "course_id = ?", courseId); err != nil {

                    return err
                  }

                  scoreAvg := 0
                  completion := false

                  if len(completionModules) > 0 {

                    completion = true

                    for _, completionModule := range completionModules {

                      moduleCompletion := false
                      for _, moduleCheck := range modules {

                        if completionModule.ModuleID == moduleCheck.ID {

                          scoreAvg += completionModule.Score
                          moduleCompletion = true
                          break
                        }
                      }

                      if !moduleCompletion {

                        completion = false
                        break
                      }
                    }
                  }

                  scoreAvg /= len(completionModules)

                  if completion {

                    if course, err = courseRepo.Find("id = ?", courseId); course != nil {

                      var completionCourse *models.CompletionCourses

                      if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); completionCourse != nil {

                        completionCourse.Score = scoreAvg

                        if err = completionCourseRepo.Update(completionCourse, "id = ?", completionCourse.ID); err != nil {

                          return err
                        }

                      } else {

                        // create completion
                        if _, err = completionCourseRepo.Create(&models.CompletionCourses{
                          Model:    &easy.Model{},
                          UserID:   userModel.ID,
                          CourseID: courseId,
                          Score:    scoreAvg,
                        }); err != nil {

                          return err
                        }
                      }

                      course.Finished += 1

                      // update course finished count
                      if err = courseRepo.Update(course, "id = ?", course.ID); err != nil {

                        return err
                      }

                      var event *models.Events

                      taskName := "completion course"
                      taskDescription := fmt.Sprintf("completion course %s by %s", course.Name, userModel.Username)

                      if event, err = eventRepo.Unscoped().Find("user_id = ? AND name = ? AND description = ?", userModel.ID, taskName, taskDescription); event != nil {

                        // hack, make it fresh again
                        currentTime := time.Now().UTC()
                        event.CreatedAt = currentTime
                        event.UpdatedAt = currentTime
                        event.DeletedAt = gorm.DeletedAt{Time: currentTime, Valid: false}

                        if err = eventRepo.Update(event, "id = ?", event.ID); err != nil {

                          return err
                        }

                      } else {

                        // create events
                        if _, err = eventRepo.Create(&models.Events{
                          Model:       &easy.Model{},
                          UserID:      userModel.ID,
                          Name:        taskName,
                          Description: taskDescription,
                        }); err != nil {

                          return err
                        }
                      }
                    }
                  }

                  return nil
                }

                once := false

                pp.Void(&once)

                var completionModule *models.CompletionModules

                if completionModule, err = completionModuleRepo.Find("user_id = ? AND module_id = ?", userModel.ID, module.ID); completionModule != nil {

                  if !once {

                    completionModule.Score = score

                    if err = completionModuleRepo.Update(completionModule, "id = ?", completionModule.ID); err != nil {

                      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                    }

                    if err = RecoverCompletionCourse(ctx); err != nil {

                      return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                    }

                    // passed by score
                    return ctx.OK(kornet.ResultNew(kornet.MessageNew(fmt.Sprintf("your score is %d passed", score), false), &m.KMap{
                      "score": score,
                    }))
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

                if err = RecoverCompletionCourse(ctx); err != nil {

                  return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                }

                // passed by score
                return ctx.OK(kornet.ResultNew(kornet.MessageNew(fmt.Sprintf("your score is %d passed", score), false), &m.KMap{
                  "score": score,
                }))
              }

              // not passed
              return ctx.OK(kornet.ResultNew(kornet.MessageNew(fmt.Sprintf("your score is %d not enough to passed", score), false), &m.KMap{
                "score": score,
              }))
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

  router.Post("/course/review", &m.KMap{
    "AuthToken":   true,
    "description": "Review Course After Completion",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "rating":  "number",
        "comment": "string",
      }),
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse request body", true))
        }

        rating := util.ValueToInt(body.Get("rating"))
        comment := m.KValueToString(body.Get("comment"))

        courseId := m.KValueToString(kReq.Query.Get("id"))

        var completionCourse *models.CompletionCourses
        var course *models.Courses

        rates := []int{0, 1, 2, 3, 4, 5}

        found := false
        for _, rate := range rates {

          if rate == rating {

            found = true
            break
          }
        }

        if !found {

          return ctx.BadRequest(kornet.Msg("invalid rating", true))
        }

        var reviewCourse *models.ReviewCourses

        if reviewCourse, err = reviewRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); reviewCourse != nil {

          return ctx.BadRequest(kornet.Msg("you have already reviewed this course", false))
        }

        if course, err = courseRepo.Find("id = ?", courseId); course != nil {

          if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); completionCourse != nil {

            if _, err = reviewRepo.Create(&models.ReviewCourses{
              Model:       &easy.Model{},
              CourseID:    course.ID,
              UserID:      userModel.ID,
              Description: comment,
              Rating:      rating,
            }); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            switch rating {
            case 1:
              course.Rating1 += 1
              break
            case 2:
              course.Rating2 += 1
              break
            case 3:
              course.Rating3 += 1
              break
            case 4:
              course.Rating4 += 1
              break
            case 5:
              course.Rating5 += 1
              break
            }

            // update course
            if err = courseRepo.Update(course, "id = ?", course.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            return ctx.OK(kornet.Msg("success", false))
          }

          return ctx.BadRequest(kornet.Msg("you doesn't completed this course", true))

        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/course/review", &m.KMap{
    "AuthToken":   true,
    "description": "Review Course After Completion",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // review id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "rating":  "number",
        "comment": "string",
      }),
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()

        body := &m.KMap{}

        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse request body", true))
        }

        rating := util.ValueToInt(body.Get("rating"))
        comment := m.KValueToString(body.Get("comment"))

        reviewId := m.KValueToString(kReq.Query.Get("id"))

        var completionCourse *models.CompletionCourses
        var course *models.Courses

        rates := []int{0, 1, 2, 3, 4, 5}

        found := false
        for _, rate := range rates {

          if rate == rating {

            found = true
            break
          }
        }

        if !found {

          return ctx.BadRequest(kornet.Msg("invalid rating", true))
        }

        var reviewCourse *models.ReviewCourses

        if reviewCourse, err = reviewRepo.Find("id = ?", reviewId); reviewCourse != nil {

          courseId := reviewCourse.CourseID

          if course, err = courseRepo.Find("id = ?", courseId); course != nil {

            if completionCourse, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); completionCourse != nil {

              ratingPrev := reviewCourse.Rating

              reviewCourse.Rating = rating
              reviewCourse.Description = comment

              if err = reviewRepo.Update(reviewCourse, "id = ?", reviewCourse.ID); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              // reset rating
              switch ratingPrev {
              case 1:
                course.Rating1 = pp.Lint(course.Rating1 > 0, course.Rating1-1, 0)
                break
              case 2:
                course.Rating2 = pp.Lint(course.Rating2 > 0, course.Rating2-1, 0)
                break
              case 3:
                course.Rating3 = pp.Lint(course.Rating3 > 0, course.Rating3-1, 0)
                break
              case 4:
                course.Rating4 = pp.Lint(course.Rating4 > 0, course.Rating4-1, 0)
                break
              case 5:
                course.Rating5 = pp.Lint(course.Rating5 > 0, course.Rating5-1, 0)
                break
              }

              // update rating
              switch rating {
              case 1:
                course.Rating1 += 1
                break
              case 2:
                course.Rating2 += 1
                break
              case 3:
                course.Rating3 += 1
                break
              case 4:
                course.Rating4 += 1
                break
              case 5:
                course.Rating5 += 1
                break
              }

              // update course
              if err = courseRepo.Update(course, "id = ?", course.ID); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              return ctx.OK(kornet.Msg("success", false))
            }

            return ctx.BadRequest(kornet.Msg("you doesn't completed this course", true))

          }

          return ctx.InternalServerError(kornet.Msg("course not found", true))
        }

        return ctx.BadRequest(kornet.Msg("review not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course/review", &m.KMap{
    "AuthToken":   true,
    "description": "Review Course After Completion",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()
        courseId := m.KValueToString(kReq.Query.Get("id"))

        var reviewCourse *models.ReviewCourses

        if reviewCourse, err = reviewRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); reviewCourse != nil {

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("success", false), &m.KMap{
            "id":      reviewCourse.ID,
            "rating":  reviewCourse.Rating,
            "comment": reviewCourse.Description,
          }))
        }

        return ctx.BadRequest(kornet.Msg("review not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course/resume", &m.KMap{
    "AuthToken":   true,
    "description": "Resume Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)

        kReq, _ := ctx.Kornet()
        courseId := m.KValueToString(kReq.Query.Get("id"))

        var URL *url.URL

        if URL, err = url.Parse(ctx.BaseURL()); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parse base url", true))
        }

        docPub := posix.KPathNew("/api/v1/public/document")

        var assign *models.Assignments

        if assign, err = assignRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); assign != nil {

          URL.Path = docPub.Copy().JoinStr(assign.Document)
          URL.RawPath = URL.Path
          assign.Document = URL.String()

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("success", false), &m.KMap{
            "id":       assign.ID,
            "document": assign.Document,
            "video":    assign.Video,
          }))
        }

        return ctx.BadRequest(kornet.Msg("assign not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Post("/course/resume", &m.KMap{
    "AuthToken":   true,
    "description": "Resume Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "video?": "string",
      }),
    },
    "response": swag.CreatedJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()
        courseId := m.KValueToString(kReq.Query.Get("id"))
        videoSource := m.KValueToString(kReq.Query.Get("video"))

        pp.Void(userModel, courseId, videoSource)

        var filename string

        filename, _ = util.GenUniqFileNameOutput("assets/public/documents", "resume.pdf")

        var assign *models.Assignments
        var course *models.Courses

        if course, err = courseRepo.Find("id = ?", courseId); course != nil {

          if assign, err = assignRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); err != nil {

            // try to create again
            if assign, err = assignRepo.Create(&models.Assignments{
              Model:    &easy.Model{},
              UserID:   userModel.ID,
              CourseID: courseId,
            }); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }
          }

          if assign.Document != "" {

            if util.SwagCheckDocumentExist(assign.Document) {

              if err = util.SwagRemoveDocument(ctx, assign.Document); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }

              statusCode := ctx.Response().StatusCode()

              if !(200 <= statusCode && statusCode < 300) {

                return nil
              }
            } else {

              assign.Document = ""
            }
          }

          if assign.Document == "" {

            assign.Document = filename

          } else {

            filename = assign.Document
          }

          assign.Video = videoSource

          if err = util.SwagSaveDocument(ctx, filename, func(filename string) error {

            assign.Document = filename

            return assignRepo.Update(assign, "id = ?", assign.ID)

          }); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          // fix swag error return

          statusCode := ctx.Response().StatusCode()

          if !(200 <= statusCode && statusCode < 300) {

            return nil
          }

          if err = assignRepo.Update(assign, "id = ?", assign.ID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          var event *models.Events

          taskName := "resume"
          taskDescription := fmt.Sprintf("resume %s by %s", course.Name, userModel.Username)

          if event, err = eventRepo.Find("user_id = ? AND name = ? AND description = ?", userModel.ID, taskName, taskDescription); event != nil {

            // hack, make it fresh again
            currentTime := time.Now().UTC()
            event.CreatedAt = currentTime
            event.UpdatedAt = currentTime
            event.DeletedAt = gorm.DeletedAt{Time: currentTime, Valid: false}

            if err = eventRepo.Update(event, "id = ?", event.ID); err != nil {

              return err
            }

          } else {

            // create events
            if _, err = eventRepo.Create(&models.Events{
              Model:       &easy.Model{},
              UserID:      userModel.ID,
              Name:        taskName,
              Description: taskDescription,
            }); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }
          }

          return ctx.Created(kornet.Msg("success create resume", false))
        }

        return ctx.BadRequest(kornet.Msg("course not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Put("/course/resume", &m.KMap{
    "AuthToken":   true,
    "description": "Resume Course",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // assign id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
      "body": swag.JSON(&m.KMap{
        "video?": "string",
      }),
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()
        assignId := m.KValueToString(kReq.Query.Get("id"))
        videoSource := m.KValueToString(kReq.Query.Get("video"))

        pp.Void(userModel, assignId, videoSource)

        var filename string

        filename, _ = util.GenUniqFileNameOutput("assets/public/documents", "resume.pdf")

        var assign *models.Assignments

        if assign, err = assignRepo.Find("id = ?", assignId); assign != nil {

          courseId := assign.CourseID

          var course *models.Courses

          if course, err = courseRepo.Find("id = ?", courseId); course != nil {

            if videoSource != "" {

              assign.Video = videoSource
            }

            if assign.Document != "" {

              if util.SwagCheckDocumentExist(assign.Document) {

                if err = util.SwagRemoveDocument(ctx, assign.Document); err != nil {

                  return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                }

                statusCode := ctx.Response().StatusCode()

                if !(200 <= statusCode && statusCode < 300) {

                  return nil
                }
              } else {

                assign.Document = ""

              }
            }

            if assign.Document == "" {

              assign.Document = filename

            } else {

              filename = assign.Document
            }

            if err = util.SwagSaveDocument(ctx, filename, func(filename string) error {

              assign.Document = filename

              return assignRepo.Update(assign, "id = ?", assign.ID)

            }); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            // fix swag error return

            statusCode := ctx.Response().StatusCode()

            if !(200 <= statusCode && statusCode < 300) {

              return nil
            }

            if err = assignRepo.Update(assign, "id = ?", assign.ID); err != nil {

              return ctx.InternalServerError(kornet.Msg(err.Error(), true))
            }

            var event *models.Events

            taskName := "resume"
            taskDescription := fmt.Sprintf("resume %s by %s", course.Name, userModel.Username)

            if event, err = eventRepo.Find("user_id = ? AND name = ? AND description = ?", userModel.ID, taskName, taskDescription); event != nil {

              // hack, make it fresh again
              currentTime := time.Now().UTC()
              event.CreatedAt = currentTime
              event.UpdatedAt = currentTime
              event.DeletedAt = gorm.DeletedAt{Time: currentTime, Valid: false}

              if err = eventRepo.Update(event, "id = ?", event.ID); err != nil {

                return err
              }

            } else {

              // create events
              if _, err = eventRepo.Create(&models.Events{
                Model:       &easy.Model{},
                UserID:      userModel.ID,
                Name:        taskName,
                Description: taskDescription,
              }); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }
            }

            return ctx.OK(kornet.Msg("success update resume", false))
          }

          return ctx.BadRequest(kornet.Msg("course not found", true))
        }

        return ctx.BadRequest(kornet.Msg("resume not found", true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })

  router.Get("/course/report", &m.KMap{
    "AuthToken":   true,
    "description": "Get Course Report",
    "request": &m.KMap{
      "params": &m.KMap{
        "id": "string", // course id
      },
      "headers": &m.KMap{
        "Authorization": "string",
      },
    },
    "response": swag.OkJSON(&kornet.Result{}),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var userModel *mo.UserModel
    var ok bool

    if ctx.Event() {

      if userModel, ok = ctx.Target().(*mo.UserModel); ok {

        kReq, _ := ctx.Kornet()
        courseId := m.KValueToString(kReq.Query.Get("id"))

        var completionCourses *models.CompletionCourses
        var completionModules []models.CompletionModules
        var course *models.Courses

        if completionCourses, err = completionCourseRepo.Find("user_id = ? AND course_id = ?", userModel.ID, courseId); err != nil {

          return ctx.BadRequest(kornet.Msg("you have not completed this course", true))
        }

        if course, err = courseRepo.Find("id = ?", courseId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        courseInfo := &m.KMap{
          "id":   course.ID,
          "name": course.Name,
        }

        if completionModules, err = completionModuleRepo.FindAll(-1, -1, "user_id = ? AND course_id = ?", userModel.ID, courseId); err != nil {

          return ctx.InternalServerError(kornet.Msg(err.Error(), true))
        }

        modules := make([]m.KMapImpl, len(completionModules))

        for i, completionModule := range completionModules {

          var module *models.Modules

          if module, err = moduleRepo.Unscoped().Find("id = ?", completionModule.ModuleID); err != nil {

            return ctx.InternalServerError(kornet.Msg(err.Error(), true))
          }

          moduleInfo := &m.KMap{
            "id":   module.ID,
            "name": module.Name,
          }

          modules[i] = &m.KMap{
            "id":     completionModule.ID,
            "module": moduleInfo,
            "score":  completionModule.Score,
          }
        }

        return ctx.OK(kornet.ResultNew(kornet.MessageNew("success get course grades", false), &m.KMap{
          "id":      completionCourses.ID,
          "course":  courseInfo,
          "score":   completionCourses.Score,
          "modules": modules,
        }))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

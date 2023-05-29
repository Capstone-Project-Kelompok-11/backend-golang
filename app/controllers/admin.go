package controllers

import (
	"encoding/json"
	"github.com/shopspring/decimal"
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

func AdminController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

	conn := pn.Connection()
	DB := conn.GORM()

	userRepo, _ := repository.UserRepositoryNew(DB)
	courseRepo, _ := repository.CourseRepositoryNew(DB)
	checkoutRepo, _ := repository.CheckoutRepositoryNew(DB)
	moduleRepo, _ := repository.ModuleRepositoryNew(DB)
	quizRepo, _ := repository.QuizzesRepositoryNew(DB)

	pp.Void(userRepo)

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

				if data, err = courseRepo.PreloadFindAllAndOrder(size, page, "name "+sort, "name LIKE ?", search); err != nil {

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
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		kReq, _ := ctx.Kornet()
		size := util.ValueToInt(kReq.Query.Get("size"))
		page := util.ValueToInt(kReq.Query.Get("page"))

		var err error
		var data []models.Checkout

		if data, err = checkoutRepo.CatchAll(size, page); err != nil {

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
			}
		}

		return ctx.OK(kornet.ResultNew(kornet.MessageNew("checkout history", false), exposed))
	})
}

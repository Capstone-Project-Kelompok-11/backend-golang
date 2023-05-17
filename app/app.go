package app

import (
  "lms/app/controllers"
  "lms/app/factory"
  "lms/app/middleware"
  "lms/app/models"
  "lms/app/tasks"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  bac "skfw/papaya/pigeon/templates/basicAuth/controllers"
  "time"
)

func App(pn papaya.NetImpl) error {

  conn := pn.Connection()
  DB := conn.GORM()

  ManageControlResourceShared(pn)

  logger := middleware.MakeLoggerMiddleware(pn)
  pn.Use(logger)

  swagger := pn.MakeSwagger(&swag.SwagInfo{
    Title:       "Academy API",
    Version:     "1.0.0",
    Description: "Academy API Documentation",
  })

  mainGroup := swagger.Group("/api/v1", "Schema")

  anonymGroup := mainGroup.Group("/public", "Anonymous")
  actionGroup := mainGroup.Group("/action", "Action")
  adminGroup := mainGroup.Group("/admin", "Administration")
  userGroup := mainGroup.Group("/users", "Authentication")

  anonymRouter := anonymGroup.Router()
  actionRouter := actionGroup.Router()
  adminRouter := adminGroup.Router()
  userRouter := userGroup.Router()

  expired := time.Hour * 24
  activeDuration := time.Hour * 4 // time to live, interval
  maxSessions := 6

  controllers.AnonymController(pn, anonymRouter)

  basicAuth := bac.BasicAuthNew(conn, expired, activeDuration, maxSessions)
  basicAuth.Bind(swagger, userRouter)

  controllers.ActionController(pn, actionRouter)
  controllers.AdminController(pn, adminRouter)

  swagger.AddTask(tasks.MakeAdminTask())

  DB.AutoMigrate(
    &models.Users{},
    &models.Sessions{},
    &models.Categories{},
    &models.Courses{},
    &models.CategoryCourses{},
    &models.Modules{},
    &models.Quizzes{},
    &models.ReviewCourses{},
    &models.CompletionCourses{},
    &models.CompletionModules{},
    &models.Transactions{},
    &models.Carts{},
  )

  factory.AdminFactory(pn) // set admin factory

  swagger.Start()

  return pn.Serve("127.0.0.1", 8000)
}

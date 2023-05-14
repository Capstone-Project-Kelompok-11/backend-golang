package app

import (
	"lms/app/controllers"
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
	gorm := conn.GORM()

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
	userGroup := mainGroup.Group("/users", "Authentication")

	anonymRouter := anonymGroup.Router()
	userRouter := userGroup.Router()

	expired := time.Hour * 24
	activeDuration := time.Hour * 4 // time to live, interval
	maxSessions := 6

	basicAuth := bac.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

	swagger.AddTask(tasks.MakeAdminTask())

	gorm.AutoMigrate(
		&models.Users{},
		&models.Sessions{},
		&models.Modules{},
		&models.Courses{},
		&models.Carts{},
		&models.Transactions{},
	)

	controllers.AnonymController(anonymRouter)

	swagger.Start()

	return pn.Serve("127.0.0.1", 8000)
}

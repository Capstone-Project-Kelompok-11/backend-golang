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

	pn.Use(middleware.MakeLoggerMiddleware(pn))

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Academy API",
		Version:     "1.0.0",
		Description: "Academy API Documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")

	testGroup := mainGroup.Group("/test", "Testing")
	anonymGroup := mainGroup.Group("/public", "Anonymous")
	adminGroup := mainGroup.Group("/admin", "Administration")
	userGroup := mainGroup.Group("/users", "User Management")

	testRouter := testGroup.Router()
	anonymRouter := anonymGroup.Router()
	adminRouter := adminGroup.Router()
	userRouter := userGroup.Router()

	expired := time.Hour * 24
	activeDuration := time.Minute * 24 // time to live, interval
	maxSessions := 24

	controllers.TestController(pn, testRouter)
	controllers.AnonymController(pn, anonymRouter)

	basicAuth := bac.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

	controllers.ActionController(pn, userRouter)
	controllers.CourseController(pn, userRouter)
	controllers.CheckoutController(pn, userRouter)
	controllers.AdminController(pn, adminRouter)

	swagger.AddTask(tasks.MakeAdminTask())

	DB.AutoMigrate(
		&models.Users{},
		&models.Events{},
		&models.Sessions{},
		&models.Categories{},
		&models.Courses{},
		&models.CategoryCourses{},
		&models.Assignments{},
		&models.Modules{},
		&models.Quizzes{},
		&models.ReviewCourses{},
		&models.CompletionCourses{},
		&models.CompletionModules{},
		&models.Checkout{},
		&models.Banners{},
	)

	factory.AdminFactory(pn) // set admin factory

	swagger.Start()

	return pn.Serve("0.0.0.0", 8080)
}

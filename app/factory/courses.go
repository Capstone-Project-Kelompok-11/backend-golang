package factory

import (
	"github.com/shopspring/decimal"
	"lms/app/models"
	"lms/app/repository"
	"skfw/papaya"
	"skfw/papaya/pigeon/easy"
)

func CourseFactory(pn papaya.NetImpl) error {

	var err error

	conn := pn.Connection()
	DB := conn.GORM()

	userRepo, _ := repository.UserRepositoryNew(DB)
	courseRepo, _ := repository.CourseRepositoryNew(DB)
	moduleRepo, _ := repository.ModuleRepositoryNew(DB)
	quizRepo, _ := repository.QuizzesRepositoryNew(DB)
	catRepo, _ := repository.CategoryRepositoryNew(DB)
	catCourseRepo, _ := repository.CategoryCourseRepositoryNew(DB)

	var user *models.Users

	if user, err = userRepo.Find("username = ?", "admin"); err != nil {
		panic("cannot find user")
	}

	var course *models.Courses

	course, _ = courseRepo.Find("name = ?", "example")

	if course != nil {
		return nil
	}

	if course, err = courseRepo.Create(&models.Courses{
		Model:       &easy.Model{},
		Name:        "example",
		Description: "this is example course",
		Thumbnail:   "example.png",
		UserID:      user.ID,
		Price:       decimal.NewFromInt(2_000_000),
		Level:       "Moderate",
	}); err != nil {
		panic("cannot create course")
	}

	var category *models.Categories

	if category, err = catRepo.Create(&models.Categories{
		Model:       &easy.Model{},
		Name:        "example",
		Description: "this is example category",
		Thumbnail:   "example.png",
	}); err != nil {
		panic("cannot create category")
	}

	catCourseRepo.Create(&models.CategoryCourses{
		Model:      &easy.Model{},
		CategoryID: category.ID,
		CourseID:   course.ID,
	})

	var module *models.Modules

	if module, err = moduleRepo.Create(&models.Modules{
		Model:       &easy.Model{},
		CourseID:    course.ID,
		Name:        "Mathematics",
		Description: "this is mathematics",
	}); err != nil {
		panic("cannot create module")
	}

	var quiz *models.Quizzes

	if quiz, err = quizRepo.Create(&models.Quizzes{
		Model:    &easy.Model{},
		ModuleID: module.ID,
		Data:     " [{ \"question\": \"1 + 2 = ?\", \"choices\": [ { \"text\": \"3\", \"valid\": true }, { \"text\": \"12\", \"valid\": false } ] }]",
	}); err != nil {
		panic("cannot create quiz")
	}

	if quiz == nil {
		panic("cannot create quiz")
	}

	return err
}

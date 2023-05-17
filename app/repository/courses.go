package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CourseRepository struct {
  Repository easy.RepositoryImpl[models.Courses]
}

type CourseRepositoryImpl interface {
  easy.RepositoryImpl[models.Courses]
}

func CourseRepositoryNew(DB *gorm.DB) (CourseRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Courses]
  if repo, err = easy.RepositoryNew[models.Courses](DB, &models.Courses{}); err != nil {

    return nil, err
  }
  courseRepo := &CourseRepository{
    Repository: repo,
  }
  return courseRepo, nil
}

func (c *CourseRepository) Init(DB *gorm.DB, model *models.Courses) error {

  return c.Repository.Init(DB, model)
}

func (c *CourseRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CourseRepository) Find(query any, args ...any) (*models.Courses, error) {

  return c.Repository.Find(query, args...)
}

func (c *CourseRepository) FindAll(size int, page int, query any, args ...any) ([]models.Courses, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CourseRepository) CatchAll(size int, page int) ([]models.Courses, error) {

  return c.Repository.CatchAll(size, page)
}

func (c *CourseRepository) Create(model *models.Courses) (*models.Courses, error) {

  return c.Repository.Create(model)
}

func (c *CourseRepository) Update(model *models.Courses, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CourseRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CourseRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

func (c *CourseRepository) Unscoped() easy.RepositoryImpl[models.Courses] {

  return c.Unscoped()
}

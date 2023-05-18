package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CategoryCourseRepository struct {
  Repository easy.RepositoryImpl[models.CategoryCourses]
}

type CategoryCourseRepositoryImpl interface {
  easy.RepositoryImpl[models.CategoryCourses]
}

func CategoryCourseRepositoryNew(DB *gorm.DB) (CategoryCourseRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.CategoryCourses]
  if repo, err = easy.RepositoryNew[models.CategoryCourses](DB, &models.CategoryCourses{}); err != nil {

    return nil, err
  }
  cateCourseRepo := &CategoryCourseRepository{
    Repository: repo,
  }
  return cateCourseRepo, nil
}

func (c *CategoryCourseRepository) Init(DB *gorm.DB, model *models.CategoryCourses) error {

  return c.Repository.Init(DB, model)
}

func (c *CategoryCourseRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CategoryCourseRepository) Find(query any, args ...any) (*models.CategoryCourses, error) {

  return c.Repository.Find(query, args...)
}

func (c *CategoryCourseRepository) FindAll(size int, page int, query any, args ...any) ([]models.CategoryCourses, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CategoryCourseRepository) CatchAll(size int, page int) ([]models.CategoryCourses, error) {

  return c.Repository.CatchAll(size, page)
}

func (c *CategoryCourseRepository) Create(model *models.CategoryCourses) (*models.CategoryCourses, error) {

  return c.Repository.Create(model)
}

func (c *CategoryCourseRepository) Update(model *models.CategoryCourses, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CategoryCourseRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CategoryCourseRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

func (c *CategoryCourseRepository) Unscoped() easy.RepositoryImpl[models.CategoryCourses] {

  return c.Repository.Unscoped()
}

func (c *CategoryCourseRepository) GORM() *gorm.DB {

  return c.Repository.GORM()
}

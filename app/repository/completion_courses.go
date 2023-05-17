package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CompletionCourseRepository struct {
  Repository easy.RepositoryImpl[models.CompletionCourses]
}

type CompletionCourseRepositoryImpl interface {
  easy.RepositoryImpl[models.CompletionCourses]
}

func CompletionCourseRepositoryNew(DB *gorm.DB) (CompletionCourseRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.CompletionCourses]
  if repo, err = easy.RepositoryNew[models.CompletionCourses](DB, &models.CompletionCourses{}); err != nil {

    return nil, err
  }
  completionCourseRepo := &CompletionCourseRepository{
    Repository: repo,
  }
  return completionCourseRepo, nil
}

func (c *CompletionCourseRepository) Init(DB *gorm.DB, model *models.CompletionCourses) error {

  return c.Repository.Init(DB, model)
}

func (c *CompletionCourseRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CompletionCourseRepository) Find(query any, args ...any) (*models.CompletionCourses, error) {

  return c.Repository.Find(query, args...)
}

func (c *CompletionCourseRepository) FindAll(size int, page int, query any, args ...any) ([]models.CompletionCourses, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CompletionCourseRepository) Create(model *models.CompletionCourses) (*models.CompletionCourses, error) {

  return c.Repository.Create(model)
}

func (c *CompletionCourseRepository) Update(model *models.CompletionCourses, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CompletionCourseRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CompletionCourseRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

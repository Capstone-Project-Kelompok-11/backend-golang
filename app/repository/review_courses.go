package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type ReviewCourseRepository struct {
  Repository easy.RepositoryImpl[models.ReviewCourses]
}

type ReviewCourseRepositoryImpl interface {
  easy.RepositoryImpl[models.ReviewCourses]
}

func ReviewCourseRepositoryNew(DB *gorm.DB) (ReviewCourseRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.ReviewCourses]
  if repo, err = easy.RepositoryNew[models.ReviewCourses](DB, &models.ReviewCourses{}); err != nil {

    return nil, err
  }
  reviewCourseRepo := &ReviewCourseRepository{
    Repository: repo,
  }
  return reviewCourseRepo, nil
}

func (m *ReviewCourseRepository) Init(DB *gorm.DB, model *models.ReviewCourses) error {

  return m.Repository.Init(DB, model)
}

func (m *ReviewCourseRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *ReviewCourseRepository) Find(query any, args ...any) (*models.ReviewCourses, error) {

  return m.Repository.Find(query, args...)
}

func (m *ReviewCourseRepository) FindAll(size int, page int, query any, args ...any) ([]models.ReviewCourses, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *ReviewCourseRepository) CatchAll(size int, page int) ([]models.ReviewCourses, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *ReviewCourseRepository) Create(model *models.ReviewCourses) (*models.ReviewCourses, error) {

  return m.Repository.Create(model)
}

func (m *ReviewCourseRepository) Update(model *models.ReviewCourses, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *ReviewCourseRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *ReviewCourseRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *ReviewCourseRepository) Unscoped() easy.RepositoryImpl[models.ReviewCourses] {

  return m.Unscoped()
}

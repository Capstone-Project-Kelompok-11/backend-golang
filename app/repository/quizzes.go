package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type QuizzesRepository struct {
  Repository easy.RepositoryImpl[models.Quizzes]
}

type QuizzesRepositoryImpl interface {
  easy.RepositoryImpl[models.Quizzes]
}

func QuizzesRepositoryNew(DB *gorm.DB) (QuizzesRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Quizzes]
  if repo, err = easy.RepositoryNew[models.Quizzes](DB, &models.Quizzes{}); err != nil {

    return nil, err
  }
  quizRepo := &QuizzesRepository{
    Repository: repo,
  }
  return quizRepo, nil
}

func (m *QuizzesRepository) Init(DB *gorm.DB, model *models.Quizzes) error {

  return m.Repository.Init(DB, model)
}

func (m *QuizzesRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *QuizzesRepository) Find(query any, args ...any) (*models.Quizzes, error) {

  return m.Repository.Find(query, args...)
}

func (m *QuizzesRepository) FindAll(size int, page int, query any, args ...any) ([]models.Quizzes, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *QuizzesRepository) CatchAll(size int, page int) ([]models.Quizzes, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *QuizzesRepository) Create(model *models.Quizzes) (*models.Quizzes, error) {

  return m.Repository.Create(model)
}

func (m *QuizzesRepository) Update(model *models.Quizzes, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *QuizzesRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *QuizzesRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *QuizzesRepository) Unscoped() easy.RepositoryImpl[models.Quizzes] {

  return m.Repository.Unscoped()
}

func (m *QuizzesRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

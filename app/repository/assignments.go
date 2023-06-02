package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type AssignmentRepository struct {
  Repository easy.RepositoryImpl[models.Assignments]
}

type AssignmentRepositoryImpl interface {
  easy.RepositoryImpl[models.Assignments]
  Grade(id string, grade int) (*models.Assignments, error)
}

func AssignmentRepositoryNew(DB *gorm.DB) (AssignmentRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Assignments]
  if repo, err = easy.RepositoryNew[models.Assignments](DB, &models.Assignments{}); err != nil {

    return nil, err
  }
  assignmentRepo := &AssignmentRepository{
    Repository: repo,
  }
  return assignmentRepo, nil
}

func (m *AssignmentRepository) Init(DB *gorm.DB, model *models.Assignments) error {

  return m.Repository.Init(DB, model)
}

func (m *AssignmentRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *AssignmentRepository) Find(query any, args ...any) (*models.Assignments, error) {

  return m.Repository.Find(query, args...)
}

func (m *AssignmentRepository) FindAll(size int, page int, query any, args ...any) ([]models.Assignments, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *AssignmentRepository) CatchAll(size int, page int) ([]models.Assignments, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *AssignmentRepository) Create(model *models.Assignments) (*models.Assignments, error) {

  return m.Repository.Create(model)
}

func (m *AssignmentRepository) Update(model *models.Assignments, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *AssignmentRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *AssignmentRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *AssignmentRepository) Unscoped() easy.RepositoryImpl[models.Assignments] {

  return m.Repository.Unscoped()
}

func (m *AssignmentRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

func (m *AssignmentRepository) Grade(id string, grade int) (*models.Assignments, error) {

  var err error
  var assign *models.Assignments

  if assign, err = m.Find("id = ?", id); err != nil {

    return nil, err
  }

  assign.Score = grade

  return assign, m.Update(assign, "id = ?", id)
}

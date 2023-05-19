package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type EventRepository struct {
  Repository easy.RepositoryImpl[models.Events]
}

type EventRepositoryImpl interface {
  easy.RepositoryImpl[models.Events]
}

func EventRepositoryNew(DB *gorm.DB) (EventRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Events]
  if repo, err = easy.RepositoryNew[models.Events](DB, &models.Events{}); err != nil {

    return nil, err
  }
  eventRepo := &EventRepository{
    Repository: repo,
  }
  return eventRepo, nil
}

func (m *EventRepository) Init(DB *gorm.DB, model *models.Events) error {

  return m.Repository.Init(DB, model)
}

func (m *EventRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *EventRepository) Find(query any, args ...any) (*models.Events, error) {

  return m.Repository.Find(query, args...)
}

func (m *EventRepository) FindAll(size int, page int, query any, args ...any) ([]models.Events, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *EventRepository) CatchAll(size int, page int) ([]models.Events, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *EventRepository) Create(model *models.Events) (*models.Events, error) {

  return m.Repository.Create(model)
}

func (m *EventRepository) Update(model *models.Events, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *EventRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *EventRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *EventRepository) Unscoped() easy.RepositoryImpl[models.Events] {

  return m.Repository.Unscoped()
}

func (m *EventRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

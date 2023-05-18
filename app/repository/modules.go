package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type ModuleRepository struct {
  Repository easy.RepositoryImpl[models.Modules]
}

type ModuleRepositoryImpl interface {
  easy.RepositoryImpl[models.Modules]
}

func ModuleRepositoryNew(DB *gorm.DB) (ModuleRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Modules]
  if repo, err = easy.RepositoryNew[models.Modules](DB, &models.Modules{}); err != nil {

    return nil, err
  }
  moduleRepo := &ModuleRepository{
    Repository: repo,
  }
  return moduleRepo, nil
}

func (m *ModuleRepository) Init(DB *gorm.DB, model *models.Modules) error {

  return m.Repository.Init(DB, model)
}

func (m *ModuleRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *ModuleRepository) Find(query any, args ...any) (*models.Modules, error) {

  return m.Repository.Find(query, args...)
}

func (m *ModuleRepository) FindAll(size int, page int, query any, args ...any) ([]models.Modules, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *ModuleRepository) CatchAll(size int, page int) ([]models.Modules, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *ModuleRepository) Create(model *models.Modules) (*models.Modules, error) {

  return m.Repository.Create(model)
}

func (m *ModuleRepository) Update(model *models.Modules, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *ModuleRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *ModuleRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *ModuleRepository) Unscoped() easy.RepositoryImpl[models.Modules] {

  return m.Repository.Unscoped()
}

func (m *ModuleRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

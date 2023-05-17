package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CompletionModuleRepository struct {
  Repository easy.RepositoryImpl[models.CompletionModules]
}

type CompletionModuleRepositoryImpl interface {
  easy.RepositoryImpl[models.CompletionModules]
}

func CompletionModuleRepositoryNew(DB *gorm.DB) (CompletionModuleRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.CompletionModules]
  if repo, err = easy.RepositoryNew[models.CompletionModules](DB, &models.CompletionModules{}); err != nil {

    return nil, err
  }
  completionModuleRepo := &CompletionModuleRepository{
    Repository: repo,
  }
  return completionModuleRepo, nil
}

func (c *CompletionModuleRepository) Init(DB *gorm.DB, model *models.CompletionModules) error {

  return c.Repository.Init(DB, model)
}

func (c *CompletionModuleRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CompletionModuleRepository) Find(query any, args ...any) (*models.CompletionModules, error) {

  return c.Repository.Find(query, args...)
}

func (c *CompletionModuleRepository) FindAll(size int, page int, query any, args ...any) ([]models.CompletionModules, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CompletionModuleRepository) CatchAll(size int, page int) ([]models.CompletionModules, error) {

  return c.Repository.CatchAll(size, page)
}

func (c *CompletionModuleRepository) Create(model *models.CompletionModules) (*models.CompletionModules, error) {

  return c.Repository.Create(model)
}

func (c *CompletionModuleRepository) Update(model *models.CompletionModules, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CompletionModuleRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CompletionModuleRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

func (c *CompletionModuleRepository) Unscoped() easy.RepositoryImpl[models.CompletionModules] {

  return c.Unscoped()
}

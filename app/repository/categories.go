package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CategoryRepository struct {
  Repository easy.RepositoryImpl[models.Categories]
}

type CategoryRepositoryImpl interface {
  easy.RepositoryImpl[models.Categories]
}

func CategoryRepositoryNew(DB *gorm.DB) (CategoryRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Categories]
  if repo, err = easy.RepositoryNew[models.Categories](DB, &models.Categories{}); err != nil {

    return nil, err
  }
  cateRepo := &CategoryRepository{
    Repository: repo,
  }
  return cateRepo, nil
}

func (c *CategoryRepository) Init(DB *gorm.DB, model *models.Categories) error {

  return c.Repository.Init(DB, model)
}

func (c *CategoryRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CategoryRepository) Find(query any, args ...any) (*models.Categories, error) {

  return c.Repository.Find(query, args...)
}

func (c *CategoryRepository) FindAll(size int, page int, query any, args ...any) ([]models.Categories, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CategoryRepository) CatchAll(size int, page int) ([]models.Categories, error) {

  return c.Repository.CatchAll(size, page)
}

func (c *CategoryRepository) Create(model *models.Categories) (*models.Categories, error) {

  return c.Repository.Create(model)
}

func (c *CategoryRepository) Update(model *models.Categories, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CategoryRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CategoryRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

func (c *CategoryRepository) Unscoped() easy.RepositoryImpl[models.Categories] {

  return c.Repository.Unscoped()
}

func (c *CategoryRepository) GORM() *gorm.DB {

  return c.Repository.GORM()
}

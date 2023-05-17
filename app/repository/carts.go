package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CartRepository struct {
  Repository easy.RepositoryImpl[models.Carts]
}

type CartRepositoryImpl interface {
  easy.RepositoryImpl[models.Carts]
}

func CartRepositoryNew(DB *gorm.DB) (CartRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Carts]
  if repo, err = easy.RepositoryNew[models.Carts](DB, &models.Carts{}); err != nil {

    return nil, err
  }
  cartRepo := &CartRepository{
    Repository: repo,
  }
  return cartRepo, nil
}

func (c *CartRepository) Init(DB *gorm.DB, model *models.Carts) error {

  return c.Repository.Init(DB, model)
}

func (c *CartRepository) SessionNew() {

  c.Repository.SessionNew()
}

func (c *CartRepository) Find(query any, args ...any) (*models.Carts, error) {

  return c.Repository.Find(query, args...)
}

func (c *CartRepository) FindAll(size int, page int, query any, args ...any) ([]models.Carts, error) {

  return c.Repository.FindAll(size, page, query, args...)
}

func (c *CartRepository) CatchAll(size int, page int) ([]models.Carts, error) {

  return c.Repository.CatchAll(size, page)
}

func (c *CartRepository) Create(model *models.Carts) (*models.Carts, error) {

  return c.Repository.Create(model)
}

func (c *CartRepository) Update(model *models.Carts, query any, args ...any) error {

  return c.Repository.Update(model, query, args...)
}

func (c *CartRepository) Remove(query any, args ...any) error {

  return c.Repository.Remove(query, args...)
}

func (c *CartRepository) Delete(query any, args ...any) error {

  return c.Repository.Delete(query, args...)
}

func (c *CartRepository) Unscoped() easy.RepositoryImpl[models.Carts] {

  return c.Unscoped()
}

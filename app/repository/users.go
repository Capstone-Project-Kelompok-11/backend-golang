package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type UserRepository struct {
  Repository easy.RepositoryImpl[models.Users]
}

type UserRepositoryImpl interface {
  easy.RepositoryImpl[models.Users]
}

func UserRepositoryNew(DB *gorm.DB) (UserRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Users]
  if repo, err = easy.RepositoryNew[models.Users](DB, &models.Users{}); err != nil {

    return nil, err
  }
  userRepo := &UserRepository{
    Repository: repo,
  }
  return userRepo, nil
}

func (u *UserRepository) Init(DB *gorm.DB, model *models.Users) error {

  return u.Repository.Init(DB, model)
}

func (u *UserRepository) SessionNew() {

  u.Repository.SessionNew()
}

func (u *UserRepository) Find(query any, args ...any) (*models.Users, error) {

  return u.Repository.Find(query, args...)
}

func (u *UserRepository) FindAll(size int, page int, query any, args ...any) ([]models.Users, error) {

  return u.Repository.FindAll(size, page, query, args...)
}

func (u *UserRepository) Create(model *models.Users) (*models.Users, error) {

  return u.Repository.Create(model)
}

func (u *UserRepository) Update(model *models.Users, query any, args ...any) error {

  return u.Repository.Update(model, query, args...)
}

func (u *UserRepository) Remove(query any, args ...any) error {

  return u.Repository.Remove(query, args...)
}

func (u *UserRepository) Delete(query any, args ...any) error {

  return u.Repository.Delete(query, args...)
}

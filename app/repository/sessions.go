package repository

import (
  "gorm.io/gorm"
  "skfw/papaya/pigeon/easy"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

type SessionRepository struct {
  Repository easy.RepositoryImpl[mo.SessionModel]
}

type SessionRepositoryImpl interface {
  easy.RepositoryImpl[mo.SessionModel]
}

func SessionRepositoryNew(DB *gorm.DB) (SessionRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[mo.SessionModel]
  if repo, err = easy.RepositoryNew[mo.SessionModel](DB, &mo.SessionModel{}); err != nil {

    return nil, err
  }
  sessionRepo := &SessionRepository{
    Repository: repo,
  }
  return sessionRepo, nil
}

func (u *SessionRepository) Init(DB *gorm.DB, model *mo.SessionModel) error {

  return u.Repository.Init(DB, model)
}

func (u *SessionRepository) SessionNew() {

  u.Repository.SessionNew()
}

func (u *SessionRepository) Find(query any, args ...any) (*mo.SessionModel, error) {

  return u.Repository.Find(query, args...)
}

func (u *SessionRepository) FindAll(size int, page int, query any, args ...any) ([]mo.SessionModel, error) {

  return u.Repository.FindAll(size, page, query, args...)
}

func (u *SessionRepository) CatchAll(size int, page int) ([]mo.SessionModel, error) {

  return u.Repository.CatchAll(size, page)
}

func (u *SessionRepository) Create(model *mo.SessionModel) (*mo.SessionModel, error) {

  return u.Repository.Create(model)
}

func (u *SessionRepository) Update(model *mo.SessionModel, query any, args ...any) error {

  return u.Repository.Update(model, query, args...)
}

func (u *SessionRepository) Remove(query any, args ...any) error {

  return u.Repository.Remove(query, args...)
}

func (u *SessionRepository) Delete(query any, args ...any) error {

  return u.Repository.Delete(query, args...)
}

func (u *SessionRepository) Unscoped() easy.RepositoryImpl[mo.SessionModel] {

  return u.Repository.Unscoped()
}

func (u *SessionRepository) GORM() *gorm.DB {

  return u.Repository.GORM()
}

package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type TransactionRepository struct {
  Repository easy.RepositoryImpl[models.Transactions]
}

type TransactionRepositoryImpl interface {
  easy.RepositoryImpl[models.Transactions]
}

func TransactionRepositoryNew(DB *gorm.DB) (TransactionRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Transactions]
  if repo, err = easy.RepositoryNew[models.Transactions](DB, &models.Transactions{}); err != nil {

    return nil, err
  }
  transactionRepo := &TransactionRepository{
    Repository: repo,
  }
  return transactionRepo, nil
}

func (t *TransactionRepository) Init(DB *gorm.DB, model *models.Transactions) error {

  return t.Repository.Init(DB, model)
}

func (t *TransactionRepository) SessionNew() {

  t.Repository.SessionNew()
}

func (t *TransactionRepository) Find(query any, args ...any) (*models.Transactions, error) {

  return t.Repository.Find(query, args...)
}

func (t *TransactionRepository) FindAll(size int, page int, query any, args ...any) ([]models.Transactions, error) {

  return t.Repository.FindAll(size, page, query, args...)
}

func (t *TransactionRepository) Create(model *models.Transactions) (*models.Transactions, error) {

  return t.Repository.Create(model)
}

func (t *TransactionRepository) Update(model *models.Transactions, query any, args ...any) error {

  return t.Repository.Update(model, query, args...)
}

func (t *TransactionRepository) Remove(query any, args ...any) error {

  return t.Repository.Remove(query, args...)
}

func (t *TransactionRepository) Delete(query any, args ...any) error {

  return t.Repository.Delete(query, args...)
}

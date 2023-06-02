package repository

import (
  "errors"
  "fmt"
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type CheckoutRepository struct {
  Repository easy.RepositoryImpl[models.Checkout]
}

type CheckoutRepositoryImpl interface {
  easy.RepositoryImpl[models.Checkout]
  PreloadVerifyByUserId(id string, paymentMethod string) error
}

func CheckoutRepositoryNew(DB *gorm.DB) (CheckoutRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Checkout]
  if repo, err = easy.RepositoryNew[models.Checkout](DB, &models.Checkout{}); err != nil {

    return nil, err
  }
  checkoutRepo := &CheckoutRepository{
    Repository: repo,
  }
  return checkoutRepo, nil
}

func (m *CheckoutRepository) Init(DB *gorm.DB, model *models.Checkout) error {

  return m.Repository.Init(DB, model)
}

func (m *CheckoutRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *CheckoutRepository) Find(query any, args ...any) (*models.Checkout, error) {

  return m.Repository.Find(query, args...)
}

func (m *CheckoutRepository) FindAll(size int, page int, query any, args ...any) ([]models.Checkout, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *CheckoutRepository) CatchAll(size int, page int) ([]models.Checkout, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *CheckoutRepository) Create(model *models.Checkout) (*models.Checkout, error) {

  return m.Repository.Create(model)
}

func (m *CheckoutRepository) Update(model *models.Checkout, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *CheckoutRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *CheckoutRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *CheckoutRepository) Unscoped() easy.RepositoryImpl[models.Checkout] {

  return m.Repository.Unscoped()
}

func (m *CheckoutRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

func (m *CheckoutRepository) PreloadVerifyByUserId(id string, paymentMethod string) error {

  m.SessionNew()

  var err error

  if err = m.GORM().
    Where("user_id = ?", id).
    Updates(map[string]any{
      "payment_method": paymentMethod,
      "verify":         true,
    }).Error; err != nil {

    return errors.New(fmt.Sprintf("unable to catch courses"))
  }

  return nil
}

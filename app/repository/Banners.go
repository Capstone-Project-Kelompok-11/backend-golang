package repository

import (
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
)

type BannerRepository struct {
  Repository easy.RepositoryImpl[models.Banners]
}

type BannerRepositoryImpl interface {
  easy.RepositoryImpl[models.Banners]
}

func BannerRepositoryNew(DB *gorm.DB) (BannerRepositoryImpl, error) {

  var err error
  var repo easy.RepositoryImpl[models.Banners]
  if repo, err = easy.RepositoryNew[models.Banners](DB, &models.Banners{}); err != nil {

    return nil, err
  }
  bannerRepo := &BannerRepository{
    Repository: repo,
  }
  return bannerRepo, nil
}

func (m *BannerRepository) Init(DB *gorm.DB, model *models.Banners) error {

  return m.Repository.Init(DB, model)
}

func (m *BannerRepository) SessionNew() {

  m.Repository.SessionNew()
}

func (m *BannerRepository) Find(query any, args ...any) (*models.Banners, error) {

  return m.Repository.Find(query, args...)
}

func (m *BannerRepository) FindAll(size int, page int, query any, args ...any) ([]models.Banners, error) {

  return m.Repository.FindAll(size, page, query, args...)
}

func (m *BannerRepository) CatchAll(size int, page int) ([]models.Banners, error) {

  return m.Repository.CatchAll(size, page)
}

func (m *BannerRepository) Create(model *models.Banners) (*models.Banners, error) {

  return m.Repository.Create(model)
}

func (m *BannerRepository) Update(model *models.Banners, query any, args ...any) error {

  return m.Repository.Update(model, query, args...)
}

func (m *BannerRepository) Remove(query any, args ...any) error {

  return m.Repository.Remove(query, args...)
}

func (m *BannerRepository) Delete(query any, args ...any) error {

  return m.Repository.Delete(query, args...)
}

func (m *BannerRepository) Unscoped() easy.RepositoryImpl[models.Banners] {

  return m.Repository.Unscoped()
}

func (m *BannerRepository) GORM() *gorm.DB {

  return m.Repository.GORM()
}

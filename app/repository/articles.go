package repository

import (
	"lms/app/models"
	"skfw/papaya/pigeon/easy"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	Repository easy.RepositoryImpl[models.Articles]
}

type ArticleRepositoryImpl interface {
	easy.RepositoryImpl[models.Articles]
}

func ArticleRepositoryNew(DB *gorm.DB) (ArticleRepositoryImpl, error) {

	var err error
	var repo easy.RepositoryImpl[models.Articles]
	if repo, err = easy.RepositoryNew[models.Articles](DB, &models.Articles{}); err != nil {

		return nil, err
	}
	articleRepo := &ArticleRepository{
		Repository: repo,
	}
	return articleRepo, nil
}

func (m *ArticleRepository) Init(DB *gorm.DB, model *models.Articles) error {

	return m.Repository.Init(DB, model)
}

func (m *ArticleRepository) SessionNew() {

	m.Repository.SessionNew()
}

func (m *ArticleRepository) Find(query any, args ...any) (*models.Articles, error) {

	return m.Repository.Find(query, args...)
}

func (m *ArticleRepository) FindAll(size int, page int, query any, args ...any) ([]models.Articles, error) {

	return m.Repository.FindAll(size, page, query, args...)
}

func (m *ArticleRepository) CatchAll(size int, page int) ([]models.Articles, error) {

	return m.Repository.CatchAll(size, page)
}

func (m *ArticleRepository) Create(model *models.Articles) (*models.Articles, error) {

	return m.Repository.Create(model)
}

func (m *ArticleRepository) Update(model *models.Articles, query any, args ...any) error {

	return m.Repository.Update(model, query, args...)
}

func (m *ArticleRepository) Remove(query any, args ...any) error {

	return m.Repository.Remove(query, args...)
}

func (m *ArticleRepository) Delete(query any, args ...any) error {

	return m.Repository.Delete(query, args...)
}

func (m *ArticleRepository) Unscoped() easy.RepositoryImpl[models.Articles] {

	return m.Repository.Unscoped()
}

func (m *ArticleRepository) GORM() *gorm.DB {

	return m.Repository.GORM()
}

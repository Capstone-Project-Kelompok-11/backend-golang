package repository

import (
  "database/sql"
  "errors"
  "fmt"
  "gorm.io/gorm"
  "lms/app/models"
  "skfw/papaya/pigeon/easy"
  "time"
)

type UserRepository struct {
  Repository easy.RepositoryImpl[models.Users]
}

type UserRepositoryImpl interface {
  easy.RepositoryImpl[models.Users]
  CatchAllCheckoutVerified(query any, args ...any) ([]models.Checkout, error)
  CatchAllCheckoutNonVerified(query any, args ...any) ([]models.Checkout, error)
  CatchAllCheckoutCancelled(query any, args ...any) ([]models.Checkout, error)
  CountUser(userCount *int64) error
  CountStudent(studentCount *int64) error
  CountNewStudent(newStudentCount *int64) error
  CountGraduate(graduateCount *int64) error
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

func (u *UserRepository) CatchAll(size int, page int) ([]models.Users, error) {

  return u.Repository.CatchAll(size, page)
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

func (u *UserRepository) Unscoped() easy.RepositoryImpl[models.Users] {

  return u.Repository.Unscoped()
}

func (u *UserRepository) CatchAllCheckoutVerified(query any, args ...any) ([]models.Checkout, error) {

  var err error

  var users []models.Users
  var checkouts []models.Checkout
  users = make([]models.Users, 0)
  checkouts = make([]models.Checkout, 0)

  DB := u.Repository.GORM()

  if err = DB.Preload("Checkout", "verify = ?", true).Where(query, args...).Find(&users).Error; err != nil {

    return checkouts, errors.New("unable to catch all checkout from users")
  }

  if len(users) > 0 {

    return users[0].Checkout, nil // checkouts
  }

  return checkouts, nil // empty checkouts
}

func (u *UserRepository) CatchAllCheckoutNonVerified(query any, args ...any) ([]models.Checkout, error) {

  var err error

  var users []models.Users
  var checkouts []models.Checkout
  users = make([]models.Users, 0)
  checkouts = make([]models.Checkout, 0)

  DB := u.Repository.GORM()

  if err = DB.Preload("Checkout", "verify = ?", false).Where(query, args...).Find(&users).Error; err != nil {

    return checkouts, errors.New("unable to catch all checkout from users")
  }

  if len(users) > 0 {

    return users[0].Checkout, nil // checkouts
  }

  return checkouts, nil // empty checkouts
}

func (u *UserRepository) CatchAllCheckoutCancelled(query any, args ...any) ([]models.Checkout, error) {

  var err error

  var users []models.Users
  var checkouts []models.Checkout
  users = make([]models.Users, 0)
  checkouts = make([]models.Checkout, 0)

  DB := u.Repository.GORM()

  if err = DB.Preload("Checkout", "verify = ? AND deleted_at IS NOT NULL", false).Where(query, args...).Find(&users).Error; err != nil {

    return checkouts, errors.New("unable to catch all checkout from users")
  }

  if len(users) > 0 {

    return users[0].Checkout, nil // checkouts
  }

  return checkouts, nil // empty checkouts
}

func (u *UserRepository) GORM() *gorm.DB {

  return u.Repository.GORM()
}

func (u *UserRepository) CountUser(userCount *int64) error {

  var err error
  var prepared *gorm.DB
  var row *sql.Row

  u.SessionNew()

  prepared = u.GORM().
    Select("COUNT(*) AS user").
    Where("admin = ?", false)
  if err = prepared.Error; err != nil {

    return errors.New("unable to catch user count")
  }

  if prepared != nil {
    if row = prepared.Row(); row != nil {
      if err = row.Scan(userCount); err != nil {

        if err != sql.ErrNoRows {

          return errors.New("unable to catch user count")
        }
        *userCount = 0

      }

      return nil
    }
  }

  return errors.New("unable to catch user count")
}

func (u *UserRepository) CountStudent(studentCount *int64) error {

  var err error
  var prepared *gorm.DB
  var row *sql.Row

  u.SessionNew()

  prepared = u.GORM().
    Select("COUNT(*) AS student").
    Where("admin = ?", false).
    Joins("INNER JOIN checkout ON checkout.user_id = users.id").
    Group("checkout.user_id")
  if err = prepared.Error; err != nil {

    return errors.New("unable to catch user count")
  }

  if prepared != nil {
    if row = prepared.Row(); row != nil {
      if err = row.Scan(studentCount); err != nil {

        if err != sql.ErrNoRows {

          return errors.New("unable to catch user count")
        }
        *studentCount = 0

      }

      return nil
    }
  }

  return errors.New("unable to catch student count")
}

func (u *UserRepository) CountNewStudent(newStudentCount *int64) error {

  var err error
  var prepared *gorm.DB
  var row *sql.Row

  u.SessionNew()

  currentTime := time.Now().UTC()

  beforeTime := currentTime
  afterTime := beforeTime.AddDate(0, 0, -7) // 1 week ago

  prepared = u.GORM().
    Select("COUNT(*) AS student").
    Where("admin = ?", false).
    Where("users.created_at BETWEEN ? AND ?", afterTime, beforeTime).
    Joins("INNER JOIN checkout ON checkout.user_id = users.id").
    Group("checkout.user_id")
  if err = prepared.Error; err != nil {

    return errors.New("unable to catch user count")
  }

  if prepared != nil {
    if row = prepared.Row(); row != nil {
      if err = row.Scan(newStudentCount); err != nil {

        if err != sql.ErrNoRows {

          return errors.New("unable to catch user count")
        }
        *newStudentCount = 0

      }

      return nil
    }
  }

  return errors.New("unable to catch new student count")
}

func (u *UserRepository) CountGraduate(graduateCount *int64) error {

  var err error
  var prepared *gorm.DB
  var row *sql.Row

  u.SessionNew()

  prepared = u.GORM().
    Select("COUNT(*) AS graduate").
    Joins("INNER JOIN completion_courses ON completion_courses.user_id = users.id").
    Group("completion_courses.user_id")
  if err = prepared.Error; err != nil {

    fmt.Println(err)

    return errors.New("unable to catch graduate count")
  }

  if prepared != nil {
    if row = prepared.Row(); row != nil {
      if err = row.Scan(graduateCount); err != nil {

        if err != sql.ErrNoRows {

          return errors.New("unable to catch graduate count")
        }
        *graduateCount = 0

      }
      return nil
    }
  }

  return errors.New("unable to catch graduate count")
}

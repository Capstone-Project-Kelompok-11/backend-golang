package factory

import (
  "github.com/shopspring/decimal"
  "lms/app/models"
  "lms/app/repository"
  "skfw/papaya"
  "skfw/papaya/pigeon/easy"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
  "skfw/papaya/pigeon/templates/basicAuth/util"
)

func AdminFactory(pn papaya.NetImpl) error {

  var err error

  conn := pn.Connection()
  DB := conn.GORM()

  var pass string
  userRepo, _ := repository.UserRepositoryNew(DB)

  if _, err = userRepo.Find("username = ?", "admin"); err != nil {

    if pass, _ = util.HashPassword("Admin@1234"); pass != "" {

      _, err = userRepo.Create(&models.Users{
        UserModel: &mo.UserModel{
          Model:    &easy.Model{},
          Username: "admin",
          Email:    "admin@mail.co",
          Password: pass,
          Admin:    true,
        },
        Balance: decimal.NewFromInt(500),
      })
    }

    if err != nil {

      pn.Logger().Error(err)
    }
  }

  return err
}

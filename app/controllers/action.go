package controllers

import (
  "lms/app/models"
  "lms/app/repository"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/koala/mapping"
  mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

func ActionController(pn papaya.NetImpl, router swag.SwagRouterImpl) {

  conn := pn.Connection()
  DB := conn.GORM()

  userRepo, _ := repository.UserRepositoryNew(DB)

  router.Get("/userinfo", &mapping.KMap{
    "AuthToken":   true,
    "description": "Catch User Information",
    "request":     nil,
    "responses": swag.OkJSON(&kornet.Result{
      Data: &mapping.KMap{
        "name":         "string",
        "username":     "string",
        "email":        "string",
        "gender":       "string",
        "phone":        "string",
        "dob":          "string",
        "address":      "string",
        "country_code": "string",
        "city":         "string",
        "postal_code":  "string",
        "verify":       "boolean",
        "admin":        "boolean",
        "balance":      "decimal",
      },
    }),
  }, func(ctx *swag.SwagContext) error {

    var err error
    var user *models.Users

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        // get full user information
        if user, err = userRepo.Find("id = ?", userModel.ID); user != nil {

          return ctx.OK(kornet.ResultNew(kornet.MessageNew("successful get user information", false), &mapping.KMap{
            "name":         user.Name.String,
            "username":     user.Username,
            "email":        user.Email,
            "gender":       user.Gender.String,
            "phone":        user.Phone.String,
            "dob":          user.DOB.Time,
            "address":      user.Address.String,
            "country_code": user.CountryCode.String,
            "city":         user.City.String,
            "postal_code":  user.PostalCode.String,
            "verify":       user.Verify,
            "admin":        user.Admin,
            "balance":      user.Balance.BigInt(),
          }))
        }

        return ctx.InternalServerError(kornet.Msg(err.Error(), true))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  })
}

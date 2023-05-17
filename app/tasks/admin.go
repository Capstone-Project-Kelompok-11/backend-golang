package tasks

import (
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/pigeon/templates/basicAuth/models"
)

func MakeAdminTask() *swag.SwagTask {

  return swag.MakeSwagTask("Admin", func(ctx *swag.SwagContext) error {

    if ctx.Event() {

      if user, ok := ctx.Target().(*models.UserModel); ok {

        if user.Admin {

          return nil
        }
      }
    }

    ctx.Prevent()
    return ctx.Unauthorized(kornet.Msg("Access Denied", true))
  })
}

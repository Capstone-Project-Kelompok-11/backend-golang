package middleware

import (
  "encoding/json"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "skfw/papaya"
  "skfw/papaya/koala"
)

func MakeLoggerMiddleware(pn papaya.NetImpl) fiber.Handler {

  return func(ctx *fiber.Ctx) error {

    data := make(map[string]interface{})
    json.Unmarshal(ctx.Body(), &data)
    fmt.Println(data)
    fmt.Println(string(ctx.Body()))

    pn.Logger().Log(
      pn.Logger().Text(fmt.Sprintf("[%s]", ctx.Method()), koala.ColorCyan, koala.ColorBlack, koala.StyleBold),
      pn.Logger().Text(ctx.BaseURL(), koala.ColorYellow, koala.ColorBlack, koala.StyleBold),
      pn.Logger().Text(ctx.OriginalURL(), koala.ColorYellow, koala.ColorBlack, koala.StyleBold),
    )

    return ctx.Next()
  }
}

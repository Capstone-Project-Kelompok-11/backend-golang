package util

import (
  "mime"
  "mime/multipart"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
)

func SwagSaveImageX256(ctx *swag.SwagContext, name string) error {

  var err error

  var form *multipart.Form
  var extensions []string
  var ext string

  if form, err = ctx.MultipartForm(); err != nil {

    return ctx.BadRequest(kornet.Msg("request has no multipart/form-data", true))
  }

  images := m.Keys([]string{
    "image/jpeg",
    "image/png",
  })

  var found bool

  found = false

  for k, h := range form.File {

    switch k {
    case "img", "image", "draw", "drawing":

      if len(h) > 0 {

        header := h[0]
        cTy := header.Header.Get("Content-Type")
        cTy, _ = kornet.KSafeContentTy(cTy)

        if name == "" {

          name = SafePathName(header.Filename)
        }

        if images.Contain(cTy) {

          if extensions, err = mime.ExtensionsByType(cTy); err != nil {

            return ctx.InternalServerError(kornet.Msg("unable to get name of extension", true))
          }

          n := len(extensions)

          if n > 0 {

            ext = ".png" // force use PNG formatter
            output := "assets/public/images/" + name + ext

            return SaveImageX256(k, ext, output)(ctx.Ctx)
          }
        }
      }

      found = true
      break
    }

    if found {

      break
    }
  }

  if !found {

    return ctx.BadRequest(kornet.Msg("image data not found or incorrectly applied key", true))
  }

  return ctx.InternalServerError(kornet.Msg("something wrong", true))
}

func SwagSaveImage(ctx *swag.SwagContext, name string) error {

  var err error

  var form *multipart.Form
  var extensions []string
  var ext string

  if form, err = ctx.MultipartForm(); err != nil {

    return ctx.BadRequest(kornet.Msg("request is not form-data", true))
  }

  images := m.Keys([]string{
    "image/jpeg",
    "image/png",
  })

  var found bool

  found = false

  for k, h := range form.File {

    switch k {
    case "img", "image", "draw", "drawing":

      if len(h) > 0 {

        header := h[0]
        cTy := header.Header.Get("Content-Type")
        cTy, _ = kornet.KSafeContentTy(cTy) // maybe some other value was embedded in the Content-Type header like 'size='

        if name == "" {

          name = SafePathName(header.Filename)
        }

        if images.Contain(cTy) {

          if extensions, err = mime.ExtensionsByType(cTy); err != nil {

            return ctx.InternalServerError(kornet.Msg("unable to get name of extension", true))
          }

          n := len(extensions)

          if n > 0 {

            ext = extensions[n-1] // the last thing maybe a good choice
            output := "assets/public/images/" + name + ext

            return SaveImage(k, ext, output)(ctx.Ctx)
          }
        }
      }

      found = true
      break
    }

    if found {

      break
    }
  }

  if !found {

    return ctx.BadRequest(kornet.Msg("image data not found or incorrectly applied key", true))
  }

  return ctx.InternalServerError(kornet.Msg("something wrong", true))
}

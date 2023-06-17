package util

import (
  "mime"
  "mime/multipart"
  "os"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func SwagSaveImageX256(ctx *swag.SwagContext, name string, catchFileNameCallback CatchFileNameCallback) error {

  var err error

  var form *multipart.Form
  var extensions []string
  var ext string
  var fileNameChange bool

  prefName := name
  name = RemoveExtensionFromFileName(name)

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

    fileNameChange = false

    switch k {
    case "img", "image", "draw", "drawing":

      if len(h) > 0 {

        header := h[0]
        cTy := header.Header.Get("Content-Type")
        cTy, _ = kornet.KSafeContentTy(cTy)

        if SourceNetOrEmpty(name) {

          name = header.Filename
          fileNameChange = true
        }

        name = SafePathName(name)

        if images.Contain(cTy) {

          if extensions, err = mime.ExtensionsByType(cTy); err != nil {

            return ctx.InternalServerError(kornet.Msg("unable to get name of extension", true))
          }

          n := len(extensions)

          if n > 0 {

            ext = ".png" // force use PNG formatter
            filename := name + ext
            dir := "assets/public/images/"
            output := posix.KPathNew(dir).JoinStr(filename)

            if fileNameChange {

              filename, output = GenUniqFileNameOutput(dir, filename)

              if catchFileNameCallback != nil {

                if err = catchFileNameCallback(filename); err != nil {

                  return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                }
              }
            } else {

              // maybe pref-extension was wrong
              if prefName != filename {

                if catchFileNameCallback != nil {

                  if err = catchFileNameCallback(filename); err != nil {

                    return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                  }
                }
              }
            }

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

func SwagSaveImage(ctx *swag.SwagContext, name string, catchFileNameCallback CatchFileNameCallback) error {

  var err error

  var form *multipart.Form
  var extensions []string
  var ext string
  var fileNameChange bool

  prefName := name
  name = RemoveExtensionFromFileName(name)

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

    fileNameChange = false

    switch k {
    case "img", "image", "draw", "drawing":

      if len(h) > 0 {

        header := h[0]
        cTy := header.Header.Get("Content-Type")
        cTy, _ = kornet.KSafeContentTy(cTy) // maybe some other value was embedded in the Content-Type header like 'size='

        if SourceNetOrEmpty(name) {

          name = header.Filename
          fileNameChange = true
        }

        name = SafePathName(name)

        if images.Contain(cTy) {

          if extensions, err = mime.ExtensionsByType(cTy); err != nil {

            return ctx.InternalServerError(kornet.Msg("unable to get name of extension", true))
          }

          n := len(extensions)

          if n > 0 {

            ext = extensions[n-1] // the last thing maybe a good choice
            filename := name + ext
            dir := "assets/public/images/"
            output := posix.KPathNew(dir).JoinStr(filename)

            if fileNameChange {

              filename, output = GenUniqFileNameOutput(dir, filename)

              if catchFileNameCallback != nil {

                if err = catchFileNameCallback(filename); err != nil {

                  return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                }
              }
            } else {

              // maybe pref-extension was wrong
              if prefName != filename {

                if catchFileNameCallback != nil {

                  if err = catchFileNameCallback(filename); err != nil {

                    return ctx.InternalServerError(kornet.Msg(err.Error(), true))
                  }
                }
              }
            }

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

func SwagRemoveImage(ctx *swag.SwagContext, filename string) error {

  var err error

  dir := "assets/public/images/"
  p := posix.KPathNew(dir).JoinStr(filename)
  f := kio.KFileNew(p)
  if f.IsExist() {

    if f.IsFile() {

      if err = os.Remove(p); err != nil {

        return ctx.InternalServerError(kornet.Msg("unable to remove image", true))
      }

      return ctx.OK(kornet.Msg("image has been removed", false))
    }

    return ctx.InternalServerError(kornet.Msg("something error", true))
  }

  return ctx.BadRequest(kornet.Msg("image already removed", true))
}

func SwagCheckImageExist(filename string) bool {

  dir := "assets/public/images/"
  p := posix.KPathNew(dir).JoinStr(filename)
  f := kio.KFileNew(p)
  if f.IsExist() {

    if f.IsFile() {

      return true
    }
  }

  return false
}

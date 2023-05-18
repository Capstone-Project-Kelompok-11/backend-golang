package util

import (
  "mime"
  "mime/multipart"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
)

func SwagSaveDocument(ctx *swag.SwagContext, name string, catchFileNameCallback CatchFileNameCallback) error {

  var err error

  var form *multipart.Form
  var extensions []string
  var ext string

  var fileNameChange bool

  if form, err = ctx.MultipartForm(); err != nil {

    return ctx.BadRequest(kornet.Msg("request is not form-data", true))
  }

  documents := m.Keys([]string{
    "text/plain",
    "text/html",
    "text/css",
    "text/javascript",
    "text/markdown",
    "application/rtf",
    "application/pdf",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    "application/vnd.ms-excel",
    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    "application/vnd.ms-powerpoint",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "application/msword",
  })

  var found bool

  found = false

  for k, h := range form.File {

    fileNameChange = false

    switch k {
    case "doc", "document":

      if len(h) > 0 {

        header := h[0]
        cTy := header.Header.Get("Content-Type")
        cTy, _ = kornet.KSafeContentTy(cTy) // maybe some other value was embedded in the Content-Type header like 'size='

        if name == "" {

          name = SafePathName(header.Filename)
          fileNameChange = true
        }

        if documents.Contain(cTy) {

          if extensions, err = mime.ExtensionsByType(cTy); err != nil {

            return ctx.InternalServerError(kornet.Msg("unable to get name of extension", true))
          }

          n := len(extensions)

          if n > 0 {

            ext = extensions[n-1] // the last thing maybe a good choice
            output := "assets/public/documents/" + name + ext

            if fileNameChange {

              if err = catchFileNameCallback(name + ext); err != nil {

                return ctx.InternalServerError(kornet.Msg(err.Error(), true))
              }
            }

            return SaveDocument(k, output)(ctx.Ctx)
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

    return ctx.BadRequest(kornet.Msg("document data not found or incorrectly applied key", true))
  }

  return ctx.InternalServerError(kornet.Msg("something wrong", true))
}

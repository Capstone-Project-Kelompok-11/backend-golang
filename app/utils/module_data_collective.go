package util

import (
  "lms/app/models"
  "net/url"
  "skfw/papaya/bunny/swag"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func ModuleDataCollective(ctx *swag.SwagContext, data []models.Modules) []m.KMapImpl {

  var err error
  res := make([]m.KMapImpl, 0)

  var URL *url.URL

  if URL, err = url.Parse(ctx.BaseURL()); err != nil {

    URL = &url.URL{}
  }

  imagePub := posix.KPathNew("/api/v1/public/image")
  documentPub := posix.KPathNew("/api/v1/public/document")
  videoPub := posix.KPathNew("/api/v1/public/video")

  for _, module := range data {

    if module.Thumbnail != "" {

      URL.Path = imagePub.Copy().JoinStr(module.Thumbnail)
      URL.RawPath = URL.Path

      module.Thumbnail = URL.String()
    }

    if module.Document != "" {

      URL.Path = documentPub.Copy().JoinStr(module.Document)
      URL.RawPath = URL.Path

      module.Document = URL.String()
    }

    if module.Video != "" {

      URL.Path = videoPub.Copy().JoinStr(module.Video)
      URL.RawPath = URL.Path

      module.Video = URL.String()
    }

    mm := &m.KMap{
      "name":        module.Name,
      "description": module.Description,
      "thumbnail":   module.Thumbnail,
      "video":       module.Video,
      "document":    module.Document,
      "created_at":  module.CreatedAt,
      "updated_at":  module.UpdatedAt,
    }

    if len(module.Quizzes) > 0 {

      mm.Put("quizzes", QuizzesDataCollective(ctx, module.Quizzes))
    }

    if len(module.CompletionModules) > 0 {

      mm.Put("completion_modules", CompletionModulesDataCollective(ctx, module.CompletionModules))
    }

    res = append(res, mm)
  }

  return res
}

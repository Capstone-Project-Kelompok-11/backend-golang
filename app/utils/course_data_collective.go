package util

import (
  "lms/app/models"
  "lms/app/repository"
  "net/url"
  "skfw/papaya/bunny/swag"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func CourseDataCollective(ctx *swag.SwagContext, userRepo repository.UserRepositoryImpl, data []models.Courses) []m.KMapImpl {

  var err error
  res := make([]m.KMapImpl, 0)

  var URL *url.URL

  if URL, err = url.Parse(ctx.BaseURL()); err != nil {

    URL = &url.URL{}
  }

  imagePub := posix.KPathNew("/api/v1/public/image")
  documentPub := posix.KPathNew("/api/v1/public/document")
  videoPub := posix.KPathNew("/api/v1/public/video")

  for _, course := range data {

    var user *models.Users

    if user, err = userRepo.Find("id = ?", course.UserID); err != nil {

      continue
    }

    if user.Image != "" {

      URL.Path = imagePub.Copy().JoinStr(user.Image)
      URL.RawPath = URL.Path

      user.Image = URL.String()
    }

    if course.Thumbnail != "" {

      URL.Path = imagePub.Copy().JoinStr(course.Thumbnail)
      URL.RawPath = URL.Path

      course.Thumbnail = URL.String()
    }

    if course.Document != "" {

      URL.Path = documentPub.Copy().JoinStr(course.Document)
      URL.RawPath = URL.Path

      course.Document = URL.String()
    }

    if course.Video != "" {

      URL.Path = videoPub.Copy().JoinStr(course.Video)
      URL.RawPath = URL.Path

      course.Video = URL.String()
    }

    mm := &m.KMap{
      "id": course.Model.ID,
      "create_by": &m.KMap{
        "name":     user.Name.String,
        "username": user.Username,
        "image":    user.Image,
      },
      "name":        course.Name,
      "description": course.Description,
      "thumbnail":   course.Thumbnail,
      "video":       course.Video,
      "document":    course.Document,
      "price":       course.Price.BigInt(),
      "level":       course.Level,
      "rating": RatingView(Rating{
        Rating1: course.Rating1,
        Rating2: course.Rating2,
        Rating3: course.Rating3,
        Rating4: course.Rating4,
        Rating5: course.Rating5,
      }),
      "finished":         course.Finished,
      "member_count":     course.MemberCount,
      "category_courses": course.CategoryCourses,
      "created_at":       course.CreatedAt,
      "update_at":        course.UpdatedAt,
    }

    if course.Modules != nil {

      mm.Put("modules", ModuleDataCollective(ctx, course.Modules))
    }

    if course.Reviews != nil {

      mm.Put("reviews", course.Reviews)
    }

    if course.Assignments != nil {

      mm.Put("assignments", course.Assignments)
    }

    if course.Checkout != nil {

      mm.Put("checkout", course.Checkout)
    }

    res = append(res, mm)
  }

  return res
}

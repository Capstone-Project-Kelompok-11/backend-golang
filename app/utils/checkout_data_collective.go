package util

import (
  "lms/app/models"
  "lms/app/repository"
  "net/url"
  "skfw/papaya/bunny/swag"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/tools/posix"
)

func CheckoutDataCollective(ctx *swag.SwagContext, userRepo repository.UserRepositoryImpl, courseRepo repository.CourseRepositoryImpl, data []models.Checkout) []m.KMapImpl {

  var err error
  res := make([]m.KMapImpl, 0)

  var URL *url.URL

  if URL, err = url.Parse(ctx.BaseURL()); err != nil {

    URL = &url.URL{}
  }

  imagePub := posix.KPathNew("/api/v1/public/image")

  for _, checkout := range data {

    var user *models.Users
    var course *models.Courses

    if user, _ = userRepo.Find("id = ?", checkout.UserID); user != nil {
      if course, _ = courseRepo.Find("id = ?", checkout.CourseID); course != nil {

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

        mm := &m.KMap{
          "id": checkout.Model.ID,
          "user": &m.KMap{
            "name":     user.Name.String,
            "username": user.Username,
            "image":    user.Image,
          },
          "course": &m.KMap{
            "id":        course.ID,
            "name":      course.Name,
            "thumbnail": course.Thumbnail,
            "price":     course.Price.BigInt(),
          },
          "paid":           checkout.Verify,
          "cancel":         checkout.DeletedAt.Valid,
          "payment_method": checkout.PaymentMethod,
          "created_at":     checkout.CreatedAt,
          "update_at":      checkout.UpdatedAt,
          "deleted_at":     checkout.DeletedAt,
        }

        res = append(res, mm)
      }
    }
  }

  return res
}

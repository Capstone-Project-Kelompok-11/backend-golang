package util

import (
  "lms/app/models"
  "lms/app/repository"
  m "skfw/papaya/koala/mapping"
)

func CheckoutDataCollective(userRepo repository.UserRepositoryImpl, courseRepo repository.CourseRepositoryImpl, data []models.Checkout) []m.KMapImpl {

  var err error
  res := make([]m.KMapImpl, 0)

  for _, checkout := range data {

    var user *models.Users
    var course *models.Courses

    if user, err = userRepo.Find("id = ?", checkout.UserID); err != nil {

      user = &models.Users{}
    }

    if course, err = courseRepo.Find("id = ?", checkout.CourseID); err != nil {

      course = &models.Courses{}
    }

    mm := &m.KMap{
      "id":   checkout.Model.ID,
      "data": checkout,
      "user": &m.KMap{
        "id":       user.ID,
        "name":     user.Name.String,
        "username": user.Username,
        "email":    user.Email,
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

  return res
}

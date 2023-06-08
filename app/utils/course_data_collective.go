package util

import (
  "lms/app/models"
  "lms/app/repository"
  m "skfw/papaya/koala/mapping"
)

func CourseDataCollective(userRepo repository.UserRepositoryImpl, data []models.Courses) []m.KMapImpl {

  var err error
  res := make([]m.KMapImpl, 0)

  for _, course := range data {

    var user *models.Users

    if user, err = userRepo.Find("id = ?", course.UserID); err != nil {

      continue
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

      mm.Put("modules", course.Modules)
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

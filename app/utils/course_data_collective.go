package util

import (
  "lms/app/models"
  m "skfw/papaya/koala/mapping"
)

func CourseDataCollective(data []models.Courses) []m.KMapImpl {

  res := make([]m.KMapImpl, 0)

  for _, course := range data {

    mm := &m.KMap{
      "id":          course.Model.ID,
      "name":        course.Name,
      "description": course.Description,
      "thumbnail":   course.Thumbnail,
      "video":       course.Video,
      "document":    course.Document,
      "price":       course.Price,
      "level":       course.Level,
      "rating": RatingView(Rating{
        Rating1: course.Rating1,
        Rating2: course.Rating2,
        Rating3: course.Rating3,
        Rating4: course.Rating4,
        Rating5: course.Rating5,
      }),
      "finished":     course.Finished,
      "member_count": course.MemberCount,
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

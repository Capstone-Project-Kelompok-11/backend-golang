package util

import (
  "lms/app/models"
  m "skfw/papaya/koala/mapping"
)

func CourseDataCollective(data []models.Courses) []m.KMapImpl {

  res := make([]m.KMapImpl, 0)

  for _, course := range data {

    res = append(res, &m.KMap{
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
      "finished":   course.Finished,
      "user_count": course.UserCount,
    })
  }

  return res
}

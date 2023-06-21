package util

import (
  "lms/app/models"
  "lms/app/repository"
  "skfw/papaya/koala/mapping"
)

func CategoryDataCollective(categoryRepo repository.CategoryRepositoryImpl, categories []string, exposed []mapping.KMapImpl) []mapping.KMapImpl {

  reduced := make([]mapping.KMapImpl, 0)

  for _, course := range exposed {

    categoriesCourse := make([]string, 0)

    if categoryCourseModels, ok := course.Get("category_courses").([]models.CategoryCourses); ok {

      for _, categoryCourseModel := range categoryCourseModels {

        if categoryModel, _ := categoryRepo.Find("id", categoryCourseModel.CategoryID); categoryModel != nil {

          categoriesCourse = append(categoriesCourse, categoryModel.Name)
        }
      }
    }

    course.Put("categories", categoriesCourse)
    course.Del("category_courses")

    categoryIncluded := true             // always true, maybe category comma separator is empty
    for _, context := range categories { // can handle none if category comma separator is empty

      if context == "all" {

        categoryIncluded = true
        break
      }

      found := false
      for _, categoryCourse := range categoriesCourse {

        if categoryCourse == context {

          found = true
          break
        }
      }

      if !found {

        categoryIncluded = false
        break
      }
    }

    if categoryIncluded {

      reduced = append(reduced, course)
    }
  }

  return reduced
}

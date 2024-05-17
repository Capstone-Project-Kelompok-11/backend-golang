package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"lms/app/models"
	"skfw/papaya/pigeon/easy"
)

type CourseRepository struct {
	Repository easy.RepositoryImpl[models.Courses]
}

type CourseRepositoryImpl interface {
	easy.RepositoryImpl[models.Courses]
	FindAllAndOrder(size int, page int, sort string, query any, args ...any) ([]models.Courses, error)
	FindAllAndOrderPopular(size int, page int, sort string, query any, args ...any) ([]models.Courses, error)
	PreFindAllAndOrder(size int, page int, sort string, query any, args ...any) ([]models.Courses, error)
	PreFind(query any, args ...any) (*models.Courses, error)
	UpdateMemberCountById(id string) error
	PreFindByCheckUserAndCourseId(userId string, courseId string) (*models.Courses, error)
	PreCatchAll(size int, page int) ([]models.Courses, error)
	PreCatchAllPopular(size int, page int) ([]models.Courses, error)
	CountCourse(courseCount *int64) error
	CatchIncomeAndOrders(args ...any) error
}

func CourseRepositoryNew(DB *gorm.DB) (CourseRepositoryImpl, error) {

	var err error
	var repo easy.RepositoryImpl[models.Courses]
	if repo, err = easy.RepositoryNew[models.Courses](DB, &models.Courses{}); err != nil {

		return nil, err
	}
	courseRepo := &CourseRepository{
		Repository: repo,
	}
	return courseRepo, nil
}

func (c *CourseRepository) Init(DB *gorm.DB, model *models.Courses) error {

	return c.Repository.Init(DB, model)
}

func (c *CourseRepository) SessionNew() {

	c.Repository.SessionNew()
}

func (c *CourseRepository) Find(query any, args ...any) (*models.Courses, error) {

	return c.Repository.Find(query, args...)
}

func (c *CourseRepository) FindAll(size int, page int, query any, args ...any) ([]models.Courses, error) {

	return c.Repository.FindAll(size, page, query, args...)
}

func (c *CourseRepository) CatchAll(size int, page int) ([]models.Courses, error) {

	return c.Repository.CatchAll(size, page)
}

func (c *CourseRepository) Create(model *models.Courses) (*models.Courses, error) {

	return c.Repository.Create(model)
}

func (c *CourseRepository) Update(model *models.Courses, query any, args ...any) error {

	return c.Repository.Update(model, query, args...)
}

func (c *CourseRepository) Remove(query any, args ...any) error {

	return c.Repository.Remove(query, args...)
}

func (c *CourseRepository) Delete(query any, args ...any) error {

	return c.Repository.Delete(query, args...)
}

func (c *CourseRepository) Unscoped() easy.RepositoryImpl[models.Courses] {

	return c.Repository.Unscoped()
}

func (c *CourseRepository) GORM() *gorm.DB {

	return c.Repository.GORM()
}

func (c *CourseRepository) PreCatchAll(size int, page int) ([]models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if page > 0 {

		offset := size * (page - 1)
		limit := size

		if err = c.GORM().
			Preload("CategoryCourses").
			Offset(offset).
			Limit(limit).
			Find(&data).
			Error; err != nil {

			return data, errors.New(fmt.Sprintf("unable to catch courses"))
		}
	}

	return data, nil
}

func (c *CourseRepository) PreCatchAllPopular(size int, page int) ([]models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if page > 0 {

		offset := size * (page - 1)
		limit := size

		if err = c.GORM().
			Preload("CategoryCourses").
			Joins("INNER JOIN (SELECT courses.id AS id, (k / NULLIF(n, 0))::FLOAT AS rating FROM courses INNER JOIN (SELECT id, (rating1 + rating2 + rating3 + rating4 + rating5)::FLOAT AS n, (rating1 + (rating2 * 2) + (rating3 * 3) + (rating4 * 4) + (rating5 * 5))::FLOAT AS k FROM courses) AS modify ON courses.id = modify.id WHERE n >= 0 ORDER BY rating DESC, member_count DESC) AS popular ON courses.id = popular.id").
			Offset(offset).
			Limit(limit).
			Find(&data).
			Error; err != nil {

			return data, errors.New(fmt.Sprintf("unable to catch courses"))
		}
	}

	return data, nil
}

func (c *CourseRepository) FindAllAndOrder(size int, page int, sort string, query any, args ...any) ([]models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if page > 0 {

		offset := size * (page - 1)
		limit := size

		if err = c.GORM().
			Preload("CategoryCourses").
			Where(query, args...).
			Order(sort).
			Offset(offset).
			Limit(limit).
			Find(&data).
			Error; err != nil {

			return data, errors.New(fmt.Sprintf("unable to catch courses"))
		}
	}

	return data, nil
}

func (c *CourseRepository) FindAllAndOrderPopular(size int, page int, sort string, query any, args ...any) ([]models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if page > 0 {

		offset := size * (page - 1)
		limit := size

		if err = c.GORM().
			Preload("CategoryCourses").
			Joins("INNER JOIN (SELECT courses.id AS id, (k / n)::FLOAT AS rating FROM courses INNER JOIN (SELECT id, (rating1 + rating2 + rating3 + rating4 + rating5)::FLOAT AS n, (rating1 + (rating2 * 2) + (rating3 * 3) + (rating4 * 4) + (rating5 * 5))::FLOAT AS k FROM courses) AS modify ON courses.id = modify.id WHERE n > 0 ORDER BY rating DESC, member_count DESC) AS popular ON courses.id = popular.id").
			Where(query, args...).
			Order(sort).
			Offset(offset).
			Limit(limit).
			Find(&data).
			Error; err != nil {

			return data, errors.New(fmt.Sprintf("unable to catch courses"))
		}
	}

	return data, nil
}

func (c *CourseRepository) PreFindAllAndOrder(size int, page int, sort string, query any, args ...any) ([]models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if page > 0 {

		offset := size * (page - 1)
		limit := size

		if err = c.GORM().
			Preload("Checkout").
			Preload("Modules").
			Preload("CategoryCourses").
			//Joins("INNER JOIN checkout ON courses.id = checkout.course_id").
			//Joins("INNER JOIN modules ON courses.id = modules.course_id").
			Where(query, args...).
			Order(sort).
			Offset(offset).
			Limit(limit).
			Find(&data).
			Error; err != nil {

			return data, errors.New("unable to catch courses")
		}
	}

	return data, nil
}

func (c *CourseRepository) PreFind(query any, args ...any) (*models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if err = c.GORM().
		Preload("Checkout").
		Preload("Modules").
		Preload("CategoryCourses").
		//Joins("INNER JOIN checkout ON courses.id = checkout.course_id").
		//Joins("INNER JOIN modules ON courses.id = modules.course_id").
		Where(query, args...).
		Find(&data).
		Error; err != nil {

		return nil, errors.New(fmt.Sprintf("unable to catch courses"))
	}

	if len(data) > 0 {

		return &data[0], nil
	}

	return nil, errors.New("course is empty")
}

func (c *CourseRepository) PreFindByCheckUserAndCourseId(userId string, courseId string) (*models.Courses, error) {

	c.SessionNew()

	var err error

	data := make([]models.Courses, 0)

	if err = c.GORM().
		Preload("Checkout", "id IN (?)", c.GORM().
			Table("checkout").
			Select("id").
			Where("verify = true AND user_id = ? AND course_id = ?", userId, courseId)).
		Preload("Modules").
		Preload("CategoryCourses").
		Where("id = ?", courseId).
		Find(&data).
		Error; err != nil {

		return nil, errors.New(fmt.Sprintf("unable to catch courses"))
	}

	if len(data) > 0 {

		return &data[0], nil
	}

	return nil, errors.New("course is empty")
}

func (c *CourseRepository) UpdateMemberCountById(id string) error {

	var err error
	c.SessionNew()

	var course *models.Courses
	memberCount := 0

	if course, err = c.Find("id = ?", id); err != nil {

		return err
	}

	memberCount = course.MemberCount + 1

	if err = c.GORM().
		Where("id = ?", id).
		Updates(map[string]interface{}{"member_count": memberCount}).
		Error; err != nil {

		return errors.New("unable to update member count")
	}

	return nil
}

func (c *CourseRepository) CountCourse(courseCount *int64) error {

	var err error
	var prepared *gorm.DB
	var row *sql.Row

	c.SessionNew()

	prepared = c.GORM().
		Select("COUNT(*) AS course")
	if err = prepared.Error; err != nil {

		return errors.New("unable to catch course count")
	}

	if prepared != nil {
		if row = prepared.Row(); row != nil {
			if err = row.Scan(courseCount); err != nil {

				if err != sql.ErrNoRows {

					return errors.New("unable to catch course count")
				}
				*courseCount = 0
			}

			return nil
		}
	}

	return errors.New("unable to catch course count")
}

func (c *CourseRepository) CatchIncomeAndOrders(args ...any) error {

	var err error

	var SQL *sql.DB
	var prepared *sql.Rows

	if SQL, err = c.GORM().DB(); err != nil {

		return err
	}

	if prepared, err = SQL.Query("SELECT SUM(CAST(price AS INT)) AS income, COUNT(name) AS orders FROM courses INNER JOIN checkout ON courses.id = checkout.course_id WHERE checkout.verify = true"); err != nil {

		return err
	}

	for prepared.Next() {

		if err = prepared.Scan(args...); err != nil {

			return err
		}
	}

	return nil
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseRepository interface {
	GetCourse(ctx context.Context, courseId string) (model.ICourse, error)
	ListCourses(ctx context.Context, filters model.CourseFilters, paginate Paginate) ([]model.ICourse, error)
}

type pgCourseRepository struct {
	db  *database.DB
	log logger.Logger
}

func NewPgCourseRepository(db *database.DB, log logger.Logger) *pgCourseRepository {
	return &pgCourseRepository{
		db:  db,
		log: log,
	}
}

func (c *pgCourseRepository) GetCourse(ctx context.Context, courseId string) (model.ICourse, error) {
	var course model.ICourse

	result, err := c.getCourseQuery(courseId)
	if err != nil {
		return course, err
	}

	err = c.db.Pool.QueryRow(context.TODO(), result.Query, result.Args...).Scan(
		&course.Id, &course.CategoryId, &course.Title, &course.Description, &course.Uri,
		&course.Course, &course.DateAdded,
	)
	if err != nil {
		c.log.Error("get course query error")
		return course, fmt.Errorf("failed to get course %d", courseId)
	}

	return course, nil
}

func (c *pgCourseRepository) ListCourses(ctx context.Context, filters model.CourseFilters, paginate Paginate) ([]model.ICourse, error) {
	var courses []model.ICourse

	result, err := c.listCoursesQuery(filters, paginate)
	if err != nil {
		return courses, err
	}

	rows, err := c.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		c.log.Error("list course query error")
		return courses, fmt.Errorf("failed to list courses")
	}

	for rows.Next() {
		course := model.ICourse{}
		err := rows.Scan(
			&course.Id, &course.CategoryId, &course.Title, &course.Description, &course.Uri,
			&course.Course, &course.DateAdded,
		)
		if err != nil {
			c.log.Error("scan course error")
			return courses, fmt.Errorf("failed to list course categories")
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (c *pgCourseRepository) getCourseQuery(courseId string) (util.SqlBuilderResult, error) {
	var courseUuid pgtype.UUID
	if err := courseUuid.Scan(courseId); err != nil {
		c.log.Warn("scan course id error")
		return util.SqlBuilderResult{}, fmt.Errorf("invalid course id %s", courseId)
	}

	qb := util.NewSqlBuilder(
		"select id, category_id, title, description, uri, course, date_added",
		"from courses",
	)
	qb = qb.Concat("where id = $%d", courseUuid)

	return qb.Result(), nil
}

func (c *pgCourseRepository) listCoursesQuery(filters model.CourseFilters, paginate Paginate) (util.SqlBuilderResult, error) {
	qb := util.NewSqlBuilder(
		"select id, category_id, title, description, uri, course, date_added",
		"from courses",
	)
	queryFilters := []string{}
	args := []any{}

	if filters.Name != "" {
		queryFilters = append(queryFilters, "name ilike $%d")
		args = append(args, "%"+filters.Name+"%")
	}

	if filters.CategoryId != "" {
		var categoryUUID pgtype.UUID
		if err := categoryUUID.Scan(filters.CategoryId); err != nil {
			c.log.Warn("invalid category id")
			return util.SqlBuilderResult{}, fmt.Errorf("invalid category id %s", filters.CategoryId)
		}
		queryFilters = append(queryFilters, "category_id = $%d")
		args = append(args, categoryUUID)
	}

	if filters.Course != "" {
		queryFilters = append(queryFilters, "course ilike $%d")
		args = append(args, filters.Course+"%")
	}

	if len(queryFilters) > 0 {
		orClause := "(" + strings.Join(queryFilters, " or ") + ")"
		qb = qb.Concat("where "+orClause, args...)
	}

	return qb.Result(), nil
}

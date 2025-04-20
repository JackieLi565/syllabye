package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseSchema struct {
	Id          string
	CategoryId  string
	Title       string
	Description sql.NullString
	Uri         string
	Course      string
	DateAdded   time.Time
}

type CourseFilters struct {
	Search     string
	CategoryId string
}

type CourseRepository interface {
	GetCourse(ctx context.Context, courseId string) (CourseSchema, error)
	ListCourses(ctx context.Context, filters CourseFilters, paginate util.Paginate) ([]CourseSchema, error)
}

type pgCourseRepository struct {
	db  *database.PostgresDb
	log logger.Logger
}

func NewPgCourseRepository(db *database.PostgresDb, log logger.Logger) *pgCourseRepository {
	return &pgCourseRepository{
		db:  db,
		log: log,
	}
}

func (c *pgCourseRepository) GetCourse(ctx context.Context, courseId string) (CourseSchema, error) {
	var course CourseSchema

	result, err := c.getCourseQuery(courseId)
	if err != nil {
		return course, err
	}

	err = c.db.Pool.QueryRow(context.TODO(), result.Query, result.Args...).Scan(
		&course.Id, &course.CategoryId, &course.Title, &course.Description, &course.Uri,
		&course.Course, &course.DateAdded,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return course, util.ErrNotFound
		}
		c.log.Error("get course query error", logger.Err(err))
		return course, util.ErrInternal
	}

	return course, nil
}

func (c *pgCourseRepository) ListCourses(ctx context.Context, filters CourseFilters, paginate util.Paginate) ([]CourseSchema, error) {
	var courses []CourseSchema

	result, err := c.listCoursesQuery(filters, paginate)
	if err != nil {
		return courses, err
	}

	rows, err := c.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return courses, nil
		}
		c.log.Error("list course query error", logger.Err(err))
		return courses, util.ErrInternal
	}

	for rows.Next() {
		course := CourseSchema{}
		err := rows.Scan(
			&course.Id, &course.CategoryId, &course.Title, &course.Description, &course.Uri,
			&course.Course, &course.DateAdded,
		)
		if err != nil {
			c.log.Error("scan course error", logger.Err(err))
			return courses, util.ErrInternal
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

func (c *pgCourseRepository) listCoursesQuery(filters CourseFilters, paginate util.Paginate) (util.SqlBuilderResult, error) {
	qb := util.NewSqlBuilder(
		"select id, category_id, title, description, uri, course, date_added",
		"from courses",
		"where 1 = 1",
	)

	if filters.CategoryId != "" {
		qb.Concat("and category_id = $%d", filters.CategoryId)
	}
	if filters.Search != "" {
		qb.Concat("and course ilike $%d or title ilike $%d", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	qb.Concat("limit $%d", paginate.Size)
	offset := (paginate.Page - 1) * paginate.Size
	qb.Concat("offset $%d", offset)

	return qb.Result(), nil // TODO: remove non-existent errors from API
}

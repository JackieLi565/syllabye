package repository

import (
	"context"
	"fmt"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseCategoryRepository interface {
	GetCourseCategory(ctx context.Context, categoryId string) (model.ICourseCategory, error)
	// Not need for pagination since dataset is very small
	ListCourseCategories(ctx context.Context, nameFilter string) ([]model.ICourseCategory, error)
}

type pgCourseCategoryRepository struct {
	db  *database.DB
	log logger.Logger
}

func NewPgCourseCategoryRepository(db *database.DB, log logger.Logger) *pgCourseCategoryRepository {
	return &pgCourseCategoryRepository{
		db:  db,
		log: log,
	}
}

func (cc *pgCourseCategoryRepository) GetCourseCategory(ctx context.Context, categoryId string) (model.ICourseCategory, error) {
	var category model.ICourseCategory

	result, err := cc.getCourseCategoryQuery(categoryId)
	if err != nil {
		return category, err
	}

	err = cc.db.Pool.QueryRow(context.TODO(), result.Query, result.Args...).Scan(
		&category.Id,
		&category.Name,
		&category.DateAdded,
	)
	if err != nil {
		cc.log.Error("get category query error")
		return category, fmt.Errorf("failed to get course category %d", categoryId)
	}

	return category, nil
}

func (cc *pgCourseCategoryRepository) ListCourseCategories(ctx context.Context, nameFilter string) ([]model.ICourseCategory, error) {
	var categories []model.ICourseCategory

	result := cc.listCourseCategoriesQuery(nameFilter)

	rows, err := cc.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		cc.log.Error("list category query error")
		return categories, fmt.Errorf("failed to list course categories")
	}

	for rows.Next() {
		category := model.ICourseCategory{}
		err := rows.Scan(&category.Id, &category.Name, &category.DateAdded)
		if err != nil {
			cc.log.Error("scan category error")
			return categories, fmt.Errorf("failed to list course categories")
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (cc *pgCourseCategoryRepository) getCourseCategoryQuery(categoryId string) (util.SqlBuilderResult, error) {
	var categoryUuid pgtype.UUID
	if err := categoryUuid.Scan(categoryId); err != nil {
		cc.log.Warn("scan category id error")
		return util.SqlBuilderResult{}, fmt.Errorf("invalid course category id %s", categoryId)
	}

	qb := util.NewSqlBuilder(
		"select id, name, date_added",
		"from course_categories",
	)
	qb = qb.Concat("where id = $%d", categoryUuid)

	return qb.Result(), nil
}

func (cc *pgCourseCategoryRepository) listCourseCategoriesQuery(nameFilter string) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select id, name, date_added",
		"from course_categories",
		"where 1 = 1",
	)

	if nameFilter != "" {
		qb = qb.Concat("and name ilike", "%"+nameFilter+"%")
	}

	return qb.Result()
}

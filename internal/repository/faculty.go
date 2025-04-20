package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type FacultySchema struct {
	Id        string
	Name      string
	DateAdded time.Time
}

type FacultyRepository interface {
	GetFaculty(ctx context.Context, facultyId string) (FacultySchema, error)
	// Not need for pagination since dataset is very small
	ListFaculties(ctx context.Context, nameFilter string) ([]FacultySchema, error)
}

type pgFacultyRepository struct {
	db  *database.PostgresDb
	log logger.Logger
}

func NewPgFacultyRepository(db *database.PostgresDb, log logger.Logger) *pgFacultyRepository {
	return &pgFacultyRepository{
		db:  db,
		log: log,
	}
}

func (f *pgFacultyRepository) GetFaculty(ctx context.Context, facultyId string) (FacultySchema, error) {
	var faculty FacultySchema

	result, err := f.getFacultyQuery(facultyId)

	err = f.db.Pool.QueryRow(context.TODO(), result.Query, result.Args...).Scan(
		&faculty.Id, &faculty.Name, &faculty.DateAdded,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return faculty, util.ErrNotFound
		}

		f.log.Error("get faculty query error", logger.Err(err))
		return faculty, fmt.Errorf("failed to get faculty %s", facultyId)
	}

	return faculty, nil
}

func (f *pgFacultyRepository) ListFaculties(ctx context.Context, nameFilter string) ([]FacultySchema, error) {
	var faculties []FacultySchema

	result := f.listFacultiesQuery(nameFilter)

	rows, err := f.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return faculties, nil
		}

		f.log.Error("list faculty query error", logger.Err(err))
		return faculties, fmt.Errorf("failed to list faculties")
	}

	for rows.Next() {
		faculty := FacultySchema{}
		err := rows.Scan(&faculty.Id, &faculty.Name, &faculty.DateAdded)
		if err != nil {
			f.log.Error("scan faculty query error", logger.Err(err))
			return faculties, fmt.Errorf("failed to list faculties")
		}

		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

func (f *pgFacultyRepository) getFacultyQuery(facultyId string) (util.SqlBuilderResult, error) {
	var facultyUuid pgtype.UUID
	if err := facultyUuid.Scan(facultyId); err != nil {
		return util.SqlBuilderResult{}, fmt.Errorf("invalid faculty %s id", facultyId)
	}

	qb := util.NewSqlBuilder(
		"select id, name, date_added",
		"from faculties",
	)
	qb = qb.Concat("where id = $%d", facultyUuid)

	return qb.Result(), nil
}

func (f *pgFacultyRepository) listFacultiesQuery(nameFilter string) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select id, name, date_added",
		"from faculties",
		"where 1 = 1",
	)

	if nameFilter != "" {
		qb = qb.Concat("and name ilike", "%"+nameFilter+"%")
	}

	return qb.Result()
}

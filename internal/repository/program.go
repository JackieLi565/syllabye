package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProgramRepository interface {
	GetProgram(ctx context.Context, programId string) (model.IProgram, error)
	// Not need for pagination since dataset is very small
	ListPrograms(ctx context.Context, filters model.ProgramFilters) ([]model.IProgram, error)
}

type pgProgramRepository struct {
	db  *database.DB
	log logger.Logger
}

func NewPgProgramRepository(db *database.DB, log logger.Logger) *pgProgramRepository {
	return &pgProgramRepository{
		db:  db,
		log: log,
	}
}

func (p *pgProgramRepository) GetProgram(ctx context.Context, programId string) (model.IProgram, error) {
	var program model.IProgram

	res, err := p.getProgramQuery(programId)
	if err != nil {
		return program, err
	}

	err = p.db.Pool.QueryRow(context.TODO(), res.Query, res.Args...).Scan(
		&program.Id, &program.FacultyId, &program.Name, &program.Uri, &program.DateAdded,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return program, util.ErrNotFound
		}

		p.log.Error("failed to get program", logger.Err(err))
		return program, fmt.Errorf("failed to get program %s", programId)
	}

	return program, nil
}

func (p *pgProgramRepository) ListPrograms(ctx context.Context, filters model.ProgramFilters) ([]model.IProgram, error) {
	var programs []model.IProgram

	res, err := p.listProgramsQuery(filters)
	if err != nil {
		return programs, err
	}

	rows, err := p.db.Pool.Query(context.TODO(), res.Query, res.Args...)
	if err != nil {
		p.log.Error("list programs query error")
		return programs, fmt.Errorf("failed to query programs")
	}

	for rows.Next() {
		program := model.IProgram{}
		err := rows.Scan(&program.Id, &program.FacultyId, &program.Name, &program.Uri, &program.DateAdded)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return programs, nil
			}

			p.log.Error("scan programs query error", logger.Err(err))
			return programs, fmt.Errorf("failed to scan program %s", program.Id)
		}

		programs = append(programs, program)
	}

	return programs, nil
}

func (p *pgProgramRepository) getProgramQuery(programId string) (util.SqlBuilderResult, error) {
	var programUuid pgtype.UUID
	if err := programUuid.Scan(programId); err != nil {
		return util.SqlBuilderResult{}, fmt.Errorf("invalid program uuid format %s", programId)
	}

	qb := util.NewSqlBuilder(
		"select id, faculty_id, name, uri, date_added",
		"from programs",
	)
	qb = qb.Concat("where id = $%d", programUuid)

	return qb.Result(), nil
}

func (p *pgProgramRepository) listProgramsQuery(filters model.ProgramFilters) (util.SqlBuilderResult, error) {
	qb := util.NewSqlBuilder(
		"select id, faculty_id, name, uri, date_added",
		"from programs",
		"where 1 = 1",
	)

	if filters.FacultyId != "" {
		var facultyUuid pgtype.UUID
		err := facultyUuid.Scan(filters.FacultyId)
		if err != nil {
			return util.SqlBuilderResult{}, fmt.Errorf("failed to decode faculty id %s", filters.FacultyId)
		}
		qb = qb.Concat("and faculty_id = $%d", facultyUuid)
	}

	if filters.Name != "" {
		qb.Concat("and name ilike $%d", "%"+filters.Name+"%")
	}

	return qb.Result(), nil
}

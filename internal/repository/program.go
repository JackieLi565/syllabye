package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type IProgramRepository interface {
	ListPrograms(filters ProgramFilters, paginate Paginate) ([]model.Program, error)
}

type ProgramRepository struct {
	DB *database.DB
}

type ProgramFilters struct {
	Faculty *pgtype.UUID
	Name    string
}

func (p ProgramRepository) ListPrograms(filters ProgramFilters, paginate Paginate) ([]model.Program, error) {
	qb := util.NewSqlBuilder(
		"select id, faculty_id, name, uri, date_added",
		"from programs",
	)

	if filters.Faculty != nil && filters.Faculty.Valid {
		qb = qb.Concat("and faculty_id = $%d", filters.Faculty.String())
	}

	if filters.Name != "" {
		qb.Concat("and name ilike $%d", "%"+filters.Name+"%")
	}

	paginate.parsePaginate()

	qb.Concat("limit $%d", paginate.size)
	offset := (paginate.page - 1) * paginate.size
	qb.Concat("offset $%d", offset)

	var programs []model.Program
	rows, err := p.DB.Pool.Query(context.Background(), qb.Build(), qb.GetArgs()...)
	if err != nil {
		log.Println(err)
		return programs, fmt.Errorf("failed to query programs")
	}

	for rows.Next() {
		var program model.Program
		err := rows.Scan(&program.Id, &program.FacultyId, &program.Name, &program.URI, &program.DateAdded)
		if err != nil {
			log.Println(err)
			return programs, fmt.Errorf("failed to scan program %s", program.Id)
		}

		programs = append(programs, program)
	}

	return programs, nil
}

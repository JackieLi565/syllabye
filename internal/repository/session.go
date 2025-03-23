package repository

import (
	"context"
	"errors"
	"time"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository interface {
	CreateSession(userId string, exp time.Time) (string, error)
	GetSession(sessionId string) (model.ISession, error)
}

type pgSessionRepository struct {
	db  *database.DB
	log logger.Logger
}

func NewPgSessionRepository(db *database.DB, log logger.Logger) *pgSessionRepository {
	return &pgSessionRepository{
		db:  db,
		log: log,
	}
}

func (s *pgSessionRepository) CreateSession(userId string, exp time.Time) (string, error) {
	var sessionId string
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return sessionId, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("insert into sessions (user_id, date_expires)")
	qb = qb.Concat("values ($%d, $%d)", userUuid, exp)
	qb = qb.Concat("returning id")

	err = s.db.Pool.QueryRow(context.TODO(), qb.Build(), qb.GetArgs()...).Scan(
		&sessionId,
	)
	if err != nil {
		if database.IsErrConflict(err) {
			return sessionId, util.ErrConflict
		}

		s.log.Error("failed to create session", logger.Err(err))
		return sessionId, util.ErrInternal
	}

	return sessionId, nil
}

func (s *pgSessionRepository) GetSession(sessionId string) (model.ISession, error) {
	var session model.ISession
	var sessionUuid pgtype.UUID
	err := sessionUuid.Scan(sessionId)
	if err != nil {
		return session, util.ErrMalformed
	}

	qb := util.NewSqlBuilder(
		"select id, user_id, date_added, date_expires",
		"from sessions",
	)
	qb = qb.Concat("where id = $%d", sessionUuid)

	err = s.db.Pool.QueryRow(context.TODO(), qb.Build(), qb.GetArgs()...).Scan(
		&session.Id,
		&session.UserId,
		&session.DateAdded,
		&session.DateExpires,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session, util.ErrNotFound
		}

		s.log.Error("get session query error", logger.Err(err))
		return session, util.ErrInternal
	}

	return session, nil
}

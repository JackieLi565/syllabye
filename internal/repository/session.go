package repository

import (
	"context"
	"errors"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userId string) (string, error)
	GetSession(ctx context.Context, sessionId string) (model.ISession, error)
}

type pgSessionRepository struct {
	db  *database.PostgresDb
	log logger.Logger
}

func NewPgSessionRepository(db *database.PostgresDb, log logger.Logger) *pgSessionRepository {
	return &pgSessionRepository{
		db:  db,
		log: log,
	}
}

func (s *pgSessionRepository) CreateSession(ctx context.Context, userId string) (string, error) {
	var sessionId string

	result, err := s.createSessionQuery(userId)
	if err != nil {
		return sessionId, err
	}

	err = s.db.Pool.QueryRow(ctx, result.Query, result.Args...).Scan(
		&sessionId,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == database.PgConflictErrCode {
			s.log.Error("session conflict query error")
			return sessionId, util.ErrConflict
		}

		s.log.Error("failed to create session", logger.Err(err))
		return sessionId, util.ErrInternal
	}

	return sessionId, nil
}

func (s *pgSessionRepository) GetSession(ctx context.Context, sessionId string) (model.ISession, error) {
	var session model.ISession
	var sessionUuid pgtype.UUID
	err := sessionUuid.Scan(sessionId)
	if err != nil {
		return session, util.ErrMalformed
	}

	qb := util.NewSqlBuilder(
		"select id, user_id, date_added",
		"from sessions",
	)
	qb = qb.Concat("where id = $%d", sessionUuid)

	err = s.db.Pool.QueryRow(ctx, qb.Build(), qb.GetArgs()...).Scan(
		&session.Id,
		&session.UserId,
		&session.DateAdded,
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

func (s *pgSessionRepository) createSessionQuery(userId string) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("insert into sessions (user_id)")
	qb = qb.Concat("values ($%d)", userUuid)
	qb = qb.Concat("returning id")

	return qb.Result(), nil
}

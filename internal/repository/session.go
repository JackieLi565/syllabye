package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository interface {
	CreateSession(userId string, exp time.Time) (model.Session, error)
	FindSession(sessionId string) (model.Session, error)
}

type PgSessionRepository struct {
	DB *database.DB
}

func (s PgSessionRepository) CreateSession(userId string, exp time.Time) (model.Session, error) {
	var session model.Session
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return session, fmt.Errorf("invalid user id")
	}

	qb := util.NewSqlBuilder("insert into sessions (user_id, date_expires)")
	qb = qb.Concat("values ($%d, $%d)", userUuid, exp)
	qb = qb.Concat("returning id, user_id, date_added, date_expires;")

	err = s.DB.Pool.QueryRow(context.TODO(), qb.Build(), qb.GetArgs()...).Scan(
		&session.Id,
		&session.UserId,
		&session.DateAdded,
		&session.DateExpires,
	)
	if err != nil {
		log.Println(err)
		return session, fmt.Errorf("failed to create session")
	}

	return session, nil
}

func (s PgSessionRepository) FindSession(sessionId string) (model.Session, error) {
	var session model.Session
	var sessionUuid pgtype.UUID
	err := sessionUuid.Scan(sessionId)
	if err != nil {
		return session, fmt.Errorf("invalid user id")
	}

	qb := util.NewSqlBuilder(
		"select id, user_id, date_added, date_expires",
		"from sessions",
	)
	qb = qb.Concat("where id = $%d", sessionUuid)

	err = s.DB.Pool.QueryRow(context.TODO(), qb.Build(), qb.GetArgs()...).Scan(
		&session.Id,
		&session.UserId,
		&session.DateAdded,
		&session.DateExpires,
	)
	if err != nil {
		log.Println(err)
		return session, fmt.Errorf("failed to query session")
	}

	return session, nil
}

package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	GetUser(ctx context.Context, userId string) (model.IUser, error)
	LoginOrRegisterUser(ctx context.Context, openId openid.StandardClaims) (string, error)
	UpdateUser(ctx context.Context, userId string, entity model.TUser) error
}

type pgUserRepository struct {
	log logger.Logger
	db  *database.PostgresDb
}

func NewPgUserRepository(db *database.PostgresDb, log logger.Logger) *pgUserRepository {
	return &pgUserRepository{
		db:  db,
		log: log,
	}
}

func (u *pgUserRepository) GetUser(ctx context.Context, userId string) (model.IUser, error) {
	var user model.IUser

	res, err := u.getUserQuery(userId)
	if err != nil {
		return user, err
	}

	err = u.db.Pool.QueryRow(ctx, res.Query, res.Args...).Scan(
		&user.Id, &user.ProgramId, &user.FullName, &user.Nickname, &user.CurrentYear, &user.Gender, &user.Email,
		&user.Picture, &user.IsActive, &user.DateAdded, &user.DateModified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, util.ErrNotFound
		}

		u.log.Error("get user query failed", logger.Err(err))
		return user, util.ErrInternal
	}

	return user, nil
}

func (u *pgUserRepository) UpdateUser(ctx context.Context, userId string, entity model.TUser) error {
	res, err := u.updateUserQuery(userId, entity)
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, res.Query, res.Args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == database.PgConflictErrCode {
				u.log.Warn("conflict user update error", logger.Err(err))
				return util.ErrConflict
			}
			if pgErr.Code == database.PgCheckErrCode {
				return util.ErrMalformed
			}
		}

		u.log.Error("unknown user update error", logger.Err(err))
		return util.ErrInternal
	}

	return nil
}

func (u *pgUserRepository) LoginOrRegisterUser(ctx context.Context, openId openid.StandardClaims) (string, error) {
	var userId string

	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		u.log.Error("failed to begin transaction", logger.Err(err))
		return "", util.ErrInternal
	}
	defer tx.Rollback(ctx)

	res := u.getUserByEmailQuery(openId.Email)

	err = tx.QueryRow(ctx, res.Query, res.Args...).Scan(
		&userId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			res := u.registerUserQuery(openId)

			err := tx.QueryRow(ctx, res.Query, res.Args...).Scan(&userId)
			if err != nil {
				var pgErr *pgconn.PgError

				if errors.As(err, &pgErr) {
					if pgErr.Code == database.PgConflictErrCode {
						u.log.Error("user create conflict error - user already exists?", logger.Err(err))
						return "", util.ErrConflict
					}
				}

				u.log.Error("insert user query failed", logger.Err(err))
				return "", util.ErrInternal
			}

			u.log.Info("user registered", logger.Err(err))
		} else {
			u.log.Error("user query failed", logger.Err(err))
			return "", util.ErrInternal
		}
	} else {
		// TODO: update user if open id parameters change
	}

	if err := tx.Commit(ctx); err != nil {
		u.log.Error("failed to commit transaction", logger.Err(err))
		return "", util.ErrInternal
	}

	return userId, nil
}

func (u *pgUserRepository) getUserQuery(userId string) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder(
		"select id, program_id, full_name, nickname, current_year, gender, email, picture, is_active, date_added, date_modified",
		"from users",
	)
	qb = qb.Concat("where id = $%d", userUuid)

	return qb.Result(), nil
}

func (u *pgUserRepository) updateUserQuery(userId string, entity model.TUser) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		u.log.Info("invalid user id")
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("update users")
	qb = qb.Concat("set date_modified = $%d", time.Now())

	if entity.ProgramId != "" {
		var programUuid pgtype.UUID
		err := programUuid.Scan(entity.ProgramId)
		if err != nil {
			u.log.Info("invalid program id")
			return util.SqlBuilderResult{}, util.ErrMalformed
		}

		qb = qb.Concat(",program_id = $%d", programUuid)
	}
	if entity.Nickname != "" {
		qb = qb.Concat(",nickname = $%d", entity.Nickname)
	}
	if entity.CurrentYear != 0 {
		qb = qb.Concat(",current_year = $%d", entity.CurrentYear)
	}
	if entity.Gender != "" {
		qb = qb.Concat(",gender = $%d", entity.Gender)
	}
	if entity.Picture != "" {
		qb = qb.Concat(",picture = $%d", entity.Gender)
	}

	qb = qb.Concat("where id = $%d", userUuid)

	return qb.Result(), nil
}

func (u *pgUserRepository) getUserByEmailQuery(email string) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select id",
		"from users",
	)
	qb = qb.Concat("where lower(email) = $%d", strings.ToLower(email))

	return qb.Result()
}

func (u *pgUserRepository) registerUserQuery(openId openid.StandardClaims) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("insert into users (full_name, email, picture, is_active)")
	qb = qb.Concat("values ($%d, $%d, $%d, $%d)", openId.Name, openId.Email, openId.Picture, true)
	qb = qb.Concat("returning id")

	return qb.Result()
}

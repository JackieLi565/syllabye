package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	LoginOrRegisterUser(openId openid.StandardClaims) (string, error)
	CompleteUserSignUp(userId string, payload model.UserSignUpRequest) error
}

type pgUserRepository struct {
	log logger.Logger
	db  *database.DB
}

func NewPgUserRepository(db *database.DB, log logger.Logger) *pgUserRepository {
	return &pgUserRepository{
		db:  db,
		log: log,
	}
}

func (u *pgUserRepository) LoginOrRegisterUser(openId openid.StandardClaims) (string, error) {
	var userId string

	err := u.db.RunTransaction(context.TODO(), func(tx pgx.Tx) error {
		res := u.getUserByEmailQuery(openId.Email)

		err := tx.QueryRow(context.TODO(), res.Query, res.Args...).Scan(
			&userId,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				res := u.registerUserQuery(openId)
				err := tx.QueryRow(context.TODO(), res.Query, res.Args...).Scan(&userId)
				if err != nil {
					if database.IsErrConflict(err) {
						return util.ErrConflict
					}

					u.log.Error("insert user query failed", logger.Err(err))
					return util.ErrInternal
				}

				u.log.Info("insert user query failed", logger.Err(err))
			} else {
				u.log.Error("user query failed", logger.Err(err))
				return util.ErrInternal
			}
		}

		return nil
	})

	return userId, err
}

func (u *pgUserRepository) CompleteUserSignUp(userId string, payload model.UserSignUpRequest) error {
	res, err := u.updateUserSignUpQuery(userId, payload)
	_, err = u.db.Pool.Exec(context.TODO(), res.Query, res.Args...)
	if err != nil {
		if database.IsErrConflict(err) {
			return util.ErrConflict
		}

		u.log.Error("update user query failed", logger.Err(err))
		return util.ErrInternal
	}

	return nil
}

func (u *pgUserRepository) getUserByEmailQuery(email string) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select id",
		"from users u",
	)
	qb = qb.Concat("where lower(email) = $%d", strings.ToLower(email))

	return qb.Result()
}

func (u *pgUserRepository) registerUserQuery(openId openid.StandardClaims) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("insert into users (full_name, email, picture, is_active)")
	qb = qb.Concat("values ($1, $2, $3, $4)", openId.Name, openId.Email, openId.Picture, true)
	qb = qb.Concat("returning id")

	return qb.Result()
}

func (u *pgUserRepository) updateUserSignUpQuery(userId string, payload model.UserSignUpRequest) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	var programUuid pgtype.UUID
	err = programUuid.Scan(payload.ProgramId)

	if payload.Gender != "Male" || payload.Gender != "Female" || payload.Gender != "Other" {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("update users")
	qb = qb.Concat(
		"set program_id = $%d, nickname = $%d, current_year = $%d, gender = $%d",
		programUuid,
		payload.Nickname,
		payload.CurrentYear,
		payload.Gender,
	)
	qb = qb.Concat("where id = $%d", userId)

	return qb.Result(), nil
}

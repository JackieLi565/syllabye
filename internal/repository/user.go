package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	LoginOrRegisterUser(openId openid.StandardClaims) (string, error)
	CompleteUserSignUp(userId string, payload model.UserSignUpRequest) error
}

type PgUserRepository struct {
	DB *database.DB
}

func (u *PgUserRepository) LoginOrRegisterUser(openId openid.StandardClaims) (string, error) {
	var userId string

	err := u.DB.RunTransaction(context.TODO(), func(tx pgx.Tx) error {
		userQb := util.NewSqlBuilder(
			"select id",
			"from users u",
		)
		userQb = userQb.Concat("where lower(email) = $%d", strings.ToLower(openId.Email))

		err := tx.QueryRow(context.TODO(), userQb.Build(), userQb.GetArgs()...).Scan(
			&userId,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				userQb := util.NewSqlBuilder("insert into users (full_name, email, picture, is_active)")
				userQb = userQb.Concat("values ($1, $2, $3, $4)", openId.Name, openId.Email, openId.Picture, true)
				userQb = userQb.Concat("returning id")

				err := tx.QueryRow(context.TODO(), userQb.Build(), userQb.GetArgs()...).Scan(&userId)
				if err != nil {
					log.Println(err)
					return fmt.Errorf("failed to create new user")
				}

				log.Println("User not found, a new user was created")
			} else {
				log.Println(err)
				return fmt.Errorf("failed to query for user")
			}
		}

		return nil
	})

	return userId, err
}

func (u *PgUserRepository) CompleteUserSignUp(userId string, payload model.UserSignUpRequest) error {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return fmt.Errorf("Invalid user ID value.")
	}

	var programUuid pgtype.UUID
	err = programUuid.Scan(payload.ProgramId)

	err = u.checkGender(payload.Gender)
	if err != nil {
		return err
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

	_, err = u.DB.Pool.Exec(context.TODO(), qb.Build(), qb.GetArgs()...)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Unable to complete user sign up.")
	}

	return nil
}

func (u *PgUserRepository) checkGender(gender string) error {
	if gender != "Male" || gender != "Female" || gender != "Other" {
		return fmt.Errorf("Invalid gender value.")
	}

	return nil
}

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
)

type UserRepository interface {
	LoginOrRegisterUser(openId openid.StandardClaims) (model.User, error)
}

type PgUserRepository struct {
	DB *database.DB
}

func (u PgUserRepository) LoginOrRegisterUser(openId openid.StandardClaims) (model.User, error) {
	var user model.User

	err := u.DB.RunTransaction(context.TODO(), func(tx pgx.Tx) error {
		userQb := util.NewSqlBuilder(
			"select id, program_id, full_name, nickname, current_year, gender, email, picture, is_active",
			"from users",
		)
		userQb = userQb.Concat("where lower(email) = $%d", strings.ToLower(openId.Email))

		err := tx.QueryRow(context.TODO(), userQb.Build(), userQb.GetArgs()...).Scan(
			&user.Id,
			&user.ProgramId,
			&user.FullName,
			&user.Nickname,
			&user.CurrentYear,
			&user.Gender,
			&user.Email,
			&user.Picture,
			&user.IsActive,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				userQb := util.NewSqlBuilder("insert into users (full_name, email, picture, is_active)")
				userQb = userQb.Concat("values ($1, $2, $3, $4)", openId.Name, openId.Email, openId.Picture, true)
				userQb = userQb.Concat("returning id")

				var newUserId string
				err := tx.QueryRow(context.TODO(), userQb.Build(), userQb.GetArgs()...).Scan(&newUserId)
				if err != nil {
					log.Println(err)
					return fmt.Errorf("failed to create new user")
				}

				user.Id = newUserId
				user.FullName = openId.Name
				user.Email = openId.Email
				user.Picture.Scan(openId.Picture)

				log.Println("User not found, a new user was created")
			} else {
				log.Println(err)
				return fmt.Errorf("failed to query for user")
			}
		}

		return nil
	})

	return user, err
}

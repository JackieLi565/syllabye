package repository

import (
	"context"
	"errors"
	"fmt"
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

	AddUserCourse(ctx context.Context, userId string, entity model.TUserCourse) error
	DeleteUserCourse(ctx context.Context, userId string, courseId string) error
	UpdateUserCourse(ctx context.Context, userId string, courseId string, entity model.TUserCourse) error
	ListUserCourses(ctx context.Context, userId string, filters model.CourseFilters, paginate util.Paginate) ([]model.ICourse, error)
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

			u.log.Info(fmt.Sprintf("user registered with id %s", userId))
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

func (u *pgUserRepository) AddUserCourse(ctx context.Context, userId string, entity model.TUserCourse) error {
	result, err := u.addUserCourseQuery(userId, entity)
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, result.Query, result.Args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == database.PgCheckErrCode || pgErr.Code == database.PgInvalidTextRepErrCode {
				return util.ErrMalformed
			} else if pgErr.Code == database.PgConflictErrCode || pgErr.Code == database.PgFKeyViolationErrCode {
				return util.ErrConflict
			}
		}

		u.log.Error("un-handled add user course query error", logger.Err(err))
		return util.ErrInternal
	}

	u.log.Info(fmt.Sprintf("user %s added course %s", userId, entity.CourseId))
	return nil
}

func (u *pgUserRepository) addUserCourseQuery(userId string, entity model.TUserCourse) (util.SqlBuilderResult, error) {
	courseUuid, err := database.ParsePgUuid(entity.CourseId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("insert into user_courses (user_id, course_id, year_taken, semester_taken)")
	qb.Concat("values ($%d, $%d, $%d, $%d)", userId, courseUuid, entity.YearTaken, entity.SemesterTaken)

	return qb.Result(), nil
}

func (u *pgUserRepository) DeleteUserCourse(ctx context.Context, userId string, courseId string) error {
	result, err := u.deleteUserCourseQuery(userId, courseId)
	if err != nil {
		return err
	}

	err = u.db.Pool.QueryRow(ctx, result.Query, result.Args...).Scan(new(interface{}))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.ErrNotFound
		}

		u.log.Error("un-handled delete user course query error", logger.Err(err))
		return util.ErrInternal
	}

	u.log.Info(fmt.Sprintf("user %s removed course %s", userId, courseId))
	return nil
}

func (u *pgUserRepository) deleteUserCourseQuery(userId string, courseId string) (util.SqlBuilderResult, error) {
	courseUuid, err := database.ParsePgUuid(courseId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("delete from user_courses")
	qb.Concat("where user_id = $%d and course_id = $%d", userId, courseUuid)
	qb.Concat("returning course_id")

	return qb.Result(), nil
}

func (u *pgUserRepository) UpdateUserCourse(ctx context.Context, userId string, courseId string, entity model.TUserCourse) error {
	result, err := u.updateUserCourseQuery(userId, courseId, entity)
	if err != nil {
		return err
	}

	err = u.db.Pool.QueryRow(ctx, result.Query, result.Args...).Scan(new(interface{}))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == database.PgCheckErrCode || pgErr.Code == database.PgInvalidTextRepErrCode {
				return util.ErrMalformed
			} else if pgErr.Code == database.PgFKeyViolationErrCode {
				return util.ErrConflict
			}
		}

		u.log.Error("un-handled update user course query error", logger.Err(err))
		return util.ErrInternal
	}

	return nil
}

func (u *pgUserRepository) updateUserCourseQuery(userId string, courseId string, entity model.TUserCourse) (util.SqlBuilderResult, error) {
	courseUuid, err := database.ParsePgUuid(courseId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("update user_courses")
	qb.Concat("set date_modified = $%d", time.Now())

	if entity.SemesterTaken != nil {
		qb.Concat(",semester_taken = $%d", entity.SemesterTaken)
	}

	if entity.YearTaken != nil {
		qb.Concat(",year_taken = $%d", entity.YearTaken)
	}

	qb.Concat("where user_id = $%d and course_id = $%d", userId, courseUuid)
	qb.Concat("returning course_id")

	return qb.Result(), nil
}

func (u *pgUserRepository) ListUserCourses(ctx context.Context, userId string, filters model.CourseFilters, paginate util.Paginate) ([]model.ICourse, error) {
	result := u.listUserCoursesQuery(userId, filters, paginate)

	rows, err := u.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.ICourse{}, nil
		}

		u.log.Error("un-handled list user course query error", logger.Err(err))
		return []model.ICourse{}, util.ErrInternal
	}

	var courses []model.ICourse
	for rows.Next() {
		course := model.ICourse{}
		err := rows.Scan(
			&course.Id, &course.CategoryId, &course.Title, &course.Description, &course.Uri,
			&course.Course, &course.DateAdded,
		)
		if err != nil {
			u.log.Error("scan internal course error", logger.Err(err))
			return courses, util.ErrInternal
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (u *pgUserRepository) listUserCoursesQuery(userId string, filters model.CourseFilters, paginate util.Paginate) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select id, category_id, title, description, uri, course, date_added",
		"from courses",
	)
	qb.Concat("where exists (")
	qb.Concat("select 1 from user_courses uc where uc.course_id = id and uc.user_id = $%d", userId)
	qb.Concat(")")

	if filters.CategoryId != "" {
		qb.Concat("and category_id = $%d", filters.CategoryId)
	}
	if filters.Search != "" {
		qb.Concat("and course ilike $%d or title ilike $%d", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	qb.Concat("limit $%d", paginate.Size)
	offset := (paginate.Page - 1) * paginate.Size
	qb.Concat("offset $%d", offset)

	return qb.Result()
}

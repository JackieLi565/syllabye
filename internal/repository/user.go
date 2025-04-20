package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/nullable"
)

type UserSchema struct {
	Id           string
	ProgramId    sql.NullString
	FullName     string
	Nickname     sql.NullString
	CurrentYear  sql.NullInt16
	Gender       sql.NullString
	Email        string
	Bio          sql.NullString
	IgHandle     sql.NullString
	Picture      sql.NullString
	IsActive     bool
	DateAdded    time.Time
	DateModified time.Time
}

type UpdateUser struct {
	ProgramId   nullable.Nullable[string]
	Nickname    nullable.Nullable[string]
	CurrentYear nullable.Nullable[int16]
	Gender      nullable.Nullable[string]
	Bio         nullable.Nullable[string]
	IgHandle    nullable.Nullable[string]
}

type UserCourseSchema struct {
	UserId        string
	CourseId      string
	Title         string
	Course        string
	YearTaken     sql.NullInt16
	SemesterTaken sql.NullString
	DateAdded     time.Time
	DateModified  time.Time
}

type InsertUserCourse struct {
	CourseId      string
	YearTaken     *int16
	SemesterTaken *string
}

type UpdateUserCourse struct {
	YearTaken     nullable.Nullable[int16]
	SemesterTaken nullable.Nullable[string]
}

type UserRepository interface {
	GetUser(ctx context.Context, userId string) (UserSchema, error)
	LoginOrRegisterUser(ctx context.Context, openId openid.StandardClaims) (string, error)
	UpdateUser(ctx context.Context, userId string, entity UpdateUser) error
	SearchUserNickname(ctx context.Context, nickname string) (bool, error)

	AddUserCourse(ctx context.Context, userId string, entity InsertUserCourse) error
	DeleteUserCourse(ctx context.Context, userId string, courseId string) error
	UpdateUserCourse(ctx context.Context, userId string, courseId string, entity UpdateUserCourse) error
	ListUserCourses(ctx context.Context, userId string, filters CourseFilters, paginate util.Paginate) ([]UserCourseSchema, error)
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

func (u *pgUserRepository) GetUser(ctx context.Context, userId string) (UserSchema, error) {

	res, err := u.getUserQuery(userId)
	if err != nil {
		return UserSchema{}, err
	}

	var userSchema UserSchema
	err = u.db.Pool.QueryRow(ctx, res.Query, res.Args...).Scan(
		&userSchema.Id, &userSchema.ProgramId, &userSchema.FullName, &userSchema.Nickname, &userSchema.CurrentYear, &userSchema.Gender, &userSchema.Email,
		&userSchema.Picture, &userSchema.IsActive, &userSchema.DateAdded, &userSchema.DateModified, &userSchema.Bio, &userSchema.IgHandle,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserSchema{}, util.ErrNotFound
		}

		u.log.Error("get user query failed", logger.Err(err))
		return UserSchema{}, util.ErrInternal
	}

	return userSchema, nil
}

func (u *pgUserRepository) getUserQuery(userId string) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder(
		"select id, program_id, full_name, nickname, current_year, gender, email, picture, is_active, date_added, date_modified, bio, ig_handle",
		"from users",
	)
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

func (u *pgUserRepository) UpdateUser(ctx context.Context, userId string, entity UpdateUser) error {
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

func (u *pgUserRepository) updateUserQuery(userId string, entity UpdateUser) (util.SqlBuilderResult, error) {
	var userUuid pgtype.UUID
	err := userUuid.Scan(userId)
	if err != nil {
		u.log.Info("invalid user id")
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("update users set id = id")

	if entity.ProgramId.IsSpecified() {
		stm := ",program_id = $%d"
		programId, err := entity.ProgramId.Get()
		if err != nil {
			qb.Concat(stm, nil)
		} else {
			var programUuid pgtype.UUID
			err := programUuid.Scan(programId)
			if err != nil {
				u.log.Info("invalid program id")
				return util.SqlBuilderResult{}, util.ErrMalformed
			}

			qb = qb.Concat(stm, programUuid)
		}
	}
	if entity.Nickname.IsSpecified() {
		stm := ",nickname = $%d"
		nickname, err := entity.Nickname.Get()
		if err != nil {
			qb.Concat(stm, nil)
		} else {
			qb.Concat(stm, nickname)
		}
	}
	if entity.CurrentYear.IsSpecified() {
		stm := ",current_year = $%d"
		currentYear, err := entity.CurrentYear.Get()
		if err != nil {
			qb.Concat(stm, nil)
		} else {
			qb.Concat(stm, currentYear)
		}
	}
	if entity.Gender.IsSpecified() {
		stm := ",gender = $%d"
		gender, err := entity.Gender.Get()
		if err != nil {
			qb.Concat(stm, nil)
		} else {
			qb.Concat(stm, gender)
		}
	}
	if entity.Bio.IsSpecified() {
		bio, err := entity.Bio.Get()
		if err != nil {
			qb.Concat(",bio = null")
		} else {
			qb.Concat(",bio = $%d", bio)
		}
	}
	if entity.IgHandle.IsSpecified() {
		ig, err := entity.IgHandle.Get()
		if err != nil {
			qb.Concat(",ig_handle = null")
		} else {
			qb.Concat(",ig_handle = $%d", ig)
		}
	}

	qb = qb.Concat("where id = $%d", userUuid)

	return qb.Result(), nil
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

func (u *pgUserRepository) registerUserQuery(openId openid.StandardClaims) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("insert into users (full_name, email, picture, is_active)")
	qb = qb.Concat("values ($%d, $%d, $%d, $%d)", openId.Name, openId.Email, openId.Picture, true)
	qb = qb.Concat("returning id")

	return qb.Result()
}

func (u *pgUserRepository) AddUserCourse(ctx context.Context, userId string, entity InsertUserCourse) error {
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

func (u *pgUserRepository) addUserCourseQuery(userId string, entity InsertUserCourse) (util.SqlBuilderResult, error) {
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

func (u *pgUserRepository) UpdateUserCourse(ctx context.Context, userId string, courseId string, entity UpdateUserCourse) error {
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

func (u *pgUserRepository) updateUserCourseQuery(userId string, courseId string, entity UpdateUserCourse) (util.SqlBuilderResult, error) {
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

func (u *pgUserRepository) ListUserCourses(ctx context.Context, userId string, filters CourseFilters, paginate util.Paginate) ([]UserCourseSchema, error) {
	result := u.listUserCoursesQuery(userId, filters, paginate)

	rows, err := u.db.Pool.Query(context.TODO(), result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []UserCourseSchema{}, nil
		}

		u.log.Error("un-handled list user course query error", logger.Err(err))
		return []UserCourseSchema{}, util.ErrInternal
	}

	var courses []UserCourseSchema
	for rows.Next() {
		course := UserCourseSchema{}
		err := rows.Scan(
			&course.UserId,
			&course.CourseId,
			&course.Title,
			&course.Course,
			&course.YearTaken,
			&course.SemesterTaken,
		)
		if err != nil {
			u.log.Error("scan internal user course error", logger.Err(err))
			return courses, util.ErrInternal
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (u *pgUserRepository) listUserCoursesQuery(userId string, filters CourseFilters, paginate util.Paginate) util.SqlBuilderResult {
	qb := util.NewSqlBuilder(
		"select uc.user_id, uc.course_id, c.title, c.course, uc.year_taken, uc.semester_taken",
		"from user_courses uc",
		"inner join courses c on c.id = uc.course_id",
	)
	qb.Concat("where uc.user_id = $%d", userId)

	if filters.CategoryId != "" {
		qb.Concat("and c.category_id = $%d", filters.CategoryId)
	}
	if filters.Search != "" {
		qb.Concat("and c.course ilike $%d or c.title ilike $%d", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	qb.Concat("limit $%d", paginate.Size)
	offset := (paginate.Page - 1) * paginate.Size
	qb.Concat("offset $%d", offset)

	return qb.Result()
}

func (u *pgUserRepository) SearchUserNickname(ctx context.Context, nickname string) (bool, error) {
	result := u.searchNicknameQuery(nickname)

	var existsFlag int8
	err := u.db.Pool.QueryRow(ctx, result.Query, result.Args...).Scan(&existsFlag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		u.log.Error("un-handled search nickname query error", logger.Err(err))
		return false, util.ErrInternal
	}

	return existsFlag == 1, nil
}

func (u *pgUserRepository) searchNicknameQuery(nickname string) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("select 1 from users")
	qb.Concat("where nickname = $%d", strings.ToLower(nickname))

	return qb.Result()
}

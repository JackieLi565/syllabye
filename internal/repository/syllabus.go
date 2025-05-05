package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/nullable"
)

type SyllabusSchema struct {
	Id          string
	UserId      string
	CourseId    string
	File        string
	FileSize    int
	ContentType string
	Year        int16
	Semester    string
	DateAdded   time.Time
	DateSynced  sql.NullTime
}

type InsertSyllabus struct {
	UserId      string
	CourseId    string
	File        string
	FileSize    int
	ContentType string
	Checksum    string
	Year        int16
	Semester    string
}

type UpdateSyllabus struct {
	Year     nullable.Nullable[int16]
	Semester nullable.Nullable[string]
}

type SyllabusFilters struct {
	UserId   string
	CourseId string
	Year     *int16
	Semester string
}

type SyllabusLikeSchema struct {
	SyllabusId string
	UserId     string
	IsDislike  bool
	DateAdded  time.Time
}

type SyllabusMeta struct {
	Id        string
	Course    string
	UserId    string
	UserName  string
	UserEmail string
}

type SyllabusRepository interface {
	GetAndViewSyllabus(ctx context.Context, userId string, syllabusId string) (SyllabusSchema, error)
	CreateSyllabus(ctx context.Context, syllabus InsertSyllabus) (string, error)
	ListSyllabi(ctx context.Context, userId string, filters SyllabusFilters, paginate util.Paginate) ([]SyllabusSchema, error)
	DeleteSyllabus(ctx context.Context, userId string, syllabusId string) error
	UpdateSyllabus(ctx context.Context, userId string, syllabusId string, syllabus UpdateSyllabus) error
	// SyncSyllabus updates a syllabus with a valid date_synced value.
	SyncSyllabus(ctx context.Context, syllabusId string) error
	VerifySyllabus(ctx context.Context, syllabusId string) (bool, SyllabusMeta, error)
	ListSyllabusLikes(ctx context.Context, syllabusId string) ([]SyllabusLikeSchema, error)
	LikeSyllabus(ctx context.Context, userId string, syllabusId string, dislike bool) error
	DeleteSyllabusLike(ctx context.Context, userId string, syllabusId string) error
}

type pgSyllabusRepository struct {
	db  *database.PostgresDb
	log logger.Logger
}

func NewPgSyllabusRepository(db *database.PostgresDb, log logger.Logger) *pgSyllabusRepository {
	return &pgSyllabusRepository{
		db:  db,
		log: log,
	}
}

func (s *pgSyllabusRepository) GetAndViewSyllabus(ctx context.Context, userId string, syllabusId string) (SyllabusSchema, error) {
	getResult, err := s.getActiveSyllabusQuery(userId, syllabusId)
	viewResult, _ := s.incrementSyllabusView(userId, syllabusId) // No need to handel err (pre handle id in getSyllabusQuery)
	if err != nil {
		return SyllabusSchema{}, err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Err(err))
		return SyllabusSchema{}, util.ErrInternal
	}
	defer tx.Rollback(ctx)

	syllabus := SyllabusSchema{}
	err = tx.QueryRow(ctx, getResult.Query, getResult.Args...).Scan(
		&syllabus.Id,
		&syllabus.UserId,
		&syllabus.CourseId,
		&syllabus.File,
		&syllabus.FileSize,
		&syllabus.ContentType,
		&syllabus.Year,
		&syllabus.Semester,
		&syllabus.DateAdded,
		&syllabus.DateSynced,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Info("syllabus not found")
			return SyllabusSchema{}, util.ErrNotFound
		}

		s.log.Error("un-handled get syllabus error", logger.Err(err))
		return SyllabusSchema{}, util.ErrInternal
	}

	if _, err := tx.Exec(ctx, viewResult.Query, viewResult.Args...); err != nil {
		s.log.Error("un-handled view syllabus error", logger.Err(err))
		return SyllabusSchema{}, util.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", logger.Err(err))
		return SyllabusSchema{}, util.ErrInternal
	}

	return syllabus, nil
}

func (s *pgSyllabusRepository) incrementSyllabusView(userId string, syllabusId string) (util.SqlBuilderResult, error) {
	var syllabusUuid pgtype.UUID
	if err := syllabusUuid.Scan(syllabusId); err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("insert into syllabus_views (syllabus_id, user_id)")
	qb.Concat("values ($%d, $%d) on conflict (syllabus_id, user_id) do nothing;", syllabusUuid, userId)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) getActiveSyllabusQuery(userId string, syllabusId string) (util.SqlBuilderResult, error) {
	var syllabusUuid pgtype.UUID
	if err := syllabusUuid.Scan(syllabusId); err != nil {
		s.log.Info("invalid syllabus id")
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("select id, user_id, course_id, file, file_size, content_type, year, semester, date_added, date_synced from syllabi")
	qb.Concat("where id = $%d", syllabusId)
	qb.Concat("and (date_synced is not null or user_id = $%d)", userId)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) CreateSyllabus(ctx context.Context, syllabus InsertSyllabus) (string, error) {
	result := s.createSyllabusQuery(syllabus)

	var syllabusId string
	err := s.db.Pool.QueryRow(ctx, result.Query, result.Args...).Scan(&syllabusId)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == database.PgCheckErrCode {
			s.log.Info("syllabus failed database check")
			return "", util.ErrMalformed
		}

		s.log.Error("un-handled create syllabus query error", logger.Err(err))
		return "", util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("syllabus %s created", syllabusId))
	return syllabusId, nil
}

func (s *pgSyllabusRepository) createSyllabusQuery(sy InsertSyllabus) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("insert into syllabi (user_id, course_id, file, file_size, content_type, year, semester)")
	qb.Concat("values ($%d, $%d, $%d, $%d, $%d, $%d, $%d)", sy.UserId, sy.CourseId, sy.File, sy.FileSize, sy.ContentType, sy.Year, sy.Semester)
	qb.Concat("returning id")

	return qb.Result()
}

func (s *pgSyllabusRepository) ListSyllabi(ctx context.Context, userId string, filters SyllabusFilters, paginate util.Paginate) ([]SyllabusSchema, error) {
	result, err := s.listSyllabiQuery(userId, filters, paginate)
	if err != nil {
		return []SyllabusSchema{}, err
	}

	rows, err := s.db.Pool.Query(ctx, result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []SyllabusSchema{}, nil
		}
		s.log.Error("un-handled list syllabi query error", logger.Err(err))
		return []SyllabusSchema{}, util.ErrInternal
	}

	var syllabi []SyllabusSchema
	for rows.Next() {
		syllabus := SyllabusSchema{}
		err := rows.Scan(
			&syllabus.Id,
			&syllabus.UserId,
			&syllabus.CourseId,
			&syllabus.File,
			&syllabus.FileSize,
			&syllabus.ContentType,
			&syllabus.Year,
			&syllabus.Semester,
			&syllabus.DateAdded,
			&syllabus.DateSynced,
		)
		if err != nil {
			s.log.Error("scan syllabus error", logger.Err(err))
			return []SyllabusSchema{}, util.ErrInternal
		}

		syllabi = append(syllabi, syllabus)
	}

	return syllabi, nil
}

func (s *pgSyllabusRepository) listSyllabiQuery(userId string, filters SyllabusFilters, paginate util.Paginate) (util.SqlBuilderResult, error) {
	qb := util.NewSqlBuilder("select id, user_id, course_id, file, file_size, content_type, year, semester, date_added, date_synced from syllabi")
	qb.Concat("where (date_synced is not null or user_id = $%d)", userId)

	if filters.UserId != "" {
		var userUuid pgtype.UUID
		if err := userUuid.Scan(filters.UserId); err != nil {
			s.log.Info("invalid user id")
			return util.SqlBuilderResult{}, util.ErrMalformed
		}
		qb.Concat("and user_id = $%d", userUuid)
	}

	if filters.CourseId != "" {
		var courseUuid pgtype.UUID
		if err := courseUuid.Scan(filters.CourseId); err != nil {
			s.log.Info("invalid course id")
			return util.SqlBuilderResult{}, util.ErrMalformed
		}
		qb.Concat("and course_id = $%d", courseUuid)
	}

	if filters.Year != nil {
		if *filters.Year <= 0 {
			s.log.Info("year zero or less was passed as a syllabus query filter")
		}
		qb.Concat("and year = $%d", *filters.Year)
	}

	if filters.Semester != "" {
		qb.Concat("and semester = $%d", filters.Semester)
	}

	qb.Concat("limit $%d", paginate.Size)
	offset := (paginate.Page - 1) * paginate.Size
	qb.Concat("offset $%d", offset)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) DeleteSyllabus(ctx context.Context, userId string, syllabusId string) error {
	result, err := s.deleteSyllabusQuery(syllabusId)
	if err != nil {
		return err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Err(err))
		return util.ErrInternal
	}
	defer tx.Rollback(ctx)

	var createUserId string
	err = tx.QueryRow(ctx, result.Query, result.Args...).Scan(&createUserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.ErrNotFound
		}

		s.log.Error("un-handled delete syllabus query error", logger.Err(err))
		return util.ErrInternal
	}

	if createUserId != userId {
		s.log.Info(fmt.Sprintf("user %s attempted to delete user %s syllabus %s", userId, createUserId, syllabusId))
		return util.ErrForbidden
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("syllabus %s deleted", syllabusId))
	return nil
}

func (s *pgSyllabusRepository) deleteSyllabusQuery(syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := database.ParsePgUuid(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("delete from syllabi")
	qb.Concat("where id = $%d", syllabusUuid)
	qb.Concat("returning user_id")

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) UpdateSyllabus(ctx context.Context, userId string, syllabusId string, syllabus UpdateSyllabus) error {
	result, err := s.updateSyllabusQuery(syllabusId, syllabus)
	if err != nil {
		return err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Err(err))
		return util.ErrInternal
	}
	defer tx.Rollback(ctx)

	var createUserId string
	err = tx.QueryRow(ctx, result.Query, result.Args...).Scan(&createUserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.ErrNotFound
		}

		s.log.Error("un-handled update syllabus query error", logger.Err(err))
		return util.ErrInternal
	}

	if createUserId != userId {
		s.log.Info(fmt.Sprintf("user %s attempted to update user %s syllabus %s", userId, createUserId, syllabusId))
		return util.ErrForbidden
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("syllabus %s updated", syllabusId))
	return nil
}

func (s *pgSyllabusRepository) updateSyllabusQuery(syllabusId string, syllabus UpdateSyllabus) (util.SqlBuilderResult, error) {
	var syllabusUuid pgtype.UUID
	if err := syllabusUuid.Scan(syllabusId); err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("update syllabi")
	qb.Concat("set date_modified = $%d", time.Now())

	if syllabus.Year.IsSpecified() {
		year, err := syllabus.Year.Get()
		if err != nil {
			qb.Concat(",year = $%d", year)
		}
	}

	if syllabus.Semester.IsSpecified() {
		semester, err := syllabus.Semester.Get()
		if err != nil {
			qb.Concat(",semester = $%d", semester)
		}
	}

	qb.Concat("where id = $%d", syllabusUuid)
	qb.Concat("returning user_id")

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) SyncSyllabus(ctx context.Context, syllabusId string) error {
	result, err := s.syncSyllabusQuery(syllabusId)
	if err != nil {
		return err
	}

	_, err = s.db.Pool.Exec(ctx, result.Query, result.Args...)
	if err != nil {
		s.log.Error("un-handled sync syllabus query error", logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("syllabus %s synced", syllabusId))
	return nil
}

func (s *pgSyllabusRepository) syncSyllabusQuery(syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := s.validateSyllabusId(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("update syllabi")
	qb.Concat("set date_synced = $%d", time.Now())
	qb.Concat("where id = $%d", syllabusUuid)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) LikeSyllabus(ctx context.Context, userId string, syllabusId string, dislike bool) error {
	deleteResult, err := s.deleteSyllabusLikeQuery(userId, syllabusId)
	createResult, _ := s.createSyllabusLikeQuery(userId, syllabusId, dislike)

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Err(err))
		return util.ErrInternal
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteResult.Query, deleteResult.Args...)
	if err != nil {
		s.log.Error("un-handled syllabus like query error", logger.Err(err))
		return util.ErrInternal
	}

	_, err = tx.Exec(ctx, createResult.Query, createResult.Args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == database.PgConflictErrCode {
			s.log.Warn("duplicate syllabus like entry")
			return util.ErrConflict
		}

		s.log.Error("un-handled syllabus like query error", logger.Err(err))
		return util.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("like syllabus %s reaction", syllabusId))
	return nil
}

func (s *pgSyllabusRepository) createSyllabusLikeQuery(userId string, syllabusId string, dislike bool) (util.SqlBuilderResult, error) {
	syllabusUuid, err := s.validateSyllabusId(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("insert into syllabus_likes (syllabus_id, user_id, is_dislike)")
	qb.Concat("values ($%d, $%d, $%d)", syllabusUuid, userId, dislike)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) DeleteSyllabusLike(ctx context.Context, userId string, syllabusId string) error {
	result, err := s.deleteSyllabusLikeQuery(userId, syllabusId)
	if err != nil {
		return err
	}

	_, err = s.db.Pool.Exec(ctx, result.Query, result.Args...)
	if err != nil {
		s.log.Error("un-handled syllabus delete like query error", logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("user %s removed like from syllabus %s", userId, syllabusId))
	return nil
}

func (s *pgSyllabusRepository) deleteSyllabusLikeQuery(userId string, syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := s.validateSyllabusId(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("delete from syllabus_likes")
	qb.Concat("where syllabus_id = $%d and user_id = $%d", syllabusUuid, userId)

	return qb.Result(), nil
}

func (s *pgSyllabusRepository) ListSyllabusLikes(ctx context.Context, syllabusId string) ([]SyllabusLikeSchema, error) {
	result, err := s.listSyllabusLikesQuery(syllabusId)
	if err != nil {
		return []SyllabusLikeSchema{}, err
	}

	rows, err := s.db.Pool.Query(ctx, result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []SyllabusLikeSchema{}, nil
		}

		s.log.Error("un-handled list syllabus likes query error", logger.Err(err))
		return []SyllabusLikeSchema{}, util.ErrInternal
	}

	likes := []SyllabusLikeSchema{}
	for rows.Next() {
		like := SyllabusLikeSchema{}
		err := rows.Scan(
			&like.SyllabusId,
			&like.UserId,
			&like.IsDislike,
			&like.DateAdded,
		)
		if err != nil {
			s.log.Error(fmt.Sprintf("an error occurred when scanning for syllabus likes on syllabus %s", syllabusId), logger.Err(err))
			return []SyllabusLikeSchema{}, util.ErrInternal
		}

		likes = append(likes, like)
	}

	return likes, nil
}

func (s *pgSyllabusRepository) listSyllabusLikesQuery(syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := s.validateSyllabusId(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("select syllabus_id, user_id, is_dislike, date_added from syllabus_likes")
	qb.Concat("where syllabus_id = $%d", syllabusUuid)

	return qb.Result(), nil
}

// Deprecated - use database database.ParsePgUuid()
func (s *pgSyllabusRepository) validateSyllabusId(syllabusId string) (pgtype.UUID, error) {
	var syllabusUuid pgtype.UUID
	if err := syllabusUuid.Scan(syllabusId); err != nil {
		return pgtype.UUID{}, util.ErrMalformed
	}

	return syllabusUuid, nil
}

func (s *pgSyllabusRepository) VerifySyllabus(ctx context.Context, syllabusId string) (bool, SyllabusMeta, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin verify syllabus transaction", logger.Err(err))
		return false, SyllabusMeta{}, util.ErrInternal
	}
	defer tx.Rollback(ctx)

	getResult, err := s.getSyllabusMetaQuery(syllabusId)
	if err != nil {
		return false, SyllabusMeta{}, err
	}

	var dateSynced sql.NullTime
	meta := SyllabusMeta{}
	err = tx.QueryRow(ctx, getResult.Query, getResult.Args...).Scan(
		&meta.Id,
		&dateSynced,
		&meta.UserId,
		&meta.UserName,
		&meta.UserEmail,
		&meta.Course,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Info(fmt.Sprintf("syllabus %s not longer exists", syllabusId))
			return false, SyllabusMeta{}, util.ErrNotFound
		}

		s.log.Error("un-handled verify syllabus query error", logger.Err(err))
		return false, SyllabusMeta{}, util.ErrInternal
	}

	isVerified := true
	// Delete syllabus as it was not uploaded within a given time span
	if !dateSynced.Valid {
		deleteResult, _ := s.deleteSyllabusQuery(syllabusId)
		_, err = tx.Exec(ctx, deleteResult.Query, deleteResult.Args...)
		if err != nil {
			s.log.Error("un-handled delete syllabus query error", logger.Err(err))
			return false, SyllabusMeta{}, util.ErrInternal
		}

		s.log.Info(fmt.Sprintf("syllabus %s removed", syllabusId))
		isVerified = false
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit verify syllabus transaction", logger.Err(err))
		return false, SyllabusMeta{}, util.ErrInternal
	}

	return isVerified, meta, nil
}

func (s *pgSyllabusRepository) getSyllabusMetaQuery(syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := database.ParsePgUuid(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder(
		"select s.id, s.date_synced, u.id as user_id, u.full_name, u.email, c.course",
		"from syllabi s",
		"inner join users u on u.id = s.user_id",
		"inner join courses c on c.id = s.course_id",
	)
	qb.Concat("where s.id = $%d", syllabusUuid)

	return qb.Result(), nil
}

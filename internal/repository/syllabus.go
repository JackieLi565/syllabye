package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type SyllabusRepository interface {
	GetAndViewSyllabus(ctx context.Context, userId string, syllabusId string) (model.ISyllabus, error)
	CreateSyllabus(ctx context.Context, syllabus model.TSyllabus) (string, error)
	ListSyllabi(ctx context.Context, userId string, filters model.SyllabusFilters, paginate util.Paginate) ([]model.ISyllabus, error)
	DeleteSyllabus(ctx context.Context, userId string, syllabusId string) error
	UpdateSyllabus(ctx context.Context, userId string, syllabusId string, syllabus model.TSyllabus) error
	// SyncSyllabus updates a syllabus with a valid date_synced value.
	SyncSyllabus(ctx context.Context, syllabusId string) error
	// VerifySyllabus verifies if a syllabus has a valid sync date, otherwise the syllabus will be removed with a return value of false.
	VerifySyllabus(ctx context.Context, syllabusId string) (bool, error)
	ListSyllabusLikes(ctx context.Context, syllabusId string) ([]model.ISyllabusLike, error)
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

func (s *pgSyllabusRepository) GetAndViewSyllabus(ctx context.Context, userId string, syllabusId string) (model.ISyllabus, error) {
	getResult, err := s.getActiveSyllabusQuery(userId, syllabusId)
	viewResult, _ := s.incrementSyllabusView(userId, syllabusId) // No need to handel err (pre handle id in getSyllabusQuery)
	if err != nil {
		return model.ISyllabus{}, err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Err(err))
		return model.ISyllabus{}, util.ErrInternal
	}
	defer tx.Rollback(ctx)

	syllabus := model.ISyllabus{}
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
			return model.ISyllabus{}, util.ErrNotFound
		}

		s.log.Error("un-handled get syllabus error", logger.Err(err))
		return model.ISyllabus{}, util.ErrInternal
	}

	if _, err := tx.Exec(ctx, viewResult.Query, viewResult.Args...); err != nil {
		s.log.Error("un-handled view syllabus error", logger.Err(err))
		return model.ISyllabus{}, util.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", logger.Err(err))
		return model.ISyllabus{}, util.ErrInternal
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

func (s *pgSyllabusRepository) CreateSyllabus(ctx context.Context, syllabus model.TSyllabus) (string, error) {
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

func (s *pgSyllabusRepository) createSyllabusQuery(sy model.TSyllabus) util.SqlBuilderResult {
	qb := util.NewSqlBuilder("insert into syllabi (user_id, course_id, file, file_size, content_type, year, semester)")
	qb.Concat("values ($%d, $%d, $%d, $%d, $%d, $%d, $%d)", sy.UserId, sy.CourseId, sy.File, sy.FileSize, sy.ContentType, sy.Year, sy.Semester)
	qb.Concat("returning id")

	return qb.Result()
}

func (s *pgSyllabusRepository) ListSyllabi(ctx context.Context, userId string, filters model.SyllabusFilters, paginate util.Paginate) ([]model.ISyllabus, error) {
	result, err := s.listSyllabiQuery(userId, filters, paginate)
	if err != nil {
		return []model.ISyllabus{}, err
	}

	rows, err := s.db.Pool.Query(ctx, result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.ISyllabus{}, nil
		}
		s.log.Error("un-handled list syllabi query error", logger.Err(err))
		return []model.ISyllabus{}, util.ErrInternal
	}

	var syllabi []model.ISyllabus
	for rows.Next() {
		syllabus := model.ISyllabus{}
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
			return []model.ISyllabus{}, util.ErrInternal
		}

		syllabi = append(syllabi, syllabus)
	}

	return syllabi, nil
}

func (s *pgSyllabusRepository) listSyllabiQuery(userId string, filters model.SyllabusFilters, paginate util.Paginate) (util.SqlBuilderResult, error) {
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

func (s *pgSyllabusRepository) UpdateSyllabus(ctx context.Context, userId string, syllabusId string, syllabus model.TSyllabus) error {
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

func (s *pgSyllabusRepository) updateSyllabusQuery(syllabusId string, syllabus model.TSyllabus) (util.SqlBuilderResult, error) {
	var syllabusUuid pgtype.UUID
	if err := syllabusUuid.Scan(syllabusId); err != nil {
		return util.SqlBuilderResult{}, util.ErrMalformed
	}

	qb := util.NewSqlBuilder("update syllabi")
	qb.Concat("set date_modified = $%d", time.Now())

	if syllabus.Year > 0 {
		qb.Concat(",year = $%d", syllabus.Year)
	}

	if syllabus.Semester != "" {
		qb.Concat(",semester = $%d", syllabus.Semester)
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

func (s *pgSyllabusRepository) ListSyllabusLikes(ctx context.Context, syllabusId string) ([]model.ISyllabusLike, error) {
	result, err := s.listSyllabusLikesQuery(syllabusId)
	if err != nil {
		return []model.ISyllabusLike{}, err
	}

	rows, err := s.db.Pool.Query(ctx, result.Query, result.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.ISyllabusLike{}, nil
		}

		s.log.Error("un-handled list syllabus likes query error", logger.Err(err))
		return []model.ISyllabusLike{}, util.ErrInternal
	}

	likes := []model.ISyllabusLike{}
	for rows.Next() {
		like := model.ISyllabusLike{}
		err := rows.Scan(
			&like.SyllabusId,
			&like.UserId,
			&like.IsDislike,
			&like.DateAdded,
		)
		if err != nil {
			s.log.Error(fmt.Sprintf("an error occurred when scanning for syllabus likes on syllabus %s", syllabusId), logger.Err(err))
			return []model.ISyllabusLike{}, util.ErrInternal
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

func (s *pgSyllabusRepository) VerifySyllabus(ctx context.Context, syllabusId string) (bool, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin verify syllabus transaction", logger.Err(err))
		return false, util.ErrInternal
	}
	defer tx.Rollback(ctx)

	getResult, err := s.getDateSyncedSyllabusQuery(syllabusId)
	if err != nil {
		return false, err
	}

	var dateSynced sql.NullTime
	err = tx.QueryRow(ctx, getResult.Query, getResult.Args...).Scan(&dateSynced)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Error("attempt to verify a syllabus that no longer exists")
			return false, util.ErrNotFound
		}

		s.log.Error("un-handled verify syllabus query error", logger.Err(err))
		return false, util.ErrInternal
	}

	// Delete syllabus as it was not uploaded within a given time span
	isVerified := true // Assume the syllabus is verified unless removed
	if !dateSynced.Valid {
		deleteResult, _ := s.deleteSyllabusQuery(syllabusId)
		_, err = tx.Exec(ctx, deleteResult.Query, deleteResult.Args...)
		if err != nil {
			s.log.Error("un-handled delete syllabus query error", logger.Err(err))
			return isVerified, util.ErrInternal
		}

		isVerified = false
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit verify syllabus transaction", logger.Err(err))
		return false, util.ErrInternal
	}

	return isVerified, nil
}

func (s *pgSyllabusRepository) getDateSyncedSyllabusQuery(syllabusId string) (util.SqlBuilderResult, error) {
	syllabusUuid, err := database.ParsePgUuid(syllabusId)
	if err != nil {
		return util.SqlBuilderResult{}, err
	}

	qb := util.NewSqlBuilder("select date_synced from syllabi")
	qb.Concat("where id = $%d", syllabusUuid)

	return qb.Result(), nil
}

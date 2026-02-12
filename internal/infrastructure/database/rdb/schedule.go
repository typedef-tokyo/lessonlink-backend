package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type Schedule struct {
	c *sql.DB
}

func NewScheduleRepository(c IMySQL) repository.ScheduleRepository {
	return &Schedule{c: c.GetConn()}
}

func (f *Schedule) Save(ctx context.Context, tx *sql.Tx, rootModel *schedule.RootScheduleModel) (vo.ScheduleID, error) {

	scheduleID := vo.SCHEDULE_ID_INVALID

	// 新規登録の場合
	if rootModel.ID().IsInitial() {

		// スケジュールを登録
		scheduleDTO := f.toScheduleDTO(rootModel)
		err := scheduleDTO.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		scheduleID, err = vo.NewScheduleID(scheduleDTO.ID)
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		err = f.toItemBulkInsert(ctx, tx, scheduleID, rootModel.HistoryIndex(), rootModel.Items())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		err = f.toRoomItemBulkInsert(ctx, tx, scheduleID, rootModel.HistoryIndex(), rootModel.RoomItems())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

	} else {

		// 既存のデータを取得
		existsRecord, err := f.findByIDWithHistoryIndex(ctx, tx, rootModel.ID(), rootModel.HistoryIndex())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTrace(err)
		}

		if existsRecord == nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
		}

		// スケジュールを更新
		scheduleDTO := f.toScheduleDTO(rootModel)
		_, err = scheduleDTO.Update(ctx, tx, boil.Whitelist(
			dto.TBLScheduleColumns.Title,
			dto.TBLScheduleColumns.HistoryIndex,
			dto.TBLScheduleColumns.StartTime,
			dto.TBLScheduleColumns.EndTime,
			dto.TBLScheduleColumns.LastUpdateUser,
			dto.TBLScheduleColumns.UpdatedAt,
		))

		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		// 既存のアイテムを削除する
		if existsRecord.R != nil {

			_, err = tx.Exec(
				"DELETE FROM tbl_schedule_items WHERE schedule_id = ? AND history_index >= ? ",
				scheduleDTO.ID,
				scheduleDTO.HistoryIndex,
			)
			if err != nil {
				return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
			}

			_, err = tx.Exec(
				"DELETE FROM tbl_schedule_room_items WHERE schedule_id = ? AND history_index >= ? ",
				scheduleDTO.ID,
				scheduleDTO.HistoryIndex,
			)
			if err != nil {
				return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
			}
		}

		err = f.toItemBulkInsert(ctx, tx, rootModel.ID(), rootModel.HistoryIndex(), rootModel.Items())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		err = f.toRoomItemBulkInsert(ctx, tx, rootModel.ID(), rootModel.HistoryIndex(), rootModel.RoomItems())
		if err != nil {
			return vo.SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceInternalServerError(err)
		}

		scheduleID = rootModel.ID()
	}

	return scheduleID, nil
}

func (f *Schedule) Delete(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, deleteUserID vo.UserID) error {

	// 既存のデータを取得
	existsRecord, err := f.findByID(ctx, scheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if existsRecord.R != nil {

		_, err := existsRecord.R.ScheduleTBLScheduleItems.DeleteAll(ctx, tx)
		if err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}

		_, err = existsRecord.R.ScheduleTBLScheduleRoomItems.DeleteAll(ctx, tx)
		if err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}
	}

	// スケジュールを更新
	_, err = existsRecord.Delete(ctx, tx)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Schedule) findByIDWithHistoryIndex(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex) (*dto.TBLSchedule, error) {

	query := dto.TBLSchedules(
		dto.TBLScheduleWhere.ID.EQ(scheduleID.Value()),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleItems, dto.TBLScheduleItemWhere.HistoryIndex.EQ(historyIndex.Value())),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleRoomItems, dto.TBLScheduleRoomItemWhere.HistoryIndex.EQ(historyIndex.Value())),
	)

	var scheduleRecords *dto.TBLSchedule
	var err error
	if tx == nil {
		scheduleRecords, err = query.One(ctx, f.c)
	} else {
		scheduleRecords, err = query.One(ctx, tx)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return scheduleRecords, nil
}

func (f *Schedule) findByID(ctx context.Context, scheduleID vo.ScheduleID) (*dto.TBLSchedule, error) {

	record, err := dto.TBLSchedules(
		dto.TBLScheduleWhere.ID.EQ(scheduleID.Value()),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleItems),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleRoomItems),
	).One(ctx, f.c)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return record, nil
}

func (f *Schedule) FindByIDWithLock(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID) (*schedule.RootScheduleModel, error) {

	scheduleRecord, err := dto.TBLSchedules(
		dto.TBLScheduleWhere.ID.EQ(scheduleID.Value()),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleItems),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleRoomItems),
		qm.For("UPDATE"),
	).One(ctx, tx)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if scheduleRecord == nil {
		return nil, nil
	}

	model, err := f.toModel(scheduleRecord)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil
}

func (f *Schedule) FindByIDWithLockHistoryIndex(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex) (*schedule.RootScheduleModel, error) {

	scheduleRecord, err := dto.TBLSchedules(
		dto.TBLScheduleWhere.ID.EQ(scheduleID.Value()),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleItems, dto.TBLScheduleItemWhere.HistoryIndex.EQ(historyIndex.Value())),
		qm.Load(dto.TBLScheduleRels.ScheduleTBLScheduleRoomItems, dto.TBLScheduleRoomItemWhere.HistoryIndex.EQ(historyIndex.Value())),
		qm.For("UPDATE"),
	).One(ctx, tx)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if scheduleRecord == nil {
		return nil, nil
	}

	model, err := f.toModel(scheduleRecord)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil
}

func (f *Schedule) FindByIDWithHistoryIndex(ctx context.Context, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex) (*schedule.RootScheduleModel, error) {

	scheduleRecord, err := f.findByIDWithHistoryIndex(ctx, nil, scheduleID, historyIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if scheduleRecord == nil {
		return nil, nil
	}

	model, err := f.toModel(scheduleRecord)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil
}

func (f *Schedule) FindByID(ctx context.Context, scheduleID vo.ScheduleID) (*schedule.RootScheduleModel, error) {

	scheduleDTO, err := f.findByID(ctx, scheduleID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if scheduleDTO == nil {
		return nil, nil
	}

	historyIndex, err := vo.NewHistoryIndex(scheduleDTO.HistoryIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	scheduleLatestDTO, err := f.findByIDWithHistoryIndex(ctx, nil, scheduleID, historyIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	latestModel, err := f.toModel(scheduleLatestDTO)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return latestModel, nil
}

func (f *Schedule) toScheduleDTO(root *schedule.RootScheduleModel) *dto.TBLSchedule {

	startTime, endTime := root.ScheduleTime().Value()

	return &dto.TBLSchedule{
		ID:             root.ID().Value(),
		Campus:         root.Campus().Value(),
		Title:          root.Title().Value(),
		HistoryIndex:   root.HistoryIndex().Value(),
		StartTime:      startTime,
		EndTime:        endTime,
		CreateUser:     root.CreateUser().Value(),
		LastUpdateUser: root.LastUpdateUser().Value(),
		UpdatedAt:      time.Now(),
		CreatedAt:      time.Now(),
	}
}

func (f *Schedule) toItemBulkInsert(ctx context.Context, tx *sql.Tx, sheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, items schedule.ScheduleItemModelSlice) error {

	if len(items) == 0 {
		return nil
	}

	const COLUMUN_COUNT = 5

	placeholders := make([]string, len(items))
	values := make([]any, len(items)*COLUMUN_COUNT)

	for index, item := range items {

		placeholders[index] = "(?, ?, ?, ?, ?)"

		counter := 0
		values[index*COLUMUN_COUNT+counter] = sheduleID.Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = historyIndex.Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = item.LessonID().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = item.Identifier().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = item.Duration().Value()
	}

	query := fmt.Sprintf(`
	INSERT INTO tbl_schedule_items(
		schedule_id,
		history_index,
		lesson_id,
		identifier,
		duration
	)
	VALUES %s `, strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Schedule) toRoomItemBulkInsert(ctx context.Context, tx *sql.Tx, sheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, roomItems schedule.ScheduleRoomItemModelSlice) error {

	if len(roomItems) == 0 {
		return nil
	}

	const COLUMUN_COUNT = 11

	placeholders := make([]string, len(roomItems))
	values := make([]any, len(roomItems)*COLUMUN_COUNT)

	for index, roomItem := range roomItems {

		placeholders[index] = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		startTimeHoue, startTimeMinute := roomItem.StartTime().Value()
		endTimeHoue, endTimeMinute := roomItem.EndTime().Value()

		counter := 0
		values[index*COLUMUN_COUNT+counter] = sheduleID.Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = historyIndex.Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = roomItem.ItemTag().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = roomItem.LessonID().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = roomItem.Identifier().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = roomItem.Duration().Value()
		counter++
		values[index*COLUMUN_COUNT+counter] = startTimeHoue
		counter++
		values[index*COLUMUN_COUNT+counter] = startTimeMinute
		counter++
		values[index*COLUMUN_COUNT+counter] = endTimeHoue
		counter++
		values[index*COLUMUN_COUNT+counter] = endTimeMinute
		counter++
		values[index*COLUMUN_COUNT+counter] = roomItem.RoomIndex().Value()
	}

	query := fmt.Sprintf(`
	INSERT INTO tbl_schedule_room_items(
		schedule_id,
		history_index,
		item_tag,
		lesson_id,
		identifier,
		duration,
		start_time_hour,
		start_time_minutes,
		end_time_hour,
		end_time_minutes,
		room_index
	)
	VALUES %s `, strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Schedule) toModel(record *dto.TBLSchedule) (*schedule.RootScheduleModel, error) {

	var id vo.ScheduleID
	var campus vo.Campus
	var title vo.ScheduleTitle
	var historyIndex vo.HistoryIndex
	var createUser vo.UserID
	var lastUpdateUser vo.UserID
	items := []*schedule.ScheduleItemModel{}
	roomItems := []*schedule.ScheduleRoomItemModel{}

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&id, vo.NewScheduleID, record.ID))
	errs = errors.Join(errs, vo.SetVOConstructor(&campus, vo.NewCampus, record.Campus))
	errs = errors.Join(errs, vo.SetVOConstructor(&title, vo.NewScheduleTitle, record.Title))
	errs = errors.Join(errs, vo.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, record.HistoryIndex))
	errs = errors.Join(errs, vo.SetVOConstructor(&createUser, vo.NewUserID, record.CreateUser))
	errs = errors.Join(errs, vo.SetVOConstructor(&lastUpdateUser, vo.NewUserID, record.LastUpdateUser))

	scheduleTime, err := vo.NewScheduleTime(record.StartTime, record.EndTime)
	errs = errors.Join(errs, err)

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	if record.R != nil {

		items = make([]*schedule.ScheduleItemModel, 0, len(record.R.ScheduleTBLScheduleItems))
		for _, recordItem := range record.R.ScheduleTBLScheduleItems {

			var lessonID vo.LessonID
			var identifier vo.Identifier
			var duration vo.LessonDuration

			var errs error
			errs = errors.Join(errs, vo.SetVOConstructor(&lessonID, vo.NewLessonID, recordItem.LessonID))
			errs = errors.Join(errs, vo.SetVOConstructor(&identifier, vo.NewIdentifier, recordItem.Identifier))
			errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, recordItem.Duration))

			if errs != nil {
				return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
			}

			items = append(items, schedule.NewScheduleItemModel(
				lessonID,
				identifier,
				duration,
			))
		}

		roomItems = make([]*schedule.ScheduleRoomItemModel, 0, len(record.R.ScheduleTBLScheduleRoomItems))
		for _, recordItem := range record.R.ScheduleTBLScheduleRoomItems {

			var itemTag vo.RoomItemTag
			var lessonID vo.LessonID
			var identifier vo.Identifier
			var duration vo.LessonDuration
			var roomIndex vo.RoomIndex

			var errs error
			errs = errors.Join(errs, vo.SetVOConstructor(&itemTag, vo.NewRoomItemTag, recordItem.ItemTag))
			errs = errors.Join(errs, vo.SetVOConstructor(&lessonID, vo.NewLessonID, recordItem.LessonID))
			errs = errors.Join(errs, vo.SetVOConstructor(&identifier, vo.NewIdentifier, recordItem.Identifier))
			errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, recordItem.Duration))

			startTime, err := vo.NewScheduleLessonTime(recordItem.StartTimeHour, recordItem.StartTimeMinutes)
			errs = errors.Join(errs, err)

			endTime, err := vo.NewScheduleLessonTime(recordItem.EndTimeHour, recordItem.EndTimeMinutes)
			errs = errors.Join(errs, err)

			errs = errors.Join(errs, vo.SetVOConstructor(&roomIndex, vo.NewRoomIndex, recordItem.RoomIndex))

			if errs != nil {
				return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
			}

			roomItems = append(roomItems, schedule.NewScheduleRoomItemModel(
				itemTag,
				lessonID,
				identifier,
				duration,
				startTime,
				endTime,
				roomIndex,
			))
		}
	}

	return schedule.NewRootScheduleModel(
		id,
		campus,
		title,
		historyIndex,
		createUser,
		lastUpdateUser,
		items,
		roomItems,
		scheduleTime,
		record.CreatedAt,
		record.UpdatedAt,
	), nil

}

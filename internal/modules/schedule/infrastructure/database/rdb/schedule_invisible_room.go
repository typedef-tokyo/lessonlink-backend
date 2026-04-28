package rdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/invisible"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type ScheduleInvisibleRoom struct {
	c *sql.DB
}

func NewScheduleInvisibleRoomRepository(c rdb.IMySQL) repository.ScheduleInvisibleRoomRepository {
	return &ScheduleInvisibleRoom{c: c.GetConn()}
}

func (f *ScheduleInvisibleRoom) Save(ctx context.Context, tx *sql.Tx, sheduleID vo.ScheduleID, models []*invisible.RootScheduleInvisibleRoomModel) error {

	records, err := f.findBySheduleID(ctx, sheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	_, err = dto.TBLScheduleInvisibleRoomSlice(records).DeleteAll(ctx, tx)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	for _, model := range models {
		if err := f.toDTO(model).Insert(ctx, tx, boil.Infer()); err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}
	}

	return nil
}

func (f *ScheduleInvisibleRoom) FindBySheduleID(ctx context.Context, sheduleID vo.ScheduleID) (invisible.RootScheduleInvisibleRoomModelSlice, error) {

	records, err := f.findBySheduleID(ctx, sheduleID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	models := make([]*invisible.RootScheduleInvisibleRoomModel, 0, len(records))
	for _, record := range records {

		model, err := f.toModel(record)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		models = append(models, model)
	}

	return models, nil
}

func (f *ScheduleInvisibleRoom) findBySheduleID(ctx context.Context, sheduleID vo.ScheduleID) ([]*dto.TBLScheduleInvisibleRoom, error) {

	records, err := dto.TBLScheduleInvisibleRooms(
		dto.TBLScheduleInvisibleRoomWhere.ScheduleID.EQ(sheduleID.Value()),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return records, nil
}

func (f *ScheduleInvisibleRoom) toModel(record *dto.TBLScheduleInvisibleRoom) (*invisible.RootScheduleInvisibleRoomModel, error) {

	var scheduleID vo.ScheduleID
	var roomIndex vo.RoomIndex

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, record.ScheduleID))
	errs = errors.Join(errs, utility.SetVOConstructor(&roomIndex, vo.NewRoomIndex, record.RoomIndex))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	return invisible.NewRootScheduleInvisibleRoomModel(
		scheduleID,
		roomIndex,
	), nil

}

func (f *ScheduleInvisibleRoom) toDTO(model *invisible.RootScheduleInvisibleRoomModel) *dto.TBLScheduleInvisibleRoom {

	return &dto.TBLScheduleInvisibleRoom{
		ScheduleID: model.ScheduleID().Value(),
		RoomIndex:  model.RoomIndex().Value(),
	}
}

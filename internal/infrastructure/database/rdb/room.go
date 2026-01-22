package rdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type Room struct {
	c *sql.DB
}

func NewRoomRepository(c IMySQL) repository.RoomRepository {
	return &Room{c: c.GetConn()}
}

func (f *Room) Save(ctx context.Context, tx *sql.Tx, campus vo.Campus, slice room.RootRoomModelSlice) error {

	roomDTOs, err := f.findByCampus(ctx, tx, campus)

	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	// 全て削除
	_, err = roomDTOs.DeleteAll(ctx, tx)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	for _, model := range slice {

		err := f.toDTO(model).Insert(ctx, tx, boil.Infer())
		if err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}
	}

	return nil
}

func (f *Room) FindByCampus(ctx context.Context, campus vo.Campus) (room.RootRoomModelSlice, error) {

	records, err := f.findByCampus(ctx, nil, campus)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	roomModels := make([]*room.RootRoomModel, 0, len(records))
	for _, dto := range records {

		model, err := f.toModel(dto)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		roomModels = append(roomModels, model)
	}

	return roomModels, nil
}

func (f *Room) findByCampus(ctx context.Context, tx *sql.Tx, campus vo.Campus) (dto.DataRoomSlice, error) {

	query := dto.DataRooms(
		dto.DataRoomWhere.Campus.EQ(campus.Value()),
	)

	var records dto.DataRoomSlice
	var err error
	if tx != nil {
		records, err = query.All(ctx, tx)
	} else {
		records, err = query.All(ctx, f.c)
	}

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return records, nil
}

func (f *Room) toModel(record *dto.DataRoom) (*room.RootRoomModel, error) {

	var errs error

	var campus vo.Campus
	var index vo.RoomIndex
	var name vo.RoomName

	errs = errors.Join(errs, vo.SetVOConstructor(&campus, vo.NewCampus, record.Campus))
	errs = errors.Join(errs, vo.SetVOConstructor(&index, vo.NewRoomIndex, record.RoomIndex))
	errs = errors.Join(errs, vo.SetVOConstructor(&name, vo.NewRoomName, record.Name))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	return room.NewRootRoomModel(
		campus,
		index,
		name,
	), nil
}

func (f *Room) toDTO(model *room.RootRoomModel) *dto.DataRoom {

	return &dto.DataRoom{
		Campus:    model.Campus().Value(),
		RoomIndex: model.RoomIndex().Value(),
		Name:      model.RoomName().Value(),
	}
}

package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IRoomEditInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, campus string, editRoom RoomsEditInputDTO) error
	}
)

type (
	RoomsEditInputDTO struct {
		Rooms []RoomEditInputDTO
	}

	RoomEditInputDTO struct {
		ID    int
		Index int
		Name  string
	}
)

///////////////

type (
	RoomEditInteractor struct {
		txManager      util.TxManager
		repositoryRoom repository.RoomRepository
	}
)

func NewRoomEditInteractor(
	txManager util.TxManager,
	repositoryRoom repository.RoomRepository,
) IRoomEditInputPort {
	return &RoomEditInteractor{
		txManager:      txManager,
		repositoryRoom: repositoryRoom,
	}
}

func (r RoomEditInteractor) Execute(ctx context.Context, role vo.RoleKey, inputCampus string, editRoom RoomsEditInputDTO) error {

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	campus, roomSlice, err := r.createModel(inputCampus, editRoom)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if !roomSlice.IsUniq() {
		return log.WrapErrorWithStackTraceBadRequest(log.Errorf("教室情報が重複しています"))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.repositoryRoom.Save(ctx, tx, campus, roomSlice); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

func (r RoomEditInteractor) createModel(inputCampus string, _editRoom RoomsEditInputDTO) (vo.Campus, room.RootRoomModelSlice, error) {

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return campus, nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	models := make([]*room.RootRoomModel, 0, len(_editRoom.Rooms))

	var errs error

	for _, editRoom := range _editRoom.Rooms {

		var index vo.RoomIndex
		var name vo.RoomName

		errs = errors.Join(errs, vo.SetVOConstructor(&index, vo.NewRoomIndex, editRoom.Index))
		errs = errors.Join(errs, vo.SetVOConstructor(&name, vo.NewRoomName, editRoom.Name))

		models = append(models, room.NewRootRoomModel(
			campus,
			index,
			name,
		))
	}

	if errs != nil {
		return campus, nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return campus, models, nil
}

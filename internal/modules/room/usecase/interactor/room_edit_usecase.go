package interactor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type (
	IRoomEditInputPort interface {
		Execute(ctx context.Context, inputRole string, campus string, editRoom RoomsEditInputDTO) error
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
		txManager      database.TxManager
		repositoryRoom repository.RoomRepository
		facadeCampus   external.ICampusExistFacade
	}
)

func NewRoomEditInteractor(
	txManager database.TxManager,
	repositoryRoom repository.RoomRepository,
	facadeCampus external.ICampusExistFacade,
) IRoomEditInputPort {
	return &RoomEditInteractor{
		txManager:      txManager,
		repositoryRoom: repositoryRoom,
		facadeCampus:   facadeCampus,
	}
}

func (r RoomEditInteractor) Execute(ctx context.Context, inputRole string, inputCampus string, editRoom RoomsEditInputDTO) error {

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	campus, roomSlice, err := r.createModel(inputCampus, editRoom)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	campusExist, err := r.facadeCampus.Execute(ctx, inputCampus)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if !campusExist {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定した校舎はありません:%s", campus.Value()))
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

		errs = errors.Join(errs, utility.SetVOConstructor(&index, vo.NewRoomIndex, editRoom.Index))
		errs = errors.Join(errs, utility.SetVOConstructor(&name, vo.NewRoomName, editRoom.Name))

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

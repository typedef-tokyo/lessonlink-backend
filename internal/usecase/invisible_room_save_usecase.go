package usecase

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/invisible"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IInvisibleRoomSaveInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, inputUserID vo.UserID, inputScheduleID int, roomIndexes []int) error
	}
)

type InvisibleRoomSaveInteractor struct {
	txManager                       util.TxManager
	repositorySchedule              repository.ScheduleRepository
	repositoryUser                  repository.UserRepository
	repositoryRoom                  repository.RoomRepository
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository
	serviceScheduleEditPermission   service.IScheduleEditPermissionService
}

func NewInvisibleRoomSaveInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	repositoryRoom repository.RoomRepository,
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IInvisibleRoomSaveInputPort {
	return &InvisibleRoomSaveInteractor{
		txManager:                       txManager,
		repositorySchedule:              repositorySchedule,
		repositoryUser:                  repositoryUser,
		repositoryRoom:                  repositoryRoom,
		repositoryScheduleInvisibleRoom: repositoryScheduleInvisibleRoom,
		serviceScheduleEditPermission:   serviceScheduleEditPermission,
	}
}

func (r InvisibleRoomSaveInteractor) Execute(ctx context.Context, role vo.RoleKey, inputUserID vo.UserID, inputScheduleID int, inputRoomIndexes []int) error {

	scheduleID, err := vo.NewScheduleID(inputScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	scheduleData, err := r.repositorySchedule.FindByID(ctx, scheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	editUser, err := r.repositoryUser.FindByUserID(ctx, inputUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(scheduleData, editUser)
	if !isEnable {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	allRooms, err := r.repositoryRoom.FindByCampus(ctx, scheduleData.Campus())
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	invisibleRooms := make([]*invisible.RootScheduleInvisibleRoomModel, 0, len(inputRoomIndexes))
	for _, inputRoomIndex := range inputRoomIndexes {

		roomIndex, err := vo.NewRoomIndex(inputRoomIndex)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		if !allRooms.IsExist(scheduleData.Campus(), roomIndex) {
			return log.WrapErrorWithStackTraceBadRequest(log.Errorf("指定した教室番号は存在しません 番号:%d", roomIndex.Value()))
		}

		invisibleRooms = append(invisibleRooms, invisible.NewRootScheduleInvisibleRoomModel(scheduleData.ID(), roomIndex))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.repositoryScheduleInvisibleRoom.Save(ctx, tx, scheduleID, invisibleRooms); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

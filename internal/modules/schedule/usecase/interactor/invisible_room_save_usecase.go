package interactor

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/invisible"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type (
	IInvisibleRoomSaveInputPort interface {
		Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputRoomIndexes []int) error
	}
)

type InvisibleRoomSaveInteractor struct {
	txManager                       database.TxManager
	repositorySchedule              repository.ScheduleRepository
	facadeGetUser                   external.IUserGetFacade
	facadeExistsRoom                external.IRoomExistsFacade
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository
	serviceScheduleEditPermission   service.IScheduleEditPermissionService
}

func NewInvisibleRoomSaveInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	facadeExistsRoom external.IRoomExistsFacade,
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IInvisibleRoomSaveInputPort {
	return &InvisibleRoomSaveInteractor{
		txManager:                       txManager,
		repositorySchedule:              repositorySchedule,
		facadeGetUser:                   facadeGetUser,
		facadeExistsRoom:                facadeExistsRoom,
		repositoryScheduleInvisibleRoom: repositoryScheduleInvisibleRoom,
		serviceScheduleEditPermission:   serviceScheduleEditPermission,
	}
}

func (r InvisibleRoomSaveInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputRoomIndexes []int) error {

	scheduleID, err := vo.NewScheduleID(inputScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	scheduleData, err := r.repositorySchedule.FindByID(ctx, scheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if scheduleData == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
	}

	editUserID, err := vo.NewUserID(inputUserID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	_, outEditUserRole, err := r.facadeGetUser.Execute(ctx, editUserID.Value())
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	editUserRole, err := vo.NewRoleKey(outEditUserRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(scheduleData, editUserID, editUserRole)
	if !isEnable {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	// 全て存在する教室かチェック
	err = r.facadeExistsRoom.Execute(ctx, scheduleData.Campus().Value(), inputRoomIndexes)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	invisibleRooms := make([]*invisible.RootScheduleInvisibleRoomModel, 0, len(inputRoomIndexes))
	for _, inputRoomIndex := range inputRoomIndexes {

		roomIndex, err := vo.NewRoomIndex(inputRoomIndex)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
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

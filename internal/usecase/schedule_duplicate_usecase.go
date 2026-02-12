package usecase

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IScheduleDuplicatePort interface {
		Execute(ctx context.Context, roleKey vo.RoleKey, duplicateUser vo.UserID, inputDuplicateScheduleID int) error
	}
)

type (
	ScheduleDuplicateInteractor struct {
		txManager          util.TxManager
		repositorySchedule repository.ScheduleRepository
	}
)

func NewScheduleDuplicateInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
) IScheduleDuplicatePort {
	return &ScheduleDuplicateInteractor{
		txManager:          txManager,
		repositorySchedule: repositorySchedule,
	}
}

func (r ScheduleDuplicateInteractor) Execute(ctx context.Context, roleKey vo.RoleKey, duplicateUser vo.UserID, inputDuplicateScheduleID int) error {

	if roleKey.IsViewer() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	duplicateScheduleID, err := vo.NewScheduleID(inputDuplicateScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	schedule, err := r.repositorySchedule.FindByID(ctx, duplicateScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if schedule == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", duplicateScheduleID.Value()))
	}

	duplicateSchedule := schedule.Duplicate(duplicateUser)

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		_, err = r.repositorySchedule.Save(ctx, tx, duplicateSchedule)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

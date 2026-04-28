package interactor

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type (
	IScheduleDuplicatePort interface {
		Execute(ctx context.Context, inputRole string, inputDuplicateUser int, inputDuplicateScheduleID int) error
	}
)

type (
	ScheduleDuplicateInteractor struct {
		txManager          database.TxManager
		repositorySchedule repository.ScheduleRepository
	}
)

func NewScheduleDuplicateInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
) IScheduleDuplicatePort {
	return &ScheduleDuplicateInteractor{
		txManager:          txManager,
		repositorySchedule: repositorySchedule,
	}
}

func (r ScheduleDuplicateInteractor) Execute(ctx context.Context, inputRole string, inputDuplicateUser int, inputDuplicateScheduleID int) error {

	duplicateUser, err := vo.NewUserID(inputDuplicateUser)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if role.IsViewer() {
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

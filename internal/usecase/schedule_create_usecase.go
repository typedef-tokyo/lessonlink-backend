package usecase

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IScheduleCreateInputPort interface {
		Execute(
			ctx context.Context,
			role vo.RoleKey,
			registerUser vo.UserID,
			inputCampus string,
			inputStartTime int,
			inputEndTime int,
		) (*ScheduleCreateOutput, error)
	}
)

type (
	ScheduleCreateOutput struct {
		ScheduleID int
	}
)

type (
	ScheduleCreateInteractor struct {
		txManager          util.TxManager
		repositoryCampus   repository.CampusRepository
		repositorySchedule repository.ScheduleRepository
		repositoryLesson   repository.LessonRepository
	}
)

func NewScheduleCreateInteractor(
	txManager util.TxManager,
	repositoryCampus repository.CampusRepository,
	repositorySchedule repository.ScheduleRepository,
	repositoryLesson repository.LessonRepository,

) IScheduleCreateInputPort {
	return &ScheduleCreateInteractor{
		txManager:          txManager,
		repositoryCampus:   repositoryCampus,
		repositorySchedule: repositorySchedule,
		repositoryLesson:   repositoryLesson,
	}
}

func (r ScheduleCreateInteractor) Execute(
	ctx context.Context,
	role vo.RoleKey,
	registerUser vo.UserID,
	inputCampus string,
	inputStartTime int,
	inputEndTime int,
) (*ScheduleCreateOutput, error) {

	if role.IsEditor() {
		return nil, log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	scheduleTime, err := vo.NewScheduleTime(
		inputStartTime,
		inputEndTime,
	)

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	campuses, err := r.repositoryCampus.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if !campuses.IsExist(campus) {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定した校舎はありません:%s", campus.Value()))
	}

	var sheduleID vo.ScheduleID
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		rootSchedule := schedule.NewCreateRootScheduleModel(
			campus,
			registerUser,
			scheduleTime,
		)

		sheduleID, err = r.repositorySchedule.Save(ctx, tx, rootSchedule)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &ScheduleCreateOutput{ScheduleID: sheduleID.Value()}, nil
}

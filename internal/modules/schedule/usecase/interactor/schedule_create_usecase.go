package interactor

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type (
	IScheduleCreateInputPort interface {
		Execute(
			ctx context.Context,
			inputRole string,
			inputRegisterUser int,
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
		txManager          database.TxManager
		facadeCampus       external.ICampusExistFacade
		repositorySchedule repository.ScheduleRepository
	}
)

func NewScheduleCreateInteractor(
	txManager database.TxManager,
	facadeCampus external.ICampusExistFacade,
	repositorySchedule repository.ScheduleRepository,
) IScheduleCreateInputPort {
	return &ScheduleCreateInteractor{
		txManager:          txManager,
		facadeCampus:       facadeCampus,
		repositorySchedule: repositorySchedule,
	}
}

func (r ScheduleCreateInteractor) Execute(
	ctx context.Context,
	inputRole string,
	inputRegisterUser int,
	inputCampus string,
	inputStartTime int,
	inputEndTime int,
) (*ScheduleCreateOutput, error) {

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	if role.IsEditor() {
		return nil, log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	registerUser, err := vo.NewUserID(inputRegisterUser)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
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
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	campusExist, err := r.facadeCampus.Execute(ctx, inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if !campusExist {
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

package schedulelist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	campusRepository "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/campus"
)

type (
	IScheduleListQueryInputPort interface {
		Execute(ctx context.Context, campus string) (*ScheduleListQueryOutput, error)
	}
)

type (
	ScheduleListQueryOutput struct {
		ScheduleList []*QueryScheduleDTO
	}
)

type ScheduleListQueryInteractor struct {
	repositroryCampusQuery  campusRepository.CampusQueryRepository
	repositoryQuerySchedule ScheduleListQueryRepository
}

func NewScheduleListQueryInteractor(
	repositroryCampusQuery campusRepository.CampusQueryRepository,
	repositoryQuerySchedule ScheduleListQueryRepository,
) IScheduleListQueryInputPort {
	return &ScheduleListQueryInteractor{
		repositroryCampusQuery:  repositroryCampusQuery,
		repositoryQuerySchedule: repositoryQuerySchedule,
	}
}

func (r *ScheduleListQueryInteractor) Execute(ctx context.Context, campus string) (*ScheduleListQueryOutput, error) {

	campusModel, err := r.repositroryCampusQuery.GetByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if campusModel == nil {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したキャンパスはありません:%s", campus))
	}

	scheduleListDTO, err := r.repositoryQuerySchedule.GetListByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &ScheduleListQueryOutput{
		ScheduleList: scheduleListDTO,
	}, nil

}

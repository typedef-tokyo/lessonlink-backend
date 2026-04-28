package schedulelist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
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
	repositoryQuerySchedule ScheduleListQueryRepository
}

func NewScheduleListQueryInteractor(
	repositoryQuerySchedule ScheduleListQueryRepository,
) IScheduleListQueryInputPort {
	return &ScheduleListQueryInteractor{
		repositoryQuerySchedule: repositoryQuerySchedule,
	}
}

func (r *ScheduleListQueryInteractor) Execute(ctx context.Context, campus string) (*ScheduleListQueryOutput, error) {

	scheduleListDTO, err := r.repositoryQuerySchedule.GetListByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &ScheduleListQueryOutput{
		ScheduleList: scheduleListDTO,
	}, nil

}

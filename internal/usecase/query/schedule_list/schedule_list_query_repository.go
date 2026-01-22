package schedulelist

import (
	"context"
	"time"
)

type QueryScheduleDTO struct {
	ID                 int
	Campus             string
	Title              string
	Status             string
	CreateUserName     string
	LastUpdateUserName string
	CreateUser         int
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type ScheduleListQueryRepository interface {
	GetListByCampus(ctx context.Context, campus string) ([]*QueryScheduleDTO, error)
}

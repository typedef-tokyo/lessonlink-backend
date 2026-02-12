package presenter

import (
	"cmp"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	schedulelist "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/schedule_list"
)

type IScheduleListPresenter interface {
	Present(result *schedulelist.ScheduleListQueryOutput) *ScheduleListResponse
}

type ScheduleListPresenter struct {
}

func NewScheduleListPresenter(
	env configs.EnvConfig,
) IScheduleListPresenter {
	return &ScheduleListPresenter{}
}

type (
	ScheduleListResponse struct {
		Schedules []*ScheduleListDTO `json:"schedules"`
	}

	ScheduleListDTO struct {
		ScheduleID         int       `json:"schedule_id"`
		Title              string    `json:"title"`
		CreatedUserName    string    `json:"created_user_name"`
		LastUpdateUserName string    `json:"last_update_user_name"`
		LastUpdateDateTime time.Time `json:"last_update_date_time"`
		CreatedUserID      int       `json:"created_user_id"`
	}
)

func (h *ScheduleListPresenter) Present(result *schedulelist.ScheduleListQueryOutput) *ScheduleListResponse {

	scheduleList := lo.Map(result.ScheduleList, func(item *schedulelist.QueryScheduleDTO, _ int) *ScheduleListDTO {
		return &ScheduleListDTO{
			ScheduleID:         item.ID,
			Title:              item.Title,
			CreatedUserName:    item.CreateUserName,
			LastUpdateUserName: item.LastUpdateUserName,
			LastUpdateDateTime: item.UpdatedAt,
			CreatedUserID:      item.CreateUser,
		}
	})

	slices.SortFunc(scheduleList, func(a, b *ScheduleListDTO) int {
		return cmp.Compare(b.ScheduleID, a.ScheduleID)
	})

	return &ScheduleListResponse{
		Schedules: scheduleList,
	}
}

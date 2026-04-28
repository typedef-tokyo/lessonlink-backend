package service

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
)

type (
	IScheduleEditPermissionService interface {
		AllowsEditingBy(
			sheduleData *schedule.RootScheduleModel,
			editUserID vo.UserID,
			editUserRole vo.RoleKey,
		) bool
	}

	ScheduleEditPermissionService struct{}
)

func NewScheduleEditPermissionService() IScheduleEditPermissionService {
	return &ScheduleEditPermissionService{}
}

func (r ScheduleEditPermissionService) AllowsEditingBy(
	sheduleData *schedule.RootScheduleModel,
	editUserID vo.UserID,
	editUserRole vo.RoleKey,
) bool {

	if editUserRole.IsViewer() {
		return false
	}

	if editUserRole.IsEditor() && sheduleData.CreateUser() != editUserID {
		return false
	}

	return true

}

package service

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/user"
)

type (
	IScheduleEditPermissionService interface {
		AllowsEditingBy(
			sheduleData *schedule.RootScheduleModel,
			editUserID *user.RootUserModel,
		) bool
	}

	ScheduleEditPermissionService struct{}
)

func NewScheduleEditPermissionService() IScheduleEditPermissionService {
	return &ScheduleEditPermissionService{}
}

func (r ScheduleEditPermissionService) AllowsEditingBy(
	sheduleData *schedule.RootScheduleModel,
	editUser *user.RootUserModel,
) bool {

	if editUser.RoleKey().IsViewer() {
		return false
	}

	if editUser.RoleKey().IsEditor() && sheduleData.CreateUser() != editUser.ID() {
		return false
	}

	return true

}

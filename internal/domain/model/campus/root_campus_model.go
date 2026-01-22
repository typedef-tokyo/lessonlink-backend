package campus

import (
	"sort"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type CampusModelSlice []*RootCumpusModel

func (r CampusModelSlice) IsExist(campus vo.Campus) bool {

	_, found := lo.Find(r, func(item *RootCumpusModel) bool {
		return item.campus == campus
	})

	return found

}

func (r CampusModelSlice) SortByOrder() {

	sort.Slice(r, func(i, j int) bool {
		return r[i].orderIndex < r[j].orderIndex
	})
}

type RootCumpusModel struct {
	campus     vo.Campus
	campusName vo.CampusName
	orderIndex int
}

func NewRootCampus(
	campus vo.Campus,
	campusName vo.CampusName,
	orderIndex int,
) *RootCumpusModel {

	return &RootCumpusModel{
		campus:     campus,
		campusName: campusName,
		orderIndex: orderIndex,
	}
}

func (r RootCumpusModel) Campus() vo.Campus {
	return r.campus
}

func (r RootCumpusModel) CampusName() vo.CampusName {
	return r.campusName
}

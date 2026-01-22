package role

import (
	"sort"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type RootRoleModelSlice []*RootRoleModel

func (r RootRoleModelSlice) FindNameByKey(key vo.RoleKey) vo.RoleName {

	item, found := lo.Find(r, func(item *RootRoleModel) bool {
		return item.roleKey == key
	})

	if !found {
		return vo.ROLE_NAME_UNKNOWN
	}

	return item.roleName

}

func (r RootRoleModelSlice) Sort() {

	sort.Slice(r, func(i, j int) bool {
		return r[i].id < r[j].id
	})
}

func (r RootRoleModelSlice) IsOwner(roleKey vo.RoleKey) bool {

	item, found := lo.Find(r, func(item *RootRoleModel) bool {
		return item.roleKey == roleKey
	})

	if !found {
		return false
	}

	return item.roleKey == vo.ROLE_KEY_OWNER
}

type RootRoleModel struct {
	id       vo.RoleID
	roleKey  vo.RoleKey
	roleName vo.RoleName
}

func NewRootRoleModel(
	id vo.RoleID,
	roleKey vo.RoleKey,
	roleName vo.RoleName,
) *RootRoleModel {

	return &RootRoleModel{
		id:       id,
		roleKey:  roleKey,
		roleName: roleName,
	}
}

func (r RootRoleModel) ID() vo.RoleID {
	return r.id
}

func (r RootRoleModel) RoleKey() vo.RoleKey {
	return r.roleKey
}

func (r RootRoleModel) RoleName() vo.RoleName {
	return r.roleName
}

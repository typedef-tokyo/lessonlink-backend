package user

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/hash"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type RootUserModelSlice []*RootUserModel

func (r RootUserModelSlice) GetUserName(targetID vo.UserID) vo.UserDisplayName {

	user, found := lo.Find(r, func(item *RootUserModel) bool {
		return item.id == targetID
	})

	if !found {
		return "不明なユーザー"
	}

	return user.DisplayName()
}

func (r RootUserModelSlice) FindByUserID(userID vo.UserID) *RootUserModel {

	user, _ := lo.Find(r, func(item *RootUserModel) bool {
		return item.id == userID
	})

	return user
}

type RootUserModel struct {
	id          vo.UserID
	roleKey     vo.RoleKey
	userName    vo.UserName
	password    vo.UserPassword
	displayName vo.UserDisplayName
}

func NewRootUserModel(
	id vo.UserID,
	roleKey vo.RoleKey,
	userName vo.UserName,
	password vo.UserPassword,
	displayName vo.UserDisplayName,
) *RootUserModel {

	return &RootUserModel{
		id:          id,
		roleKey:     roleKey,
		userName:    userName,
		password:    password,
		displayName: displayName,
	}
}

func NewCreateUserModel(
	roleKey vo.RoleKey,
	userName vo.UserName,
	password vo.UserPassword,
	displayName vo.UserDisplayName,
) *RootUserModel {

	return &RootUserModel{
		id:          vo.USER_ID_INITIAL,
		roleKey:     roleKey,
		userName:    userName,
		password:    password,
		displayName: displayName,
	}
}

func (r RootUserModel) ID() vo.UserID {

	return r.id
}

func (r RootUserModel) RoleKey() vo.RoleKey {

	return r.roleKey
}

func (r RootUserModel) UserName() vo.UserName {

	return r.userName
}

func (r RootUserModel) DisplayName() vo.UserDisplayName {

	return r.displayName
}

func (r RootUserModel) Password() vo.UserPassword {

	return r.password
}

func (r RootUserModel) AuthenticatePassword(_password vo.UserPassword) bool {

	result, err := hash.VerifyPassword(_password.Value(), r.password.Value())

	if err != nil {
		fmt.Println(log.WrapErrorWithStackTrace(err))
		return false
	}

	return result
}

func (r RootUserModel) IsEnableDelete(deleteFromUserID vo.UserID) bool {

	return r.id == deleteFromUserID
}

func (r *RootUserModel) UpdateUser(
	newRoleKey vo.RoleKey,
	newPassword vo.UserPassword,
	newDisplayName vo.UserDisplayName,
	updateFromUserID vo.UserID,
	updateFromUserRoleKey vo.RoleKey,
) error {

	// 自分の情報以外はオーナーでないと変更できない
	if (r.id != updateFromUserID) && !(updateFromUserRoleKey.IsOwner()) {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	// 新しいロールをセット ロールが変更された場合はオーナー以外は設定できない
	if (r.roleKey != newRoleKey) && (!updateFromUserRoleKey.IsOwner()) {

		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	r.roleKey = newRoleKey
	r.displayName = newDisplayName

	// 新しいパスワードが設定されている場合はセット
	if !newPassword.IsNone() {
		r.password = newPassword
	}

	return nil
}

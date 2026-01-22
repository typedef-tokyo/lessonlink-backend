package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/user"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type UserRepository interface {
	FindAll(ctx context.Context) (user.RootUserModelSlice, error)
	FindByUserName(ctx context.Context, userName string) (*user.RootUserModel, error)
	FindByUserID(ctx context.Context, userId vo.UserID) (*user.RootUserModel, error)
	Save(ctx context.Context, tx *sql.Tx, user *user.RootUserModel, userID vo.UserID) error
	Delete(ctx context.Context, tx *sql.Tx, user *user.RootUserModel, deleteFromUserID vo.UserID) error
}

package repository

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/campus"
)

type CampusRepository interface {
	FindAll(ctx context.Context) (campus.CampusModelSlice, error)
}

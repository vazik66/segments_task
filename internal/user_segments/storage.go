package usersegments

import (
	"context"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID uint) (*[]UserSegment, error)
	AddSegmentsToUser(ctx context.Context, userID uint, slugs *[]string, ttl uint) error
	RemoveSegmentsFromUser(ctx context.Context, userID uint, slugs *[]string) (int, error)
    DeleteDeadUserSegments(ctx context.Context) (*[]DeadUserSegment, error)
}

package history

import (
	"context"
	"time"
)

type Repository interface {
	GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*[]HistoryRecord, error)
	CreateBulk(ctx context.Context, records *CreateHistoryRecordParams) error
}

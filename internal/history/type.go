package history

import (
	"time"
)

type Action string

const (
	Added   Action = "added"
	Removed Action = "removed"
)

type CreateHistoryRecordParams struct {
	UserID    uint      `mapstructure:"userId"`
	Slugs     *[]string `mapstructure:"slugs"`
	Action    Action    `mapstructure:"action"`
	CreatedAt time.Time `mapstructure:"createdAt"`
}

type HistoryRecord struct {
	UserID    uint      `json:"userId"`
	Slug      string    `json:"segmentSlug"`
	Action    Action    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateReportParams struct {
	UserID uint      `json:"userId"`
	Date   time.Time `json:"date" format:"date-time"`
}

type ReportUrlResponse struct {
	Url string `json:"url" example:"http://localhost:8000/files/3fb77cc1-1ffc-42d2-bfe3-01cfbf3b6e1d.csv"`
}

type CreateReportEventParams struct {
	Filename string
	UserID   uint
	Date     time.Time
}

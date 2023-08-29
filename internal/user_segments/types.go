package usersegments

import "time"

type UserSegment struct {
	UserID      uint      `json:"userId"`
	SegmentSlug string    `json:"segmentSlug"`
	Ttl         uint      `json:"ttl"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AddSegmentsToUserParams struct {
	UserID       uint      `json:"userId"`
	SegmentsSlug *[]string `json:"segmentsSlug"`
	Ttl          uint      `json:"ttl,omitempty"`
}

type RemoveSegmentsFromUserParams struct {
	UserID       uint      `json:"userId"`
	SegmentsSlug *[]string `json:"segmentsSlug"`
}

type GetUserSegmentsParams struct {
	UserID uint `json:"userId"`
}

type DeadUserSegment struct {
    UserID uint
    Slug string
    CreatedAt time.Time
}

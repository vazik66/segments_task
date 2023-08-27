package segment

import "context"

type Repository interface {
	Create(ctx context.Context, segment *CreateSegmentParams) (*Segment, error)
	GetBySlug(ctx context.Context, slug string) (*Segment, error)
	List(ctx context.Context) (*[]Segment, error)
	Update(ctx context.Context, segment *Segment) error
	Delete(ctx context.Context, slug string) error
}

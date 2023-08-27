package segment

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type SegmentService struct {
	repo Repository
}

func NewService(segmentRepo Repository) *SegmentService {
	return &SegmentService{
		repo: segmentRepo,
	}
}

func (s *SegmentService) Create(ctx context.Context, segment *CreateSegmentParams) (*Segment, error) {
	if segment.Slug == "" {
		return &Segment{}, fmt.Errorf("slug is empty")
	}
	segment.Slug = strings.ToUpper(segment.Slug)

	newSegment, err := s.repo.Create(ctx, segment)
	if err != nil {
		log.Printf("Could not create segment: %v. %v", segment, err)
		return &Segment{}, err
	}
	return newSegment, err
}

func (s *SegmentService) GetBySlug(ctx context.Context, slug *string) (*Segment, error) {
	if *slug == "" {
		return &Segment{}, fmt.Errorf("slug is empty")
	}
	*slug = strings.ToUpper(*slug)

	segment, err := s.repo.GetBySlug(ctx, *slug)
	if err != nil {
		log.Printf("Could not get segment by id: %s. %v", *slug, err)
		return &Segment{}, err
	}
	return segment, nil
}

func (s *SegmentService) List(ctx context.Context) (*[]Segment, error) {
	segments, err := s.repo.List(ctx)
	if err != nil {
		log.Printf("Could not list segments. %v", err)
		return nil, err
	}
	return segments, nil
}

func (s *SegmentService) Delete(ctx context.Context, slug *string) error {
	if *slug == "" {
		return fmt.Errorf("slug is empty")
	}
	*slug = strings.ToUpper(*slug)

	err := s.repo.Delete(ctx, *slug)
	if err != nil {
		log.Printf("Could not delete segments. %v", err)
		return err
	}
	return nil
}

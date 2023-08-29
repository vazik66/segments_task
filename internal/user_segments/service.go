package usersegments

import (
	"avito-segment/pkg/events"
	"avito-segment/internal/user"
	"context"
	"fmt"
	"log"
	"time"
)

type UserSegmentsService struct {
	repo     Repository
	userRepo user.Repository
	events   events.EventManager
}

func NewService(userSegmentRepo Repository, userRepo user.Repository, events events.EventManager) *UserSegmentsService {
    s := &UserSegmentsService{
		repo:     userSegmentRepo,
		userRepo: userRepo,
		events:   events,
	}
    s.postInit()
    return s
}

func (s *UserSegmentsService) postInit() {
    s.events.Subscribe("task.deleteUserSegmentsTTL", func(i interface{}) error { 
        log.Println("task delete ttl called")
        _, _ = s.DeleteExpired(context.Background())
        log.Println("task delete ttl finished")
        return nil
    })
}

func (s *UserSegmentsService) GetByUser(ctx context.Context, params *GetUserSegmentsParams) (*[]UserSegment, error) {
	if _, err := s.userRepo.GetByID(ctx, params.UserID); err != nil {
		return nil, fmt.Errorf("User does not exists")
	}

	userSegments, err := s.repo.GetByUserID(ctx, params.UserID)
	if err != nil {
		log.Printf("Could not get userSegments by userId: %d. %v", params.UserID, err)
		return nil, err
	}
	return userSegments, nil
}

func (s *UserSegmentsService) AddToUser(ctx context.Context, params *AddSegmentsToUserParams) error {
	if len(*params.SegmentsSlug) == 0 {
		return fmt.Errorf("Empty segments")
	}

	err := s.repo.AddSegmentsToUser(ctx, params.UserID, params.SegmentsSlug, params.Ttl)
	if err != nil {
		log.Printf("Could not add segments to user: %+v", params)
		return err
	}

	_ = s.events.Publish("usersegments.added", map[string]interface{}{
		"userId": params.UserID,
		"slugs":  params.SegmentsSlug,
	})
	return nil
}

func (s *UserSegmentsService) RemoveFromUser(ctx context.Context, params *RemoveSegmentsFromUserParams) error {
	if len(*params.SegmentsSlug) == 0 {
		return fmt.Errorf("Empty segments")
	}

	if _, err := s.userRepo.GetByID(ctx, params.UserID); err != nil {
		return fmt.Errorf("User does not exists")
	}

	removed, err := s.repo.RemoveSegmentsFromUser(ctx, params.UserID, params.SegmentsSlug)
	if err != nil {
		log.Printf("Could not remove segments: %v from user: %d", params.SegmentsSlug, params.UserID)
		return err
	}
    if removed == 0 {
        return fmt.Errorf("Slug does not exist")
    }

	_ = s.events.Publish("usersegments.deleted", map[string]interface{}{
		"userId": params.UserID,
		"slugs":  params.SegmentsSlug,
        "createdAt": time.Now().UTC(),
	})

	return nil
}

func (s *UserSegmentsService) DeleteExpired(ctx context.Context) (*[]DeadUserSegment, error) {
	deleted, err := s.repo.DeleteDeadUserSegments(ctx)
	if err != nil {
		return nil, err
	}
    
    _ = s.events.Publish("usersegments.ttlDeleted", deleted)

	return deleted, err
}

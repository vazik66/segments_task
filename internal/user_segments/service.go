package usersegments

import (
	"avito-segment/internal/user"
	"context"
	"fmt"
	"log"
)

type UserSegmentsService struct {
	repo Repository
    userRepo user.Repository
}

func NewService(userSegmentRepo Repository, userRepo user.Repository) *UserSegmentsService {
	return &UserSegmentsService{
		repo: userSegmentRepo,
        userRepo: userRepo,
	}
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

    err := s.repo.AddSegmentsToUser(ctx,params.UserID, params.SegmentsSlug, params.Ttl)
    if err != nil {
        log.Printf("Could not add segments to user: %+v", params)
        return err
    }

	return nil
}

func (s *UserSegmentsService) RemoveFromUser(ctx context.Context, params *RemoveSegmentsFromUserParams ) error {
	if len(*params.SegmentsSlug) == 0 {
		return fmt.Errorf("Empty segments")
	}

    if _, err := s.userRepo.GetByID(ctx, params.UserID); err != nil {
        return fmt.Errorf("User does not exists")
    }

    err := s.repo.RemoveSegmentsFromUser(ctx, params.UserID, params.SegmentsSlug)
    if err != nil {
        log.Printf("Could not remove segments: %v from user: %d", params.SegmentsSlug, params.UserID)
        return err
    }

	return nil
}

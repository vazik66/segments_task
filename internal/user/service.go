package user

import (
	"context"
	"log"
)

type UserService struct {
	repo Repository
}

func NewService(userRepo Repository) *UserService {
	return &UserService{
		repo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context) (*User, error) {
    user, err := s.repo.Create(ctx)
    if err != nil {
        log.Printf("Could not create user: %v. %v", user, err)
        return &User{}, err
    }
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, userID uint) (*User, error) {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        log.Printf("Could not get user with id: %d. %v", userID, err)
        return &User{}, err
    }
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, userID uint) error {
    err := s.repo.Delete(ctx, userID)
    if err != nil {
        log.Printf("Could not delete user with id: %d. %v", userID, err)
        return err
    }
    return nil
}

func (s *UserService) List(ctx context.Context) (*[]User, error) {
    users, err := s.repo.List(ctx)
    if err != nil {
        log.Printf("Error listing users. %v", err)
        return nil, err
    }
	return users, nil
}

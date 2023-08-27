package user

import "context"

type Repository interface {
	Create(ctx context.Context) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	List(ctx context.Context) (*[]User, error)
	Delete(ctx context.Context, id uint) error
}

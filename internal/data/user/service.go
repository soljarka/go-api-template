package user

import "context"

// Service is the interface for User service
type Service interface {
	Save(context.Context, *User) (ID, error)
	Find(context.Context, ID) (*User, error)
	All(context.Context) ([]*User, error)
	Delete(context.Context, ID) error
	DeleteAll(context.Context) error
}

type service struct {
	repo Repository
}

// NewService creates a new User service
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) Save(ctx context.Context, user *User) (ID, error) {
	id, err := s.repo.Save(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) Delete(ctx context.Context, id ID) error {
	err := s.repo.Delete(ctx, id)

	return err
}

func (s *service) DeleteAll(ctx context.Context) error {
	err := s.repo.DeleteAll(ctx)

	return err
}

func (s *service) Find(ctx context.Context, id ID) (*User, error) {
	user, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) All(ctx context.Context) ([]*User, error) {
	users, err := s.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

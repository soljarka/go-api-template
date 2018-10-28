package user

import "context"

// Repository is the user repository interface
type Repository interface {
	Save(context.Context, *User) (ID, error)
	Find(context.Context, ID) (*User, error)
	All(context.Context) ([]*User, error)
	Delete(context.Context, ID) error
	DeleteAll(context.Context) error
}

package core

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// userManager implements the UserManager interface.
type userManager struct {
	collection *mongo.Collection
}

// NewUserManager creates a new instance of the UserManager.
func NewUserManager(collection *mongo.Collection) UserManager {
	return &userManager{collection: collection}
}

func (u *userManager) SignIn(ctx context.Context, username string, password string) (*UserDoc, error) {
	panic("implement me")
}

func (u *userManager) SignUp(ctx context.Context, username string, password string) (string, error) {
	panic("implement me")
}

package subscription

import (
	"context"
	"gorm.io/gorm"
	"time"
)

func (s service) Create(ctx context.Context, req Request) (user User, err error) {
	user, err = s.postgres.CreateUser(User{
		ID:        1,
		Email:     "test@test.com",
		UserName:  "test",
		FullName:  "test pour",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	})
	if err != nil {
		return user, err
	}
	return user, nil
}

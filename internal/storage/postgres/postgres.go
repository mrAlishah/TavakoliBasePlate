package postgres

import (
	"GolangTraining/internal/subscription"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

var models = []interface{}{
	&subscription.User{},
}

func CreateRepository(db *gorm.DB) (*Repository, error) {
	repo := &Repository{
		database: db,
	}
	logrus.Infof("current db name: %s", db.Migrator().CurrentDatabase())
	err := db.AutoMigrate(models...)
	if err != nil {
		return repo, errors.Wrap(err, "failed to auto migrate models")
	}
	return repo, nil
}

func (r *Repository) CreateUser(user subscription.User) (subscription.User, error) {
	mu := subscription.User{
		Email:    user.Email,
		UserName: user.UserName,
		FullName: user.FullName,
	}
	err := r.database.Create(&mu).Error
	if err != nil {
		return mu, errors.Wrap(err, "failed to create a user")
	}
	return mu, nil
}

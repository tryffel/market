package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/logger"
	"github.com/tryffel/market/storage/repositories"
	"github.com/tryffel/market/storage/repository_impl"
)

type Store struct {
	db    *Database
	User  repositories.User
	Group repositories.Group
}

func NewStore(conf *config.Config, logger *logger.SqlLogger) (*Store, error) {
	var err error
	s := &Store{}
	s.db, err = NewDatabase(conf, logger)

	if err != nil {
		return s, Error.Wrap(&err, "failed to initialize database connection")
	}

	s.User = repository_impl.NewUserRepository(s.GetDbEngine())
	s.Group = repository_impl.NewGroupRepository(s.GetDbEngine())
	return s, nil
}

func (s *Store) GetDbEngine() *gorm.DB {
	return s.db.GetEngine()
}

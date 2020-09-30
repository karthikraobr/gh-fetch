package store

import (
	"log"
	"time"

	"github.com/karthikraobr/gh-fetch/internal/gh"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	db  *gorm.DB
	log *log.Logger
}

// Initializer the db store
func New(connectionString string, log *log.Logger) (*Store, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&gh.Repository{}, &gh.Commit{}); err != nil {
		return nil, err
	}
	log.Println("db init successful")
	return &Store{
		db:  db,
		log: log,
	}, nil
}

// DB represents database operations
type DB interface {
	GetRepository(id int64) (*gh.Repository, error)
	CreateRepository(r *gh.Repository) (*gh.Repository, error)
	GetRepositories(username string) ([]*gh.Repository, error)
	CreateRepositories(r []*gh.Repository) ([]*gh.Repository, error)
	GetRepositoriesOrderedBy(username string, limit int, sort string, sortBy string) ([]*gh.Repository, error)
}

// GetRepository fetches a single github repository by ID.
func (s *Store) GetRepository(id int64) (*gh.Repository, error) {
	var repo gh.Repository
	result := s.db.First(&repo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

// CreateRepository creates a repository.
func (s *Store) CreateRepository(r *gh.Repository) (*gh.Repository, error) {
	r.LastAccess = time.Now()
	result := s.db.Create(r)
	if result.Error != nil {
		return nil, result.Error
	}
	return r, nil
}

// GetRepositories fetches all the repository of a user.
func (s *Store) GetRepositories(username string) ([]*gh.Repository, error) {
	var repo []*gh.Repository
	result := s.db.Where("owner = ?", username).Find(&repo)
	if result.Error != nil {
		return nil, result.Error
	}
	return repo, nil
}

//GetRepositoriesOrderedBy fetches limit number of repositories sorted by sortby in sort order.
func (s *Store) GetRepositoriesOrderedBy(username string, limit int, sort string, sortBy string) ([]*gh.Repository, error) {
	var repo []*gh.Repository
	result := s.db.Order(sortBy+" "+sort).Limit(limit).Where("owner = ?", username).Find(&repo)
	if result.Error != nil {
		return nil, result.Error
	}
	return repo, nil
}

// CreateRepositories creates if not present or updates last access of repositories in a transaction.
func (s *Store) CreateRepositories(r []*gh.Repository) ([]*gh.Repository, error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range r {
			v.LastAccess = time.Now()
			var new gh.Repository
			if err := tx.Where(gh.Repository{ID: v.ID}).Attrs(*v).FirstOrCreate(&new).Error; err != nil {
				return err
			}
			if v.LastAccess != new.LastAccess {
				tx.Model(&new).Updates(gh.Repository{LastAccess: v.LastAccess})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return r, nil
}

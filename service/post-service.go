package service

import (
	"errors"
	"time"

	"github.com/jsparraq/api-rest/entity"
	"github.com/jsparraq/api-rest/repository"
)

// PostService is a interface
type PostService interface {
	Validate(post *entity.Post) error
	Create(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

type service struct{}

var (
	repo repository.PostRepository = repository.NewElasticsearchRepository()
)

// NewPostService function
func NewPostService() PostService {
	return &service{}
}

func (*service) Validate(post *entity.Post) error {
	if post == nil {
		err := errors.New("The post is empty")
		return err
	}
	if post.Message == "" {
		err := errors.New("The tweet message is empty")
		return err
	}
	return nil
}

func (*service) Create(post *entity.Post) (*entity.Post, error) {
	post.Created = time.Now()
	return repo.Save(post)
}

func (*service) FindAll() ([]entity.Post, error) {
	return repo.FindAll()
}

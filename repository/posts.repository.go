package repository

import "github.com/jsparraq/api-rest/entity"

// PostRepository interface
type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

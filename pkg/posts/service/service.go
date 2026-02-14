package service

import "github.com/amorindev/go-cms-tmpl/pkg/posts/port"

var _ port.PostSrv = &Service{}

type Service struct {
	PostRepo port.PostRepo
}

func NewPostSrv(postRepo port.PostRepo) *Service {
	return &Service{
		PostRepo: postRepo,
	}
}
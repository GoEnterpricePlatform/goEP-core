package handler

import "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"

type RespData struct {
	ActivePage  string
	ErrorMsg    string
	Posts       []*domain.Post
	EditingID   string
	CreateErr   string
	UpdateErr   string
	UpdateErrID string
	DeleteErr   string
	DeleteErrID string
	PageErr     string
}

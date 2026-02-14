package domain

import "time"

type Post struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Desc      *string    `json:"desc"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func FromCore(c Post) *Post {
	return &Post{
		ID:        c.ID,
		Desc:      c.Desc,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

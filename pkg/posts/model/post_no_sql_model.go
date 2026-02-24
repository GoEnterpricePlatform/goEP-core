package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostNoSqlModel struct {
	ID        bson.ObjectID `bson:"_id"`
	Title     string        `bson:"title"`
	Desc      *string       `bson:"desc,omitempty"`
	CreatedAt *time.Time    `bson:"created_at,omitempty"`
	UpdatedAt *time.Time    `bson:"updated_at,omitempty"`
}

func (m *PostNoSqlModel) ToDomain(p *domain.Post) {
	if m == nil {
		return
	}

	p.ID = m.ID.Hex()
	p.Title = m.Title
	p.Desc = m.Desc
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
}

func FromDomain(p *domain.Post, id bson.ObjectID) *PostNoSqlModel {
	if p == nil {
		return nil
	}

	return &PostNoSqlModel{
		ID:        id,
		Title:     p.Title,
		Desc:      p.Desc,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

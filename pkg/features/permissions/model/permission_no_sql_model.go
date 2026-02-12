package model

import (
	"github.com/amorindev/go-cms-tmpl/pkg/features/permissions/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PermissionNoSqlModel struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}

func (m *PermissionNoSqlModel) ToDomain(p *domain.Permission) {
	if m == nil {
		return
	}

	p.ID = m.ID.Hex()
	p.Name = m.Name
}

func FromDomain(p *domain.Permission, id bson.ObjectID) *PermissionNoSqlModel {
	if p == nil {
		return nil
	}

	return &PermissionNoSqlModel{
		ID:   id,
		Name: p.Name,
	}
}

func ToDomainList(models []*PermissionNoSqlModel) []*domain.Permission {
	if len(models) == 0 {
		return []*domain.Permission{}
	}

	permissions := make([]*domain.Permission, 0, len(models))

	for _, m := range models {
		if m == nil {
			continue
		}
		var permission domain.Permission
		m.ToDomain(&permission)
		permissions = append(permissions, &permission)
	}

	return permissions
}

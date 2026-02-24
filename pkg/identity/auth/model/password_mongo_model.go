package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/domain"
)

type UserPasswordNoSqlModel struct {
	PasswordHash string     `bson:"password_hash"`
	CreatedAt    *time.Time `bson:"created_at"`
	UpdatedAt    *time.Time `bson:"updated_at"`
}

func (m *UserPasswordNoSqlModel) ToDomain() *domain.UserPasswordAuth {
	if m == nil {
		return nil
	}

	return &domain.UserPasswordAuth{
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromDomain(u *domain.UserPasswordAuth) *UserPasswordNoSqlModel {
	if u == nil {
		return nil
	}

	return &UserPasswordNoSqlModel{
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

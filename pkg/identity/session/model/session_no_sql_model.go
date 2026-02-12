package model

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/session/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionNoSqlModel struct {
	ID                bson.ObjectID `bson:"_id"`
	UserID            bson.ObjectID `bson:"user_id"`
	RefreshTokenID    string        `bson:"refresh_token_id"`
	RefreshTokenExpIn int64         `bson:"refresh_token_expires_in"`
	Revoked           bool          `bson:"revoked"`
	RememberMe        bool          `bson:"remember_me"`
	CreatedAt         *time.Time    `bson:"created_at"`
}

func (m *SessionNoSqlModel) ToDomain(s *domain.Session) {
	if m == nil {
		return
	}

	s.ID = m.ID.Hex()
	s.UserID = m.UserID.Hex()
	s.RefreshTokenID = m.RefreshTokenID
	s.RefreshTokenExpIn = m.RefreshTokenExpIn
	s.Revoked = m.Revoked
	s.RememberMe = m.RememberMe
	s.CreatedAt = m.CreatedAt
}

func FromDomain(s *domain.Session, id bson.ObjectID, userID bson.ObjectID) *SessionNoSqlModel {
	if s == nil {
		return nil
	}

	return &SessionNoSqlModel{
		ID:                id,
		UserID:            userID,
		RefreshTokenID:    s.RefreshTokenID,
		RefreshTokenExpIn: s.RefreshTokenExpIn,
		Revoked:           s.Revoked,
		RememberMe:        s.RememberMe,
		CreatedAt:         s.CreatedAt,
	}
}

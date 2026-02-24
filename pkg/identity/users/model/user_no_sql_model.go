package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/model"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserNoSqlModel struct {
	ID                 bson.ObjectID                 `bson:"_id"`
	Email              string                        `bson:"email"`
	EmailVerified      bool                          `bson:"email_verified"`
	IsActive           bool                          `bson:"is_active"`
	ImgPath            *string                       `bson:"img_path"`
	UserPassMongoModel *model.UserPasswordNoSqlModel `bson:"pass_method"`
	RoleIDs            []bson.ObjectID               `bson:"role_ids"`
	CreatedAt          *time.Time                    `bson:"created_at"`
	UpdatedAt          *time.Time                    `bson:"updated_at"`
}

func (m *UserNoSqlModel) ToDomain(u *domain.User) {
	if m == nil {
		return
	}

	roleIDs := make([]string, 0, len(m.RoleIDs))
	for _, r := range m.RoleIDs {
		roleIDs = append(roleIDs, r.Hex())
	}

	u.ID = m.ID.Hex()
	u.Email = m.Email
	u.EmailVerified = m.EmailVerified
	u.IsActive = m.IsActive
	u.ImgPath = m.ImgPath
	u.UserPassAuth = m.UserPassMongoModel.ToDomain()
	u.RoleIDs = roleIDs
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
}

func FromDomain(u *domain.User, id bson.ObjectID) (*UserNoSqlModel, error) {
	if u == nil {
		return nil, nil
	}

	roleObjectIDs := make([]bson.ObjectID, 0, len(u.RoleIDs))

	for _, r := range u.RoleIDs {
		objID, err := bson.ObjectIDFromHex(r)
		if err != nil {
			return nil, err
		}
		roleObjectIDs = append(roleObjectIDs, objID)
	}

	return &UserNoSqlModel{
		ID:                 id,
		Email:              u.Email,
		EmailVerified:      u.EmailVerified,
		IsActive:           u.IsActive,
		ImgPath:            u.ImgPath,
		UserPassMongoModel: model.FromDomain(u.UserPassAuth),
		RoleIDs:            roleObjectIDs,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}, nil
}

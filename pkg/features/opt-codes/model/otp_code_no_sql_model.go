package model

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OtpCodeNoSqlModel struct {
	ID        bson.ObjectID         `bson:"_id"`
	UserID    bson.ObjectID         `bson:"user_id"`
	OptCode   string                `bson:"otp_code"`
	Purpose   domain.OtpCodePurpose `bson:"purpose"`
	Used      bool                  `bson:"used"`
	ExpiresAt *time.Time            `bson:"expires_at"`
	CreatedAt *time.Time            `bson:"created_at"`
}

func (m *OtpCodeNoSqlModel) ToDomain(o *domain.OtpCode) {
	if m == nil {
		return
	}

	o.ID = m.ID.Hex()
	o.UserID = m.UserID.Hex()
	o.OptCode = m.OptCode
	o.Purpose = m.Purpose
	o.Used = m.Used
	o.ExpiresAt = m.ExpiresAt
	o.CreatedAt = m.CreatedAt
}

func FromDomain(o *domain.OtpCode, id bson.ObjectID) (*OtpCodeNoSqlModel, error) {
	if o == nil {
		return nil, nil
	}

	userObjectID, err := bson.ObjectIDFromHex(o.UserID)
	if err != nil {
		return nil, err
	}

	return &OtpCodeNoSqlModel{
		ID:        id,
		UserID:    userObjectID,
		OptCode:   o.OptCode,
		Purpose:   o.Purpose,
		Used:      o.Used,
		ExpiresAt: o.ExpiresAt,
		CreatedAt: o.CreatedAt,
	}, nil
}

package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VarOptionNoSqlModel struct {
	ID          bson.ObjectID `bson:"_id"`
	VariationID bson.ObjectID `bson:"variation_id"`
	Label       string        `bson:"label"`
	Value       *string       `bson:"value"`
	CreatedAt   *time.Time    `bson:"created_at"`
	UpdatedAt   *time.Time    `bson:"updated_at"`
}

func (m *VarOptionNoSqlModel) ToDomain(o *domain.VarOption) {
	if m == nil {
		return
	}

	o.ID = m.ID.Hex()
	o.VariationID = m.VariationID.Hex()
	o.Label = m.Label
	o.Value = m.Value
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

func FromDomainOption(
	o *domain.VarOption,
	id bson.ObjectID,
) (*VarOptionNoSqlModel, error) {

	if o == nil {
		return nil, nil
	}

	variationID, err := bson.ObjectIDFromHex(o.VariationID)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	return &VarOptionNoSqlModel{
		ID:          id,
		VariationID: variationID,
		Label:       o.Label,
		Value:       o.Value,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}, nil
}

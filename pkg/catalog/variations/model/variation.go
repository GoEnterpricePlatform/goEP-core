package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VariationNoSqlModel struct {
	ID        bson.ObjectID          `bson:"_id"`
	Name      string                 `bson:"name"`
	Options   []*VarOptionNoSqlModel `bson:"options,omitempty"`
	CreatedAt *time.Time             `bson:"created_at"`
	UpdatedAt *time.Time             `bson:"updated_at"`
}

func (m *VariationNoSqlModel) ToDomain(v *domain.Variation) {
	if m == nil {
		return
	}

	options := make([]*domain.VarOption, 0, len(m.Options))

	for _, o := range m.Options {
		option := &domain.VarOption{}
		o.ToDomain(option)

		options = append(options, option)
	}

	v.ID = m.ID.Hex()
	v.Name = m.Name
	v.Options = options
	v.CreatedAt = m.CreatedAt
	v.UpdatedAt = m.UpdatedAt
}

func FromDomain(
	v *domain.Variation,
	id bson.ObjectID,
) (*VariationNoSqlModel, error) {

	if v == nil {
		return nil, nil
	}

	options := make([]*VarOptionNoSqlModel, 0, len(v.Options))

	for _, o := range v.Options {

		optionID, err := bson.ObjectIDFromHex(o.ID)
		if err != nil {
			return nil, sharedD.ErrIncorrectID
		}

		option, err := FromDomainOption(o,optionID,)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	return &VariationNoSqlModel{
		ID:        id,
		Name:      v.Name,
		Options:   options,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}, nil
}

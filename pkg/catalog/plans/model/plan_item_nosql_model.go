package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Options allows us to get and not save a planItem
type PlanItemNoSqlModel struct {
	ID           bson.ObjectID       `bson:"_id"`
	ImgPath      *string             `bson:"img_path"`
	Features     []string            `bson:"features"`
	Status       domain.PlanStatus   `bson:"status"`
	VarOptionIDs []bson.ObjectID     `bson:"var_option_ids"`
	Options      []*OptionNoSqlModel `bson:"options,omitempty"`
	CreatedAt    *time.Time          `bson:"created_at"`
	UpdatedAt    *time.Time          `bson:"updated_at"`
}

func (m *PlanItemNoSqlModel) ToDomain(o *domain.PlanItem) {
	if m == nil || o == nil {
		return
	}

	varOptionIDsStr := make([]string, 0, len(m.VarOptionIDs))

	for _, vOptionOID := range m.VarOptionIDs {
		varOptionIDsStr = append(varOptionIDsStr, vOptionOID.Hex())
	}

	o.VarOptionIDs = varOptionIDsStr

	// Convert the nested NoSQL options to domain options.
	options := make([]*domain.Option, 0, len(m.Options))

	for _, optionModel := range m.Options {
		if optionModel == nil {
			continue
		}

		option := &domain.Option{}
		optionModel.ToDomain(option)

		options = append(options, option)
	}

	o.Options = options

	o.ID = m.ID.Hex()
	o.ImgPath = m.ImgPath
	o.Features = m.Features
	o.Status = m.Status
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

func FromDomainPlanItem(o *domain.PlanItem, id bson.ObjectID) (*PlanItemNoSqlModel, error) {
	if o == nil {
		return nil, nil
	}

	varOptionIDs := make([]bson.ObjectID, 0, len(o.VarOptionIDs))

	for _, vOptionID := range o.VarOptionIDs {
		optionOID, err := bson.ObjectIDFromHex(vOptionID)
		if err != nil {
			return nil, sharedD.ErrIncorrectID
		}
		varOptionIDs = append(varOptionIDs, optionOID)
	}

	options := make([]*OptionNoSqlModel, 0, len(o.Options))

	for _, o := range o.Options {
		options = append(options, FromDomainOption(o))
	}

	return &PlanItemNoSqlModel{
		ID:           id,
		ImgPath:      o.ImgPath,
		Features:     o.Features,
		Status:       o.Status,
		VarOptionIDs: varOptionIDs,
		Options:      options,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}, nil
}

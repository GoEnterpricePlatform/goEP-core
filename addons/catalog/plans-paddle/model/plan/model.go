package model

import (
	planD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type PaddlePlanNoSqlModel struct {
	ID              bson.ObjectID          `bson:"_id"`
	PlanID          bson.ObjectID          `bson:"plan_id"`
	PaddleProductID string                 `bson:"paddle_product_id"`
	Order           int                    `bson:"order"`
	Items           []*PaddlePlanItemModel `bson:"items"`
	CreatedAt       *time.Time             `bson:"created_at"`
	UpdatedAt       *time.Time             `bson:"updated_at"`
}

type PaddlePlanItemModel struct {
	ID            bson.ObjectID    `bson:"_id"`
	PlanItemID    bson.ObjectID    `bson:"plan_item_id"`
	PaddlePriceID string           `bson:"paddle_price_id"`
	Status        planD.PlanStatus `bson:"status"`
	CreatedAt     *time.Time       `bson:"created_at"`
	UpdatedAt     *time.Time       `bson:"updated_at"`
}

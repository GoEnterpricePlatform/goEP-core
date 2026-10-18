package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const processingLease = 5 * time.Minute

var _ port.PaddleWebhookEventRepo = &WebhookEventRepository{}

type WebhookEventRepository struct {
	Collection *mongoDriver.Collection
}


func NewPaddleWebhookEventRepo(collection *mongoDriver.Collection) *WebhookEventRepository {
	return &WebhookEventRepository{Collection: collection}
}

package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const processingLease = 5 * time.Minute

var _ port.PaddleWebhookEventRepo = &WebhookEventRepository{}

type WebhookEventRepository struct {
	Collection *mongoDriver.Collection
}

type webhookEventModel struct {
	ID                  string     `bson:"_id"`
	EventID             string     `bson:"event_id"`
	EventType           string     `bson:"event_type"`
	ResourceID          string     `bson:"resource_id,omitempty"`
	OccurredAt          string     `bson:"occurred_at,omitempty"`
	Status              string     `bson:"status"`
	ReceivedAt          time.Time  `bson:"received_at"`
	LastReceivedAt      time.Time  `bson:"last_received_at"`
	ProcessingStartedAt time.Time  `bson:"processing_started_at"`
	ProcessedAt         *time.Time `bson:"processed_at,omitempty"`
	FailedAt            *time.Time `bson:"failed_at,omitempty"`
	Attempts            int        `bson:"attempts"`
	LastError           string     `bson:"last_error,omitempty"`
}

func NewPaddleWebhookEventRepo(collection *mongoDriver.Collection) *WebhookEventRepository {
	return &WebhookEventRepository{Collection: collection}
}

// StartProcessing claims a new event, retries a failed/stale event, or reports
// that an already processed event can safely be acknowledged without replaying it.
func (r *WebhookEventRepository) StartProcessing(ctx context.Context, event *domain.PaddleWebhookEvent, now time.Time) (bool, error) {
	if event == nil || event.EventID == "" {
		return false, fmt.Errorf("Paddle webhook event is missing its event ID")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	record := webhookEventModel{
		ID:                  bson.NewObjectID().Hex(),
		EventID:             event.EventID,
		EventType:           event.EventType,
		ResourceID:          event.Data.ID,
		OccurredAt:          event.OccurredAt,
		Status:              "processing",
		ReceivedAt:          now,
		LastReceivedAt:      now,
		ProcessingStartedAt: now,
		Attempts:            1,
	}
	if _, err := r.Collection.InsertOne(ctx, record); err == nil {
		return false, nil
	} else if !mongoDriver.IsDuplicateKeyError(err) {
		return false, fmt.Errorf("error inserting Paddle webhook event: %w", err)
	}

	// A duplicate delivery may retry an earlier failure. Also reclaim processing
	// records left behind by a crashed worker after the lease expires.
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{
			"event_id": event.EventID,
			"$or": bson.A{
				bson.M{"status": "failed"},
				bson.M{
					"status":                "processing",
					"processing_started_at": bson.M{"$lte": now.Add(-processingLease)},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"event_type":            event.EventType,
				"resource_id":           event.Data.ID,
				"occurred_at":           event.OccurredAt,
				"status":                "processing",
				"last_received_at":      now,
				"processing_started_at": now,
			},
			"$inc":   bson.M{"attempts": 1},
			"$unset": bson.M{"last_error": "", "processed_at": "", "failed_at": ""},
		},
	)
	if err != nil {
		return false, fmt.Errorf("error claiming Paddle webhook event retry: %w", err)
	}
	if result.ModifiedCount > 0 {
		return false, nil
	}

	var existing webhookEventModel
	if err := r.Collection.FindOne(ctx, bson.M{"event_id": event.EventID}).Decode(&existing); err != nil {
		return false, fmt.Errorf("error checking duplicate Paddle webhook event: %w", err)
	}
	if existing.Status == "processed" {
		return true, nil
	}
	return false, fmt.Errorf("Paddle webhook event %s is already being processed", event.EventID)
}

func (r *WebhookEventRepository) MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error {
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{"event_id": eventID, "status": "processing"},
		bson.M{
			"$set":   bson.M{"status": "processed", "processed_at": processedAt},
			"$unset": bson.M{"last_error": "", "failed_at": ""},
		},
	)
	if err != nil {
		return fmt.Errorf("error marking Paddle webhook event as processed: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}

func (r *WebhookEventRepository) MarkFailed(ctx context.Context, eventID string, failedAt time.Time, processingErr error) error {
	errorMessage := "unknown processing error"
	if processingErr != nil {
		errorMessage = strings.TrimSpace(processingErr.Error())
	}
	if len(errorMessage) > 2000 {
		errorMessage = errorMessage[:2000]
	}
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{"event_id": eventID, "status": "processing"},
		bson.M{"$set": bson.M{
			"status":     "failed",
			"last_error": errorMessage,
			"failed_at":  failedAt,
		}},
	)
	if err != nil {
		return fmt.Errorf("error marking Paddle webhook event as failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}

func (r *WebhookEventRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongoDriver.IndexModel{
		Keys:    bson.D{{Key: "event_id", Value: 1}},
		Options: options.Index().SetName("paddle_webhook_event_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle webhook event ID index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongoDriver.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}, {Key: "received_at", Value: -1}},
		Options: options.Index().SetName("paddle_webhook_event_status_received_at"),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle webhook event audit index: %w", err)
	}
	return nil
}

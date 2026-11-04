package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const processingLease = 5 * time.Minute

var _ port.PaddleWebhookEventRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewPaddleWebhookEventRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

func (r *Repository) CreateIndexes() error {
	ctx := context.Background()
	if _, err := r.Pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_paddle_webhook_event_status_received_at ON tb_paddle_webhook_event(status,received_at DESC)`); err != nil {
		return fmt.Errorf("create Paddle webhook event index: %w", err)
	}
	return nil
}

func (r *Repository) StartProcessing(ctx context.Context, event *domain.PaddleWebhookEvent, now time.Time) (bool, error) {
	if event == nil || event.EventID == "" {
		return false, fmt.Errorf("Paddle webhook event is missing its event ID")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	var occurredAt any
	if event.OccurredAt != "" {
		parsed, err := time.Parse(time.RFC3339Nano, event.OccurredAt)
		if err != nil {
			return false, fmt.Errorf("invalid Paddle webhook occurred_at: %w", err)
		}
		occurredAt = parsed
	}
	var recordID string
	err := r.Pool.QueryRow(ctx, `INSERT INTO tb_paddle_webhook_event(id,event_id,event_type,resource_id,occurred_at,status,received_at,last_received_at,processing_started_at,attempts) VALUES($1,$2,$3,NULLIF($4,''),$5,'processing',$6,$6,$6,1)
ON CONFLICT(event_id) DO UPDATE SET event_type=EXCLUDED.event_type,resource_id=EXCLUDED.resource_id,occurred_at=EXCLUDED.occurred_at,status='processing',last_received_at=EXCLUDED.last_received_at,processing_started_at=EXCLUDED.processing_started_at,attempts=tb_paddle_webhook_event.attempts+1,last_error=NULL,processed_at=NULL,failed_at=NULL
WHERE tb_paddle_webhook_event.status='failed' OR (tb_paddle_webhook_event.status='processing' AND tb_paddle_webhook_event.processing_started_at <= $6 - ($7 * INTERVAL '1 second')) RETURNING id::text`, uuid.New(), event.EventID, event.EventType, event.Data.ID, occurredAt, now, processingLease.Seconds()).Scan(&recordID)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("claim Paddle webhook event: %w", err)
	}
	var status string
	err = r.Pool.QueryRow(ctx, `SELECT status FROM tb_paddle_webhook_event WHERE event_id=$1`, event.EventID).Scan(&status)
	if err != nil {
		return false, fmt.Errorf("read duplicate Paddle webhook event: %w", err)
	}
	if status == "processed" {
		return true, nil
	}
	return false, fmt.Errorf("Paddle webhook event %s is already being processed", event.EventID)
}

func (r *Repository) MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error {
	result, err := r.Pool.Exec(ctx, `UPDATE tb_paddle_webhook_event SET status='processed',processed_at=$2,last_error=NULL,failed_at=NULL WHERE event_id=$1 AND status='processing'`, eventID, processedAt)
	if err != nil {
		return fmt.Errorf("mark Paddle webhook event processed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, eventID string, failedAt time.Time, processingErr error) error {
	message := "unknown processing error"
	if processingErr != nil {
		message = strings.TrimSpace(processingErr.Error())
	}
	if len(message) > 2000 {
		message = message[:2000]
	}
	result, err := r.Pool.Exec(ctx, `UPDATE tb_paddle_webhook_event SET status='failed',last_error=$3,failed_at=$2 WHERE event_id=$1 AND status='processing'`, eventID, failedAt, message)
	if err != nil {
		return fmt.Errorf("mark Paddle webhook event failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}

package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.PaddleSubscriptionRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewPaddleSubscriptionRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

func parseUUID(raw, kind string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s ID: %w", shared.ErrIncorrectID, kind, err)
	}
	return id, nil
}

func optionalUUID(raw, kind string) (any, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := parseUUID(raw, kind)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func parseTime(raw *string) (any, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	v, err := time.Parse(time.RFC3339Nano, *raw)
	if err != nil {
		return nil, fmt.Errorf("invalid Paddle timestamp %q: %w", *raw, err)
	}
	return v, nil
}

type rowScanner interface{ Scan(...any) error }

func scanSubscription(row rowScanner) (*domain.PaddleSubscription, error) {
	var s domain.PaddleSubscription
	var customerID, planID, planItemID *string
	var periodStart, periodEnd, nextBilled, canceledAt *time.Time
	var scheduledAction *string
	var scheduledEffectiveAt *time.Time
	var updatedAt time.Time
	err := row.Scan(&s.ID, &s.PaddleSubscriptionID, &s.UserID, &customerID, &planID, &planItemID, &s.Status, &periodStart, &periodEnd, &nextBilled, &canceledAt, &scheduledAction, &scheduledEffectiveAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if customerID != nil {
		s.CustomerID = *customerID
	}
	if planID != nil {
		s.PlanID = *planID
	}
	if planItemID != nil {
		s.PlanItemID = *planItemID
	}
	format := func(v *time.Time) *string {
		if v == nil {
			return nil
		}
		text := v.UTC().Format(time.RFC3339Nano)
		return &text
	}
	s.CurrentPeriodStartsAt = format(periodStart)
	s.CurrentPeriodEndsAt = format(periodEnd)
	s.NextBilledAt = format(nextBilled)
	s.CanceledAt = format(canceledAt)
	if scheduledAction != nil && scheduledEffectiveAt != nil {
		s.ScheduledChange = &domain.PaddleSubscriptionChange{Action: *scheduledAction, EffectiveAt: scheduledEffectiveAt.UTC().Format(time.RFC3339Nano)}
	}
	s.UpdatedAt = updatedAt.UTC().Format(time.RFC3339Nano)
	return &s, nil
}

const subscriptionSelect = `SELECT id::text,paddle_subscription_id,user_id::text,customer_id::text,plan_id::text,plan_item_id::text,status,current_period_starts_at,current_period_ends_at,next_billed_at,canceled_at,scheduled_change_action,scheduled_change_effective_at,updated_at FROM tb_paddle_subscription`

func (r *Repository) Upsert(ctx context.Context, s *domain.PaddleSubscription, lastEventAt time.Time) error {
	if s == nil || s.PaddleSubscriptionID == "" {
		return fmt.Errorf("Paddle subscription is missing Paddle subscription ID")
	}
	userID, err := parseUUID(s.UserID, "user")
	if err != nil {
		return err
	}
	id := uuid.New()
	if s.ID != "" {
		id, err = parseUUID(s.ID, "subscription")
		if err != nil {
			return err
		}
	}
	customerID, err := optionalUUID(s.CustomerID, "Paddle customer")
	if err != nil {
		return err
	}
	planID, err := optionalUUID(s.PlanID, "plan")
	if err != nil {
		return err
	}
	planItemID, err := optionalUUID(s.PlanItemID, "plan item")
	if err != nil {
		return err
	}
	periodStart, err := parseTime(s.CurrentPeriodStartsAt)
	if err != nil {
		return err
	}
	periodEnd, err := parseTime(s.CurrentPeriodEndsAt)
	if err != nil {
		return err
	}
	nextBilled, err := parseTime(s.NextBilledAt)
	if err != nil {
		return err
	}
	canceledAt, err := parseTime(s.CanceledAt)
	if err != nil {
		return err
	}
	var changeAction, changeEffective any
	if s.ScheduledChange != nil {
		changeAction = s.ScheduledChange.Action
		parsed, e := time.Parse(time.RFC3339Nano, s.ScheduledChange.EffectiveAt)
		if e != nil {
			return fmt.Errorf("invalid Paddle scheduled change timestamp: %w", e)
		}
		changeEffective = parsed
	}
	updatedAt := time.Now().UTC()
	if s.UpdatedAt != "" {
		updatedAt, err = time.Parse(time.RFC3339Nano, s.UpdatedAt)
		if err != nil {
			return fmt.Errorf("invalid Paddle subscription updated_at: %w", err)
		}
	}
	if lastEventAt.IsZero() {
		lastEventAt = updatedAt
	}
	if s.ID == "" { // Preserve the stable local identity when replaying an existing Paddle subscription.
		_ = r.Pool.QueryRow(ctx, `SELECT id FROM tb_paddle_subscription WHERE paddle_subscription_id=$1`, s.PaddleSubscriptionID).Scan(&id)
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_paddle_subscription(id,paddle_subscription_id,user_id,customer_id,plan_id,plan_item_id,status,current_period_starts_at,current_period_ends_at,next_billed_at,canceled_at,scheduled_change_action,scheduled_change_effective_at,updated_at,last_event_at)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
ON CONFLICT(paddle_subscription_id) DO UPDATE SET user_id=EXCLUDED.user_id,customer_id=EXCLUDED.customer_id,plan_id=EXCLUDED.plan_id,plan_item_id=EXCLUDED.plan_item_id,status=EXCLUDED.status,current_period_starts_at=EXCLUDED.current_period_starts_at,current_period_ends_at=EXCLUDED.current_period_ends_at,next_billed_at=EXCLUDED.next_billed_at,canceled_at=EXCLUDED.canceled_at,scheduled_change_action=EXCLUDED.scheduled_change_action,scheduled_change_effective_at=EXCLUDED.scheduled_change_effective_at,updated_at=EXCLUDED.updated_at,last_event_at=EXCLUDED.last_event_at
WHERE tb_paddle_subscription.last_event_at IS NULL OR tb_paddle_subscription.last_event_at <= EXCLUDED.last_event_at`, id, s.PaddleSubscriptionID, userID, customerID, planID, planItemID, s.Status, periodStart, periodEnd, nextBilled, canceledAt, changeAction, changeEffective, updatedAt, lastEventAt)
	if err != nil {
		return fmt.Errorf("upsert Paddle subscription: %w", err)
	}
	s.ID = id.String()
	return nil
}

func (r *Repository) FindByID(ctx context.Context, rawID string) (*domain.PaddleSubscription, error) {
	id, err := parseUUID(rawID, "subscription")
	if err != nil {
		return nil, err
	}
	return r.findOne(ctx, subscriptionSelect+` WHERE id=$1`, id)
}
func (r *Repository) FindByPaddleSubscriptionID(ctx context.Context, id string) (*domain.PaddleSubscription, error) {
	return r.findOne(ctx, subscriptionSelect+` WHERE paddle_subscription_id=$1`, id)
}
func (r *Repository) FindCurrentByUserID(ctx context.Context, rawID string) (*domain.PaddleSubscription, error) {
	id, err := parseUUID(rawID, "user")
	if err != nil {
		return nil, err
	}
	return r.findOne(ctx, subscriptionSelect+` WHERE user_id=$1 AND status IN ('active','trialing') ORDER BY updated_at DESC,id DESC LIMIT 1`, id)
}

func (r *Repository) findOne(ctx context.Context, query string, arg any) (*domain.PaddleSubscription, error) {
	s, err := scanSubscription(r.Pool.QueryRow(ctx, query, arg))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find Paddle subscription: %w", err)
	}
	return s, nil
}

func (r *Repository) findMany(ctx context.Context, query string, args ...any) ([]*domain.PaddleSubscription, error) {
	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list Paddle subscriptions: %w", err)
	}
	defer rows.Close()
	result := make([]*domain.PaddleSubscription, 0)
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("scan Paddle subscription: %w", err)
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Paddle subscriptions: %w", err)
	}
	return result, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, rawID string) ([]*domain.PaddleSubscription, error) {
	id, err := parseUUID(rawID, "user")
	if err != nil {
		return nil, err
	}
	return r.findMany(ctx, subscriptionSelect+` WHERE user_id=$1 ORDER BY updated_at DESC,id DESC`, id)
}
func (r *Repository) FindAll(ctx context.Context) ([]*domain.PaddleSubscription, error) {
	return r.findMany(ctx, subscriptionSelect+` ORDER BY updated_at DESC,id DESC`)
}

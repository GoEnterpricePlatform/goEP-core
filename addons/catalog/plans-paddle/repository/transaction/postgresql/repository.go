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

var _ port.PaddleTransactionRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewPaddleTransactionRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

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

type rowScanner interface{ Scan(...any) error }

func scanTransaction(row rowScanner) (*domain.PaddleTransaction, error) {
	var t domain.PaddleTransaction
	var userID, subscriptionID, customerID *string
	var currency *string
	var lastEventAt *time.Time
	var updatedAt time.Time
	err := row.Scan(&t.ID, &t.PaddleTransactionID, &userID, &subscriptionID, &customerID, &t.PlanID, &t.PlanItemID, &t.Status, &currency, &t.Subtotal, &t.Tax, &t.Total, &updatedAt, &lastEventAt, &t.PaymentDetailsProvided)
	if err != nil {
		return nil, err
	}
	if userID != nil {
		t.UserID = *userID
	}
	if subscriptionID != nil {
		t.SubscriptionID = subscriptionID
	}
	if customerID != nil {
		t.CustomerID = *customerID
	}
	if currency != nil {
		t.CurrencyCode = *currency
	}
	t.UpdatedAt = updatedAt.UTC().Format(time.RFC3339Nano)
	if lastEventAt != nil {
		t.LastEventAt = *lastEventAt
	}
	return &t, nil
}

const transactionSelect = `SELECT id::text,paddle_transaction_id,user_id::text,subscription_id::text,customer_id::text,plan_id::text,plan_item_id::text,status,currency_code,subtotal,tax,total,updated_at,last_event_at,payment_details_provided FROM tb_paddle_transaction`

func (r *Repository) Insert(ctx context.Context, t *domain.PaddleTransaction) error {
	if t == nil || t.PaddleTransactionID == "" {
		return fmt.Errorf("Paddle transaction is missing Paddle transaction ID")
	}
	id := uuid.New()
	var err error
	if t.ID != "" {
		id, err = parseUUID(t.ID, "transaction")
		if err != nil {
			return err
		}
	}
	userID, err := optionalUUID(t.UserID, "user")
	if err != nil {
		return err
	}
	subscriptionID, err := optionalUUID(stringValue(t.SubscriptionID), "subscription")
	if err != nil {
		return err
	}
	customerID, err := optionalUUID(t.CustomerID, "Paddle customer")
	if err != nil {
		return err
	}
	planID, err := parseUUID(t.PlanID, "plan")
	if err != nil {
		return err
	}
	planItemID, err := parseUUID(t.PlanItemID, "plan item")
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("invalid Paddle transaction updated_at: %w", err)
	}
	var lastEventAt any
	if !t.LastEventAt.IsZero() {
		lastEventAt = t.LastEventAt
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_paddle_transaction(id,paddle_transaction_id,user_id,subscription_id,customer_id,plan_id,plan_item_id,status,currency_code,subtotal,tax,total,updated_at,last_event_at,payment_details_provided) VALUES($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9,''),$10,$11,$12,$13,$14,$15)`, id, t.PaddleTransactionID, userID, subscriptionID, customerID, planID, planItemID, t.Status, t.CurrencyCode, t.Subtotal, t.Tax, t.Total, updatedAt, lastEventAt, t.PaymentDetailsProvided)
	if err != nil {
		return fmt.Errorf("insert Paddle transaction: %w", err)
	}
	t.ID = id.String()
	return nil
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (r *Repository) FindByPaddleTransactionID(ctx context.Context, rawID string) (*domain.PaddleTransaction, error) {
	t, err := scanTransaction(r.Pool.QueryRow(ctx, transactionSelect+` WHERE paddle_transaction_id=$1`, rawID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find Paddle transaction: %w", err)
	}
	return t, nil
}

func (r *Repository) findMany(ctx context.Context, query string, args ...any) ([]*domain.PaddleTransaction, error) {
	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list Paddle transactions: %w", err)
	}
	defer rows.Close()
	result := make([]*domain.PaddleTransaction, 0)
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("scan Paddle transaction: %w", err)
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Paddle transactions: %w", err)
	}
	return result, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, rawID string) ([]*domain.PaddleTransaction, error) {
	id, err := parseUUID(rawID, "user")
	if err != nil {
		return nil, err
	}
	return r.findMany(ctx, transactionSelect+` WHERE user_id=$1 ORDER BY updated_at DESC`, id)
}
func (r *Repository) FindAllLinkedToSubscriptions(ctx context.Context) ([]*domain.PaddleTransaction, error) {
	return r.findMany(ctx, transactionSelect+` WHERE subscription_id IS NOT NULL ORDER BY updated_at DESC`)
}

func (r *Repository) UpsertFromWebhook(ctx context.Context, t *domain.PaddleTransaction, eventAt time.Time) error {
	if t == nil || t.PaddleTransactionID == "" {
		return fmt.Errorf("Paddle webhook transaction is missing Paddle transaction ID")
	}
	if eventAt.IsZero() {
		eventAt = time.Now().UTC()
	}
	id := uuid.New()
	var err error
	if t.ID != "" {
		id, err = parseUUID(t.ID, "transaction")
		if err != nil {
			return err
		}
	}
	userID, err := optionalUUID(t.UserID, "user")
	if err != nil {
		return err
	}
	subscriptionID, err := optionalUUID(stringValue(t.SubscriptionID), "subscription")
	if err != nil {
		return err
	}
	customerID, err := optionalUUID(t.CustomerID, "Paddle customer")
	if err != nil {
		return err
	}
	planID, err := parseUUID(t.PlanID, "plan")
	if err != nil {
		return err
	}
	planItemID, err := parseUUID(t.PlanItemID, "plan item")
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("invalid Paddle transaction updated_at: %w", err)
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_paddle_transaction(id,paddle_transaction_id,user_id,subscription_id,customer_id,plan_id,plan_item_id,status,currency_code,subtotal,tax,total,updated_at,last_event_at,payment_details_provided) VALUES($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9,''),$10,$11,$12,$13,$14,$15)
ON CONFLICT(paddle_transaction_id) DO UPDATE SET user_id=COALESCE(EXCLUDED.user_id,tb_paddle_transaction.user_id),subscription_id=COALESCE(EXCLUDED.subscription_id,tb_paddle_transaction.subscription_id),customer_id=COALESCE(EXCLUDED.customer_id,tb_paddle_transaction.customer_id),plan_id=EXCLUDED.plan_id,plan_item_id=EXCLUDED.plan_item_id,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at,last_event_at=EXCLUDED.last_event_at,currency_code=CASE WHEN tb_paddle_transaction.payment_details_provided THEN tb_paddle_transaction.currency_code ELSE EXCLUDED.currency_code END,subtotal=CASE WHEN tb_paddle_transaction.payment_details_provided THEN tb_paddle_transaction.subtotal ELSE EXCLUDED.subtotal END,tax=CASE WHEN tb_paddle_transaction.payment_details_provided THEN tb_paddle_transaction.tax ELSE EXCLUDED.tax END,total=CASE WHEN tb_paddle_transaction.payment_details_provided THEN tb_paddle_transaction.total ELSE EXCLUDED.total END,payment_details_provided=tb_paddle_transaction.payment_details_provided OR EXCLUDED.payment_details_provided
WHERE tb_paddle_transaction.last_event_at IS NULL OR tb_paddle_transaction.last_event_at <= EXCLUDED.last_event_at`, id, t.PaddleTransactionID, userID, subscriptionID, customerID, planID, planItemID, t.Status, t.CurrencyCode, t.Subtotal, t.Tax, t.Total, updatedAt, eventAt, t.PaymentDetailsProvided)
	if err != nil {
		return fmt.Errorf("upsert Paddle transaction from webhook: %w", err)
	}
	t.LastEventAt = eventAt
	return nil
}

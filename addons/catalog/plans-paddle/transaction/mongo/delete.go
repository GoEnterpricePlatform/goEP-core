package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (t *Transaction) Delete(ctx context.Context, id string) error {
	session, err := t.Client.StartSession()
	if err != nil {
		return fmt.Errorf("error starting session: %w", err)
	}
	defer session.EndSession(ctx)
	if err := session.StartTransaction(); err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		if err := t.deleteLinked(ctx, id); err != nil {
			return err
		}
		return session.CommitTransaction(ctx)
	})
	if err != nil {
		_ = session.AbortTransaction(context.Background())
		return fmt.Errorf("error executing delete transaction: %w", err)
	}
	return nil
}

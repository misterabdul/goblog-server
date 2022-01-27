package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionHandlerFunc func(ctx context.Context, dbConn *mongo.Database) error

func Transaction(
	ctx context.Context,
	dbConn *mongo.Database,
	canFallback bool,
	handler TransactionHandlerFunc,
) (err error) {
	var session mongo.Session

	if session, err = dbConn.Client().StartSession(); err != nil {
		if canFallback {
			return handler(ctx, dbConn)
		}
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
		if sErr = sCtx.StartTransaction(); sErr != nil {
			return sErr
		}
		if sErr = handler(sCtx, dbConn); sErr != nil {
			return sErr
		}
		if sErr = session.CommitTransaction(sCtx); sErr != nil {
			return sErr
		}

		return nil
	})
}

package db

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Format-C-eft/utils/logger"
)

func TransactionWrapper(
	ctx context.Context,
	conn Client,
	wrappedFunc func(ctx context.Context, conn Client) error,
) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("conn.Begin err: %w", err)
	}

	err = wrapForRecoverPanic(ctx, tx, wrappedFunc)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("wrapForRecoverPanic err: %w and rollback error: %w", rollbackErr, err)
		}

		return fmt.Errorf("wrapForRecoverPanic err: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit err: %w", err)
	}

	return nil
}

func wrapForRecoverPanic(
	ctx context.Context,
	tx pgx.Tx,
	wrappedFunc func(ctx context.Context, conn Client) error,
) (err error) {
	defer Recover(ctx, func(recoveryErr error) {
		err = recoveryErr
	})

	return wrappedFunc(ctx, &connAdapterTx{client: tx})
}

func Recover(ctx context.Context, fn func(err error)) {
	if p := recover(); p != nil {
		err := fmt.Errorf("recovered from panic: %v", p)
		logger.Error(ctx, err)

		stack := strings.Split(string(debug.Stack()), "\n")
		logger.ErrorKV(ctx, fmt.Sprintf("%v", p), "stack_trace", stack)

		if fn != nil {
			fn(err)
		}
	}
}

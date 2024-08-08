package qbun

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
)

type txFunc func(ctx context.Context, tx bun.Tx) error

func RunInTx(ctx context.Context, fs ...txFunc) error {
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, f := range fs {
			err := f(ctx, tx)

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

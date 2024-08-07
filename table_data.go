package qbun

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"time"
)

type ModelSelector interface {
	Select(ctx context.Context, tx bun.Tx) error
}

type TableData struct {
	Column string
	Key    interface{}
	Value  interface{}
}

type TableDataArray []*TableData

type UID struct {
	Update TableDataArray
	Insert TableDataArray
	Delete []string
}

// UpdateCell set value of TableData.Column to TableData.Value where value of column TableData.Key is equal to pkColumn.
func (td *TableData) UpdateCell(ctx context.Context, db *bun.DB, model interface{}, pkColumn string, returnedModel ModelSelector) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error BeginTx in UpdateCell function")
		return err
	}

	_, err = tx.NewUpdate().
		Model(model).
		Set(fmt.Sprintf(`"%s" = ?`, td.Column), td.Value).
		Set("updated_at = ?", time.Now()).
		Where(fmt.Sprintf(`"%s" = ?`, pkColumn), td.Key).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("UpdateCell SQL error: update table failed")
		e := tx.Rollback()
		if e != nil {
			log.Error().Err(e).Msg("Error Rollback in UpdateCell function")
			return e
		}
		return err
	}

	err = returnedModel.Select(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Error Commit in UpdateCell function")
		return err
	}

	return nil
}

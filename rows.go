package qbun

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
)

func InsertRowsIfNotExist(ctx context.Context, tx bun.Tx, model interface{}, columns []string) error {
	_, err := tx.NewInsert().
		Model(model).
		Column(columns...).
		On(fmt.Sprintf("CONFLICT (%s) DO UPDATE", strings.Join(columns, ", "))).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUnusedRows remove rows from table (model) if any value of this table's column is not contained in other table's (fkTable) column (fkColumn).
func DeleteUnusedRows(ctx context.Context, tx bun.Tx, model interface{}, fkTable, fkColumn, column string) error {
	subqTs := tx.NewSelect().
		Model(model).
		Column("id").
		Join(fmt.Sprintf("LEFT JOIN %s ON %s = %s", fkTable, fkColumn, column)).
		Where(fmt.Sprintf("%s IS NULL", fkColumn))

	_, err := tx.NewDelete().
		With("_data", subqTs).
		Model(model).
		Table("_data").
		Where(fmt.Sprintf("%s = _data.id", column)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

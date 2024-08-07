package qbun

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

// CheckColumnType returns the type of the column.
func CheckColumnType(ctx context.Context, tx bun.Tx, model interface{}, column string) (string, error) {
	var columnType string
	err := tx.NewSelect().
		Model(model).
		ColumnExpr(fmt.Sprintf(`pg_typeof(%s)`, column)).
		Limit(1).
		Scan(ctx, &columnType)
	if err != nil {
		return "", err
	}
	return columnType, nil
}

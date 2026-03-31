package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	purchasesTableName = "purchases"
)

func init() {
	registerTableName(purchasesTableName)
}

var purchasesColumns = []string{
	idColumn,
	belongsToAccountColumn,
	"product_id",
	"amount_cents",
	"currency",
	"completed_at",
	"external_transaction_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildPaymentsPurchasesQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(purchasesColumns)
		fullSelectColumns := applyToEach(purchasesColumns, func(_ int, s string) string {
			return fullColumnName(purchasesTableName, s)
		})
		accountCondition := fmt.Sprintf("%s.%s = sqlc.arg(%s)", purchasesTableName, belongsToAccountColumn, belongsToAccountColumn)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreatePurchase",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					purchasesTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPurchase",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					purchasesTableName,
					purchasesTableName, archivedAtColumn,
					purchasesTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPurchasesForAccount",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(purchasesTableName, true, true, nil, accountCondition),
					buildTotalCountSelect(purchasesTableName, true, nil, accountCondition),
					purchasesTableName,
					purchasesTableName, archivedAtColumn,
					accountCondition,
					buildFilterConditions(purchasesTableName, true, false, accountCondition),
					buildCursorLimitClause(purchasesTableName),
				)),
			},
		}
	default:
		return nil
	}
}

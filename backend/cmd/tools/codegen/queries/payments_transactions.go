package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	paymentTransactionsTableName = "payment_transactions"
)

func init() {
	registerTableName(paymentTransactionsTableName)
}

var paymentTransactionsColumns = []string{
	idColumn,
	belongsToAccountColumn,
	"subscription_id",
	"purchase_id",
	"external_transaction_id",
	"amount_cents",
	"currency",
	"status",
	createdAtColumn,
}

func buildPaymentsTransactionsQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(paymentTransactionsColumns)
		fullSelectColumns := applyToEach(paymentTransactionsColumns, func(_ int, s string) string {
			return fullColumnName(paymentTransactionsTableName, s)
		})
		accountCondition := fmt.Sprintf("%s.%s = sqlc.arg(%s)", paymentTransactionsTableName, belongsToAccountColumn, belongsToAccountColumn)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreatePaymentTransaction",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					paymentTransactionsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPaymentTransactionsForAccount",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(paymentTransactionsTableName, false, false, nil, accountCondition),
					buildTotalCountSelect(paymentTransactionsTableName, false, nil, accountCondition),
					paymentTransactionsTableName,
					accountCondition,
					buildFilterConditions(paymentTransactionsTableName, false, false, accountCondition),
					buildCursorLimitClause(paymentTransactionsTableName),
				)),
			},
		}
	default:
		return nil
	}
}

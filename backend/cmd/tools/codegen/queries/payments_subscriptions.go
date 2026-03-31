package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	subscriptionsTableName       = "subscriptions"
	externalSubscriptionIDColumn = "external_subscription_id"
)

func init() {
	registerTableName(subscriptionsTableName)
}

var subscriptionsColumns = []string{
	idColumn,
	belongsToAccountColumn,
	"product_id",
	externalSubscriptionIDColumn,
	"status",
	"current_period_start",
	"current_period_end",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildPaymentsSubscriptionsQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(subscriptionsColumns)
		fullSelectColumns := applyToEach(subscriptionsColumns, func(_ int, s string) string {
			return fullColumnName(subscriptionsTableName, s)
		})
		accountCondition := fmt.Sprintf("%s.%s = sqlc.arg(%s)", subscriptionsTableName, belongsToAccountColumn, belongsToAccountColumn)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveSubscription",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s, %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					subscriptionsTableName,
					archivedAtColumn, currentTimeExpression,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateSubscription",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					subscriptionsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetSubscription",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					subscriptionsTableName,
					subscriptionsTableName, archivedAtColumn,
					subscriptionsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetSubscriptionByExternalID",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					subscriptionsTableName,
					subscriptionsTableName, archivedAtColumn,
					subscriptionsTableName, externalSubscriptionIDColumn, externalSubscriptionIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetSubscriptionsForAccount",
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
					buildFilterCountSelect(subscriptionsTableName, true, true, nil, accountCondition),
					buildTotalCountSelect(subscriptionsTableName, true, nil, accountCondition),
					subscriptionsTableName,
					subscriptionsTableName, archivedAtColumn,
					accountCondition,
					buildFilterConditions(subscriptionsTableName, true, false, accountCondition),
					buildCursorLimitClause(subscriptionsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateSubscription",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					subscriptionsTableName,
					strings.Join(applyToEach(filterForUpdate(subscriptionsColumns, belongsToAccountColumn, "product_id"), func(_ int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateSubscriptionStatus",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET status = sqlc.arg(status), %s = %s
WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					subscriptionsTableName,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}

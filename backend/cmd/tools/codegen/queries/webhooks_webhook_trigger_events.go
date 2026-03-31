package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	webhookTriggerEventsTableName = "webhook_trigger_events"
)

func init() {
	registerTableName(webhookTriggerEventsTableName)
}

var (
	webhookTriggerEventsColumns = []string{
		idColumn,
		nameColumn,
		descriptionColumn,
		createdAtColumn,
		lastUpdatedAtColumn,
		archivedAtColumn,
	}
)

func buildWebhookTriggerEventsQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(webhookTriggerEventsColumns)
		fullSelectColumns := applyToEach(webhookTriggerEventsColumns, func(_ int, s string) string {
			return fullColumnName(webhookTriggerEventsTableName, s)
		})

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreateWebhookTriggerEvent",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					webhookTriggerEventsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetWebhookTriggerEvent",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					webhookTriggerEventsTableName,
					webhookTriggerEventsTableName, archivedAtColumn,
					webhookTriggerEventsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckWebhookTriggerEventExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					webhookTriggerEventsTableName, idColumn,
					webhookTriggerEventsTableName,
					webhookTriggerEventsTableName, archivedAtColumn,
					webhookTriggerEventsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetWebhookTriggerEvents",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(webhookTriggerEventsTableName, true, true, nil),
					buildTotalCountSelect(webhookTriggerEventsTableName, true, nil),
					webhookTriggerEventsTableName,
					webhookTriggerEventsTableName, archivedAtColumn,
					buildFilterConditions(webhookTriggerEventsTableName, true, true),
					buildCursorLimitClause(webhookTriggerEventsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateWebhookTriggerEvent",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					webhookTriggerEventsTableName,
					strings.Join(applyToEach(filterForUpdate(webhookTriggerEventsColumns), func(_ int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveWebhookTriggerEvent",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					webhookTriggerEventsTableName,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}

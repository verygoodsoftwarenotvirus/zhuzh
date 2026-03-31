package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	webhookTriggerConfigsTableName = "webhook_trigger_configs"
	belongsToWebhookColumn         = "belongs_to_webhook"
	triggerEventColumn             = "trigger_event"
)

func init() {
	registerTableName(webhookTriggerConfigsTableName)
}

var (
	webhookTriggerConfigsColumns = []string{
		idColumn,
		triggerEventColumn,
		belongsToWebhookColumn,
		createdAtColumn,
		archivedAtColumn,
	}
)

func buildWebhookTriggerConfigsQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(webhookTriggerConfigsColumns)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreateWebhookTriggerConfig",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					webhookTriggerConfigsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveWebhookTriggerConfig",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					webhookTriggerConfigsTableName,
					archivedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
					belongsToWebhookColumn, belongsToWebhookColumn,
				)),
			},
		}
	default:
		return nil
	}
}

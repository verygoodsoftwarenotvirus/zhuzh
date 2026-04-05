package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	waitlistsTableName = "waitlists"
)

func init() {
	registerTableName(waitlistsTableName)
}

var waitlistsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	"valid_until",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildWaitlistsQueries(database string) []*Query {
	switch database {
	case postgres, sqlite:
		insertColumns := filterForInsert(waitlistsColumns)
		fullSelectColumns := applyToEach(waitlistsColumns, func(_ int, s string) string {
			return fullColumnName(waitlistsTableName, s)
		})

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreateWaitlist",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					waitlistsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateWaitlist",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					waitlistsTableName,
					strings.Join(applyToEach(filterForUpdate(waitlistsColumns), func(_ int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression(database),
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveWaitlist",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					waitlistsTableName,
					lastUpdatedAtColumn, currentTimeExpression(database),
					archivedAtColumn, currentTimeExpression(database),
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetWaitlist",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					waitlistsTableName,
					waitlistsTableName, archivedAtColumn,
					waitlistsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckWaitlistExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					waitlistsTableName, idColumn,
					waitlistsTableName,
					waitlistsTableName, archivedAtColumn,
					waitlistsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "WaitlistIsNotExpired",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.valid_until >= %s
);`,
					waitlistsTableName, idColumn,
					waitlistsTableName,
					waitlistsTableName, archivedAtColumn,
					waitlistsTableName, idColumn, idColumn,
					waitlistsTableName, currentTimeExpression(database),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetWaitlists",
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
					buildFilterCountSelect(waitlistsTableName, database, true, true, nil),
					buildTotalCountSelect(waitlistsTableName, true, nil),
					waitlistsTableName,
					waitlistsTableName, archivedAtColumn,
					buildFilterConditions(waitlistsTableName, database, true, false),
					buildCursorLimitClause(waitlistsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetActiveWaitlists",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.valid_until >= %s
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(waitlistsTableName, database, true, true, nil, fmt.Sprintf("%s.valid_until >= %s", waitlistsTableName, currentTimeExpression(database))),
					buildTotalCountSelect(waitlistsTableName, true, nil, fmt.Sprintf("%s.valid_until >= %s", waitlistsTableName, currentTimeExpression(database))),
					waitlistsTableName,
					waitlistsTableName, archivedAtColumn,
					waitlistsTableName, currentTimeExpression(database),
					buildFilterConditions(waitlistsTableName, database, true, false, fmt.Sprintf("%s.valid_until >= %s", waitlistsTableName, currentTimeExpression(database))),
					buildCursorLimitClause(waitlistsTableName),
				)),
			},
		}
	default:
		return nil
	}
}

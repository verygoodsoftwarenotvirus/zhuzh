package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	permissionsTableName = "permissions"
)

func init() {
	registerTableName(permissionsTableName)
}

var permissionsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildPermissionsQueries(database string) []*Query {
	switch database {
	case postgres:

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreatePermission",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					permissionsTableName,
					strings.Join(filterForInsert(permissionsColumns), ",\n\t"),
					strings.Join(applyToEach(filterForInsert(permissionsColumns), func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPermissionByID",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(permissionsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", permissionsTableName, s)
					}), ",\n\t"),
					permissionsTableName,
					permissionsTableName, archivedAtColumn,
					permissionsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPermissionByName",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(permissionsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", permissionsTableName, s)
					}), ",\n\t"),
					permissionsTableName,
					permissionsTableName, archivedAtColumn,
					permissionsTableName, nameColumn, nameColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetPermissions",
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
					strings.Join(applyToEach(permissionsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", permissionsTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(permissionsTableName, true, true, []string{}),
					buildTotalCountSelect(permissionsTableName, true, []string{}),
					permissionsTableName,
					permissionsTableName, archivedAtColumn,
					buildFilterConditions(permissionsTableName, true, true),
					buildCursorLimitClause(permissionsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchivePermission",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					permissionsTableName,
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

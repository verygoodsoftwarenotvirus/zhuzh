package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	userRolesTableName = "user_roles"

	scopeColumn = "scope"
)

func init() {
	registerTableName(userRolesTableName)
}

var userRolesColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	scopeColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUserRolesQueries(database string) []*Query {
	switch database {
	case postgres, sqlite:

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreateUserRole",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					userRolesTableName,
					strings.Join(filterForInsert(userRolesColumns), ",\n\t"),
					strings.Join(applyToEach(filterForInsert(userRolesColumns), func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUserRoleByID",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(userRolesColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", userRolesTableName, s)
					}), ",\n\t"),
					userRolesTableName,
					userRolesTableName, archivedAtColumn,
					userRolesTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUserRoleByName",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(userRolesColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", userRolesTableName, s)
					}), ",\n\t"),
					userRolesTableName,
					userRolesTableName, archivedAtColumn,
					userRolesTableName, nameColumn, nameColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUserRoles",
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
					strings.Join(applyToEach(userRolesColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", userRolesTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(userRolesTableName, database, true, true, []string{}),
					buildTotalCountSelect(userRolesTableName, true, []string{}),
					userRolesTableName,
					userRolesTableName, archivedAtColumn,
					buildFilterConditions(userRolesTableName, database, true, true),
					buildCursorLimitClause(userRolesTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveUserRole",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					userRolesTableName,
					archivedAtColumn, currentTimeExpression(database),
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}

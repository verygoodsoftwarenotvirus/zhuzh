package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	productsTableName = "products"
)

func init() {
	registerTableName(productsTableName)
}

var productsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	"kind",
	"amount_cents",
	"currency",
	"billing_interval_months",
	"external_product_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildPaymentsProductsQueries(database string) []*Query {
	switch database {
	case postgres:
		insertColumns := filterForInsert(productsColumns)
		fullSelectColumns := applyToEach(productsColumns, func(_ int, s string) string {
			return fullColumnName(productsTableName, s)
		})

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveProduct",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s, %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					productsTableName,
					archivedAtColumn, currentTimeExpression,
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateProduct",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					productsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckProductExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
);`,
					productsTableName, idColumn,
					productsTableName,
					productsTableName, archivedAtColumn,
					productsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetProduct",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					productsTableName,
					productsTableName, archivedAtColumn,
					productsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetProductByExternalID",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.external_product_id = sqlc.arg(external_product_id);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					productsTableName,
					productsTableName, archivedAtColumn,
					productsTableName,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetProducts",
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
					buildFilterCountSelect(productsTableName, true, true, []string{}),
					buildTotalCountSelect(productsTableName, true, []string{}),
					productsTableName,
					productsTableName, archivedAtColumn,
					buildFilterConditions(productsTableName, true, true),
					buildCursorLimitClause(productsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "SearchForProducts",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s %s
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(productsTableName, true, true, []string{}),
					buildTotalCountSelect(productsTableName, true, []string{}),
					productsTableName,
					productsTableName, archivedAtColumn,
					productsTableName, nameColumn, buildILIKEForArgument("name_query"),
					buildFilterConditions(productsTableName, true, true),
					buildCursorLimitClause(productsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateProduct",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					productsTableName,
					strings.Join(applyToEach(filterForUpdate(productsColumns), func(_ int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
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

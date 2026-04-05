package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	uploadedMediaTableName = "uploaded_media"
	storagePathColumn      = "storage_path"
	mimeTypeColumn         = "mime_type"
)

func init() {
	registerTableName(uploadedMediaTableName)
}

var uploadedMediaColumns = []string{
	idColumn,
	storagePathColumn,
	mimeTypeColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	createdByUserColumn,
}

func buildUploadedMediaQueries(database string) []*Query {
	switch database {
	case postgres, sqlite:
		insertColumns := filterForInsert(uploadedMediaColumns)
		fullSelectColumns := applyToEach(uploadedMediaColumns, func(_ int, s string) string {
			return fullColumnName(uploadedMediaTableName, s)
		})

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "CreateUploadedMedia",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					uploadedMediaTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateUploadedMedia",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					uploadedMediaTableName,
					strings.Join(applyToEach(filterForUpdate(uploadedMediaColumns, createdByUserColumn), func(_ int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression(database),
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveUploadedMedia",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					uploadedMediaTableName,
					lastUpdatedAtColumn, currentTimeExpression(database),
					archivedAtColumn, currentTimeExpression(database),
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUploadedMedia",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					uploadedMediaTableName,
					uploadedMediaTableName, archivedAtColumn,
					uploadedMediaTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckUploadedMediaExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					uploadedMediaTableName, idColumn,
					uploadedMediaTableName,
					uploadedMediaTableName, archivedAtColumn,
					uploadedMediaTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUploadedMediaWithIDs",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					uploadedMediaTableName,
					uploadedMediaTableName, archivedAtColumn,
					anyInExpression(database, fmt.Sprintf("%s.%s", uploadedMediaTableName, idColumn), "ids"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetUploadedMediaForUser",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(uploadedMediaTableName, database, true, true, nil, fmt.Sprintf("%s.%s = sqlc.arg(%s)", uploadedMediaTableName, createdByUserColumn, createdByUserColumn)),
					buildTotalCountSelect(uploadedMediaTableName, true, nil, fmt.Sprintf("%s.%s = sqlc.arg(%s)", uploadedMediaTableName, createdByUserColumn, createdByUserColumn)),
					uploadedMediaTableName,
					uploadedMediaTableName, archivedAtColumn,
					uploadedMediaTableName, createdByUserColumn, createdByUserColumn,
					buildFilterConditions(uploadedMediaTableName, database, true, false, fmt.Sprintf("%s.%s = sqlc.arg(%s)", uploadedMediaTableName, createdByUserColumn, createdByUserColumn)),
					buildCursorLimitClause(uploadedMediaTableName),
				)),
			},
		}
	default:
		return nil
	}
}

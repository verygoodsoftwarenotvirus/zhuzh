package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

var (
	allTablesHat sync.Mutex
	allTables    = map[string]bool{}
)

func registerTableName(table string) {
	allTablesHat.Lock()
	defer allTablesHat.Unlock()
	allTables[table] = true
}

func getAllTables() []string {
	allTablesHat.Lock()
	defer allTablesHat.Unlock()

	tables := make([]string, 0, len(allTables))
	for t := range allTables {
		tables = append(tables, t)
	}

	slices.Sort(tables)

	return tables
}

func buildDestroyAllDataContent(database string) string {
	switch database {
	case sqlite:
		tables := getAllTables()
		stmts := make([]string, 0, len(tables))
		for _, t := range tables {
			stmts = append(stmts, fmt.Sprintf("DELETE FROM %s;", t))
		}
		return strings.Join(stmts, "\n")
	default:
		return fmt.Sprintf(`TRUNCATE %s CASCADE;`, strings.Join(getAllTables(), ", "))
	}
}

func buildMaintenanceQueries(database string) []*Query {
	switch database {
	case postgres, sqlite:
		oneDayAgo := pastIntervalExpression(database, "1 day")
		queries := []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "DeleteExpiredOAuth2ClientTokens",
					Type: ExecRowsType,
				},
				Content: fmt.Sprintf(`DELETE FROM %s WHERE %s < %s AND %s < %s AND %s < %s;`, oauth2ClientTokensTableName, codeExpiresAtColumn, oneDayAgo, accessExpiresAtColumn, oneDayAgo, refreshExpiresAtColumn, oneDayAgo),
			},
			{
				Annotation: QueryAnnotation{
					Name: "DestroyAllData",
					Type: ExecType,
				},
				Content: buildDestroyAllDataContent(database),
			},
		}
		return append(queries, buildQueueTestMessagesQueries(database)...)
	default:
		return nil
	}
}

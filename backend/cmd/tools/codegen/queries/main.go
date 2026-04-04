package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

var (
	checkOnlyFlag = pflag.Bool("check", false, "only check if files match")
	databaseFlag  = pflag.String("database", postgres, "what database to use")
)

const (
	postgres = "postgres"
)

type queryEntry struct {
	queries func(string) []*Query
	path    string
}

var queryEntries = []queryEntry{
	{queries: buildMaintenanceQueries, path: "internalops/sqlc_queries/internalops"},
	{queries: buildOAuth2ClientTokensQueries, path: "oauth/sqlc_queries/oauth2_client_tokens"},
	{queries: buildOAuth2ClientsQueries, path: "oauth/sqlc_queries/oauth2_clients"},
	{queries: buildAccountInvitationsQueries, path: "identity/sqlc_queries/account_invitations"},
	{queries: buildAccountUserMembershipsQueries, path: "identity/sqlc_queries/account_user_memberships"},
	{queries: buildAccountsQueries, path: "identity/sqlc_queries/accounts"},
	{queries: buildAuditLogEntryQueries, path: "auditlogentries/sqlc_queries/audit_logs"},
	{queries: buildAdminQueries, path: "identity/sqlc_queries/admin"},
	{queries: buildPasswordResetTokensQueries, path: "auth/sqlc_queries/password_reset_tokens"},
	{queries: buildUserSessionsQueries, path: "auth/sqlc_queries/user_sessions"},
	{queries: buildUsersQueries, path: "identity/sqlc_queries/users"},
	{queries: buildServiceSettingQueries, path: "settings/sqlc_queries/service_settings"},
	{queries: buildServiceSettingConfigurationQueries, path: "settings/sqlc_queries/service_setting_configurations"},
	{queries: buildWebhooksQueries, path: "webhooks/sqlc_queries/webhooks"},
	{queries: buildWebhookTriggerEventsQueries, path: "webhooks/sqlc_queries/webhook_trigger_events"},
	{queries: buildWebhookTriggerConfigsQueries, path: "webhooks/sqlc_queries/webhook_trigger_configs"},
	{queries: buildUserNotificationQueries, path: "notifications/sqlc_queries/user_notifications"},
	{queries: buildWaitlistsQueries, path: "waitlists/sqlc_queries/waitlists"},
	{queries: buildWaitlistSignupsQueries, path: "waitlists/sqlc_queries/waitlist_signups"},
	{queries: buildIssueReportsQueries, path: "issuereports/sqlc_queries/issue_reports"},
	{queries: buildUploadedMediaQueries, path: "uploadedmedia/sqlc_queries/uploaded_media"},
	{queries: buildUserDataDisclosuresQueries, path: "dataprivacy/sqlc_queries/user_data_disclosures"},
	{queries: buildPaymentsProductsQueries, path: "payments/sqlc_queries/products"},
	{queries: buildPaymentsSubscriptionsQueries, path: "payments/sqlc_queries/subscriptions"},
	{queries: buildPaymentsPurchasesQueries, path: "payments/sqlc_queries/purchases"},
	{queries: buildPaymentsTransactionsQueries, path: "payments/sqlc_queries/payment_transactions"},
	{queries: buildCommentsQueries, path: "comments/sqlc_queries/comments"},
	{queries: buildUserRolesQueries, path: "identity/sqlc_queries/user_roles"},
	{queries: buildPermissionsQueries, path: "identity/sqlc_queries/permissions"},
	{queries: buildUserRolePermissionsQueries, path: "identity/sqlc_queries/user_role_permissions"},
	{queries: buildUserRoleAssignmentsQueries, path: "identity/sqlc_queries/user_role_assignments"},
	{queries: buildUserRoleHierarchyQueries, path: "identity/sqlc_queries/user_role_hierarchy"},
}

func repoDir(database string) string {
	switch database {
	case sqlite:
		return "internal/repositories/sqlite/"
	default:
		return "internal/repositories/postgres/"
	}
}

func main() {
	pflag.Parse()

	runErrors := &multierror.Error{}
	databaseToUse := *databaseFlag

	baseDir := repoDir(databaseToUse)
	queryOutput := make(map[string][]*Query, len(queryEntries))
	for _, e := range queryEntries {
		queryOutput[baseDir+e.path] = e.queries(databaseToUse)
	}

	checkOnly := *checkOnlyFlag

	for filePath, queries := range queryOutput {
		filePath += ".generated.sql"
		existingFile, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if _, err = os.Create(filePath); err != nil {
					log.Fatal(fmt.Errorf("creating file: %w", err))
				}
			}
			if err != nil {
				log.Fatal(fmt.Errorf("opening existing file: %w", err))
			}
		}

		var fileContent strings.Builder
		for i, query := range queries {
			if i != 0 {
				fileContent.WriteString("\n")
			}
			fileContent.WriteString(query.Render())
		}

		var fileOutput strings.Builder
		for line := range strings.SplitSeq(strings.TrimSpace(fileContent.String()), "\n") {
			fileOutput.WriteString(strings.TrimSuffix(line, " ") + "\n")
		}

		if string(existingFile) != fileOutput.String() && checkOnly {
			runErrors = multierror.Append(runErrors, fmt.Errorf("files don't match: %s", filePath))
		}

		if !checkOnly {
			if err = os.WriteFile(filePath, []byte(fileOutput.String()), 0o600); err != nil {
				runErrors = multierror.Append(runErrors, err)
			}
		}
	}

	if runErrors.ErrorOrNil() != nil {
		log.Fatal(runErrors)
	}
}

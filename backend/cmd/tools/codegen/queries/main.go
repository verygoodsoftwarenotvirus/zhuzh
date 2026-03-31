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

func main() {
	pflag.Parse()

	runErrors := &multierror.Error{}
	databaseToUse := *databaseFlag

	queryOutput := map[string][]*Query{
		"internal/repositories/postgres/internalops/sqlc_queries/internalops":    buildMaintenanceQueries(databaseToUse),
		"internal/repositories/postgres/oauth/sqlc_queries/oauth2_client_tokens": buildOAuth2ClientTokensQueries(databaseToUse),
		"internal/repositories/postgres/oauth/sqlc_queries/oauth2_clients":                                      buildOAuth2ClientsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/account_invitations":                              buildAccountInvitationsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/account_user_memberships":                         buildAccountUserMembershipsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/accounts":                                         buildAccountsQueries(databaseToUse),
		"internal/repositories/postgres/auditlogentries/sqlc_queries/audit_logs":                                buildAuditLogEntryQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/admin":                                            buildAdminQueries(databaseToUse),
		"internal/repositories/postgres/auth/sqlc_queries/password_reset_tokens":                                buildPasswordResetTokensQueries(databaseToUse),
		"internal/repositories/postgres/auth/sqlc_queries/user_sessions":                                        buildUserSessionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/users":                                            buildUsersQueries(databaseToUse),
		"internal/repositories/postgres/settings/sqlc_queries/service_settings":                                 buildServiceSettingQueries(databaseToUse),
		"internal/repositories/postgres/settings/sqlc_queries/service_setting_configurations":                   buildServiceSettingConfigurationQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhooks":                                         buildWebhooksQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhook_trigger_events":                           buildWebhookTriggerEventsQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhook_trigger_configs":                          buildWebhookTriggerConfigsQueries(databaseToUse),
		"internal/repositories/postgres/notifications/sqlc_queries/user_notifications":                          buildUserNotificationQueries(databaseToUse),
		"internal/repositories/postgres/waitlists/sqlc_queries/waitlists":                                       buildWaitlistsQueries(databaseToUse),
		"internal/repositories/postgres/waitlists/sqlc_queries/waitlist_signups":                                buildWaitlistSignupsQueries(databaseToUse),
		"internal/repositories/postgres/issuereports/sqlc_queries/issue_reports":                                buildIssueReportsQueries(databaseToUse),
		"internal/repositories/postgres/uploadedmedia/sqlc_queries/uploaded_media":                              buildUploadedMediaQueries(databaseToUse),
		"internal/repositories/postgres/dataprivacy/sqlc_queries/user_data_disclosures":                         buildUserDataDisclosuresQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/products":                                         buildPaymentsProductsQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/subscriptions":                                    buildPaymentsSubscriptionsQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/purchases":                                        buildPaymentsPurchasesQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/payment_transactions":                             buildPaymentsTransactionsQueries(databaseToUse),
		"internal/repositories/postgres/comments/sqlc_queries/comments":                                         buildCommentsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_roles":                                       buildUserRolesQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/permissions":                                      buildPermissionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_permissions":                            buildUserRolePermissionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_assignments":                            buildUserRoleAssignmentsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_hierarchy":                              buildUserRoleHierarchyQueries(databaseToUse),
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

package main

import (
	"fmt"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"

	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
)

const (
	defaultHTTPPort = 8000
	defaultGRPCPort = 8001
	maxAttempts     = 50
	otelServiceName = "api_server"

	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
	dataChangesTopicName              = "data_changes"
	outboundEmailsTopicName           = "outbound_emails"
	searchIndexRequestsTopicName      = "search_index_requests"
	mobileNotificationsTopicName      = "mobile_notifications"
	userDataAggregationTopicName      = "user_data_aggregation_requests"
	webhookExecutionRequestsTopicName = "webhook_execution_requests"
)

var (
	contentTypeJSON = encoding.ContentTypeToString(encoding.ContentTypeJSON)
)

func main() {
	// localdev config is generated to two locations:
	// - config_files/ for docker-compose usage
	// - kustomize/configs/ for Kubernetes usage (hostnames overridden via env vars)
	localdevConfig := buildLocalDevConfig()

	envConfigs := map[string]*config.EnvironmentConfigSet{
		"deploy/environments/localdev/config_files": {
			RootConfig:                 localdevConfig,
			MCPServiceHTTPAPIServerURL: "http://localhost:8000",
			MCPServiceGRPCAPIServerURL: ":8001",
		},
		"deploy/environments/localdev/kustomize/configs": {
			RootConfig:                 localdevConfig,
			MCPServiceHTTPAPIServerURL: "http://localhost:8000",
			MCPServiceGRPCAPIServerURL: ":8001",
		},
		"deploy/environments/testing/config_files": {
			APIServiceConfigPath: "integration-tests-config.json",
			RootConfig:           buildIntegrationTestsConfig(),
		},
		"deploy/environments/prod/kustomize/configs": {
			RootConfig:                 buildProdConfig(),
			MCPServiceHTTPAPIServerURL: "https://http-api.zhuzh.dev",
			MCPServiceGRPCAPIServerURL: "https://api.zhuzh.dev",
			ServiceDatabaseUsers: map[string]string{
				"db_cleaner":                    "db_cleaner",
				"search_data_index_scheduler":   "search_data_index_scheduler",
				"mobile_notification_scheduler": "mobile_notification_scheduler",
				"async_message_handler":         "async_message_handler",
				"queue_test":                    "queue_test",
			},
		},
	}

	for p, cfg := range envConfigs {
		if err := cfg.Render(p, true, true); err != nil {
			panic(fmt.Errorf("validating config %s: %w", p, err))
		}
	}
}

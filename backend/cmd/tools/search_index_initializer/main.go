package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/identity"
	identityindexing "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/indexing"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	textsearch "github.com/verygoodsoftwarenotvirus/platform/v4/search/text"
	"github.com/verygoodsoftwarenotvirus/platform/v4/search/text/algolia"
	textsearchcfg "github.com/verygoodsoftwarenotvirus/platform/v4/search/text/config"

	"github.com/spf13/cobra"
)

const (
	defaultBatchSize = 50
)

func main() {
	var (
		databaseURL    string
		searchProvider string
		algoliaAppID   string
		algoliaAPIKey  string
	)

	root := &cobra.Command{
		Use:   "search-index-initializer",
		Short: "Initialize search indices from database (for use with proxied production DB)",
	}

	root.PersistentFlags().StringVar(&databaseURL, "database-url", "", "Postgres connection URL (or set DATABASE_URL)")
	root.PersistentFlags().StringVar(&searchProvider, "search-provider", textsearchcfg.AlgoliaProvider, "Search provider: algolia or elasticsearch")
	root.PersistentFlags().StringVar(&algoliaAppID, "algolia-app-id", "", "Algolia app ID (or set ALGOLIA_APP_ID)")
	root.PersistentFlags().StringVar(&algoliaAPIKey, "algolia-api-key", "", "Algolia API key (or set ALGOLIA_API_KEY)")

	root.AddCommand(initCmd(&databaseURL, &searchProvider, &algoliaAppID, &algoliaAPIKey))

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initCmd(databaseURL, searchProvider, algoliaAppID, algoliaAPIKey *string) *cobra.Command {
	var indices string
	var wipe bool
	var batchSize int

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Load all data from database into search indices",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInit(*databaseURL, *searchProvider, *algoliaAppID, *algoliaAPIKey, indices, wipe, batchSize)
		},
	}

	cmd.Flags().StringVar(&indices, "indices", "", "Comma-separated indices to initialize (e.g. recipes,meals,users)")
	cmd.Flags().BoolVar(&wipe, "wipe", false, "Wipe index before reindexing")
	cmd.Flags().IntVar(&batchSize, "batch-size", defaultBatchSize, "Page size for cursor pagination")

	if err := cmd.MarkFlagRequired("indices"); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func runInit(databaseURL, searchProvider, algoliaAppID, algoliaAPIKey, indicesStr string, wipe bool, batchSize int) error {
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		return fmt.Errorf("--database-url or DATABASE_URL is required")
	}

	if algoliaAppID == "" {
		algoliaAppID = os.Getenv("ALGOLIA_APP_ID")
	}
	if algoliaAPIKey == "" {
		algoliaAPIKey = os.Getenv("ALGOLIA_API_KEY")
	}
	if searchProvider == textsearchcfg.AlgoliaProvider && (algoliaAppID == "" || algoliaAPIKey == "") {
		return fmt.Errorf("--algolia-app-id and --algolia-api-key (or env vars) are required for Algolia")
	}

	indices := strings.Split(strings.TrimSpace(indicesStr), ",")
	var trimmed []string
	for _, idx := range indices {
		if s := strings.TrimSpace(idx); s != "" {
			trimmed = append(trimmed, s)
		}
	}
	indices = trimmed
	if len(indices) == 0 {
		return fmt.Errorf("at least one index is required in --indices")
	}

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	tracerProvider := tracing.NewNoopTracerProvider()
	metricsProvider := metrics.NewNoopMetricsProvider()

	dbConfig := &databasecfg.Config{
		Provider:        databasecfg.ProviderPostgres,
		MaxPingAttempts: 10,
		PingWaitPeriod:  time.Second,
	}
	if err := dbConfig.LoadConnectionDetailsFromURL(databaseURL); err != nil {
		return fmt.Errorf("loading database config: %w", err)
	}
	dbConfig.WriteConnection = dbConfig.ReadConnection

	client, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, dbConfig, nil)
	if err != nil {
		return fmt.Errorf("initializing database client: %w", err)
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("closing database client", closeErr)
		}
	}()

	auditRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, client)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditRepo, client)

	searchCfg := &textsearchcfg.Config{
		Provider: searchProvider,
		Algolia: &algolia.Config{
			AppID:  algoliaAppID,
			APIKey: algoliaAPIKey,
		},
	}

	if batchSize < 1 {
		batchSize = defaultBatchSize
	}

	userIndexer, err := buildIndexers(ctx, logger, tracerProvider, metricsProvider, searchCfg, identityRepo)
	if err != nil {
		return fmt.Errorf("building indexers: %w", err)
	}

	for _, indexType := range indices {
		if err = runIndex(ctx, logger, indexType, identityRepo, userIndexer, searchCfg, wipe, batchSize); err != nil {
			return fmt.Errorf("indexing %s: %w", indexType, err)
		}
	}

	return nil
}

func buildIndexers(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	searchCfg *textsearchcfg.Config,
	identityRepo identity.Repository,
) (*identityindexing.UserDataIndexer, error) {
	userIdx, err := textsearchcfg.ProvideIndex[identityindexing.UserSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchCfg, identityindexing.IndexTypeUsers)
	if err != nil {
		return nil, err
	}

	userIndexer := identityindexing.NewCoreDataIndexer(logger, tracerProvider, identityRepo, userIdx)

	return userIndexer, nil
}

func runIndex(
	ctx context.Context,
	logger logging.Logger,
	indexType string,
	identityRepo identity.Repository,
	userIndexer *identityindexing.UserDataIndexer,
	searchCfg *textsearchcfg.Config,
	wipe bool,
	batchSize int,
) error {
	log.Printf("Starting index: %s", indexType)

	im, err := getIndexManager(ctx, logger, indexType, searchCfg)
	if err != nil {
		return err
	}

	if wipe && im != nil {
		log.Printf("Wiping index: %s", indexType)
		if err = im.Wipe(ctx); err != nil {
			return fmt.Errorf("wiping index: %w", err)
		}
		log.Printf("Wiped index: %s", indexType)
	}

	filter := filtering.DefaultQueryFilter()
	pageSize := min(uint8(batchSize), filtering.MaxQueryFilterLimit)
	filter.MaxResponseSize = &pageSize

	var cursor *string
	pageNum := 0
	totalIndexed := 0

	for {
		pageNum++
		filter.Cursor = cursor

		var ids []string
		switch indexType {
		case identityindexing.IndexTypeUsers:
			result, fetchErr := identityRepo.GetUsers(ctx, filter)
			if fetchErr != nil {
				return fmt.Errorf("getting users: %w", fetchErr)
			}
			for _, u := range result.Data {
				ids = append(ids, u.ID)
			}
			if len(result.Data) > 0 {
				cursor = &result.Data[len(result.Data)-1].ID
			} else {
				cursor = nil
			}
		default:
			return fmt.Errorf("unknown index type: %s", indexType)
		}

		if len(ids) == 0 {
			break
		}

		for _, id := range ids {
			req := &textsearch.IndexRequest{RowID: id, IndexType: indexType}
			switch indexType {
			case identityindexing.IndexTypeUsers:
				if err = userIndexer.HandleIndexRequest(ctx, req); err != nil {
					return fmt.Errorf("indexing %s %s: %w", indexType, id, err)
				}
			}
			totalIndexed++
		}

		log.Printf("Indexed page: %s page=%d items=%d total=%d", indexType, pageNum, len(ids), totalIndexed)

		if len(ids) < batchSize {
			break
		}
	}

	log.Printf("Finished index: %s total_indexed=%d", indexType, totalIndexed)
	return nil
}

func getIndexManager(
	ctx context.Context,
	logger logging.Logger,
	indexType string,
	searchCfg *textsearchcfg.Config,
) (textsearch.IndexManager, error) {
	tracerProvider := tracing.NewNoopTracerProvider()
	metricsProvider := metrics.NewNoopMetricsProvider()

	switch indexType {
	case identityindexing.IndexTypeUsers:
		return textsearchcfg.ProvideIndex[identityindexing.UserSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchCfg, indexType)
	default:
		return nil, fmt.Errorf("unknown index type: %s", indexType)
	}
}

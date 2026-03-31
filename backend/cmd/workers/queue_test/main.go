package main

import (
	"context"
	"fmt"
	"log"
	"time"

	queuetest "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/build/jobs/queue_test"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.QueueTestJobConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	result, err := queuetest.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building queue test job: %w", err)
	}
	defer func() {
		// Flush metrics so Prometheus receives queue_test_round_trip_ms before the pod exits.
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if shutdownErr := result.ShutdownMetrics(shutdownCtx); shutdownErr != nil {
			log.Printf("error shutting down metrics: %v", shutdownErr)
		}
	}()

	if err = result.Job.Do(ctx); err != nil {
		return fmt.Errorf("running queue test job: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}

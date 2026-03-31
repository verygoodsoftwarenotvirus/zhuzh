package main

import (
	"context"
	"fmt"
	"log"

	emaildeliverabilitytest "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/build/jobs/email_deliverability_test"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.EmailDeliverabilityTestConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	job, err := emaildeliverabilitytest.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building email deliverability test job: %w", err)
	}

	if err = job.Do(ctx); err != nil {
		return fmt.Errorf("running email deliverability test: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}

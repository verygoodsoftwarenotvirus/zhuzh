package main

import (
	"context"
	"fmt"
	"log"

	mobilenotificationscheduler "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/build/jobs/mobile_notification_scheduler"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.MobileNotificationSchedulerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	scheduler, err := mobilenotificationscheduler.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building scheduler: %w", err)
	}

	if err = scheduler.ScheduleNotifications(ctx); err != nil {
		return fmt.Errorf("error scheduling notifications: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}

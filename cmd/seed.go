package cmd

import (
	"fmt"
	"time"

	"github.com/cvetkovski98/zvax-common/pkg/redis"
	"github.com/cvetkovski98/zvax-slots/internal/config"
	"github.com/cvetkovski98/zvax-slots/internal/model"
	"github.com/cvetkovski98/zvax-slots/internal/repository"
	"github.com/spf13/cobra"
)

var seedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Migrate database",
	Long:  `Migrate database`,
	RunE:  seed,
}

func seed(cmd *cobra.Command, args []string) error {
	var startDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	var endDate = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	var slots = make([]*model.Slot, 0)
	for startDate.Before(endDate) {
		startOfDay := time.Date(
			startDate.Year(),
			startDate.Month(),
			startDate.Day(),
			0, 0, 0, 0,
			startDate.Location(),
		)
		endOfDay := time.Date(
			startDate.Year(),
			startDate.Month(),
			startDate.Day(),
			23, 59, 59, 0,
			startDate.Location(),
		)
		if startDate.Weekday() == time.Saturday || startDate.Weekday() == time.Sunday {
			startDate = startDate.Add(24 * time.Hour)
			continue
		}
		for startOfDay.Before(endOfDay) {
			//available := rand.Intn(100) < 40
			slots = append(slots, &model.Slot{
				DateTime:  startOfDay,
				Location:  "skopje",
				Available: true,
			})
			startOfDay = startOfDay.Add(time.Minute * 30)
		}
		startDate = startDate.Add(24 * time.Hour)
	}

	cfg := config.GetConfig()
	rdb, err := redis.NewRedisConn(cfg.Redis)
	if err != nil {
		return err
	}
	repo := repository.NewRedisSlotRepository(rdb)
	for _, slot := range slots {
		fmt.Println("Inserting slot: ", slot)
		_, err := repo.InsertOne(cmd.Context(), slot)
		if err != nil {
			return err
		}
	}
	return nil
}

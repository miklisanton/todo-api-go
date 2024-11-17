package jobs

import (
	"context"
	"os"
	"time"
	"todo-api/internal/services"

	"github.com/rs/zerolog/log"
)

type IDateWorker interface {
	MonitorDueDate(interval time.Duration)
}

type DateWorker struct {
	TaskService services.ITaskService
	ExitChan    chan os.Signal
}

func NewDateWorker(taskService services.ITaskService, exitChan chan os.Signal) IDateWorker {
	return &DateWorker{taskService, exitChan}
}

func (dw *DateWorker) MonitorDueDate(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), interval)

			if err := dw.TaskService.UpdateOverdue(ctx); err != nil {
				log.Logger.Error().Err(err).Msg("failed to update overdue tasks")
			}
			cancel()
		case <-dw.ExitChan:
			ticker.Stop()
			return
		}
	}
}

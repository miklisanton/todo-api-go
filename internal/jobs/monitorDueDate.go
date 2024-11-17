package jobs

import (
	"context"
	"sync"
	"time"
	"todo-api/internal/services"

	"github.com/rs/zerolog/log"
)

type IDateWorker interface {
	MonitorDueDate(ctx context.Context, interval time.Duration)
	Wait()
}

type DateWorker struct {
	TaskService services.ITaskService
	doneWg      sync.WaitGroup
	Exit        chan int
}

func NewDateWorker(taskService services.ITaskService) IDateWorker {
	return &DateWorker{taskService, sync.WaitGroup{}, make(chan int)}
}

func (dw *DateWorker) MonitorDueDate(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				dw.doneWg.Add(1)
				ctx, cancel := context.WithTimeout(context.Background(), interval)

				if err := dw.TaskService.UpdateOverdue(ctx); err != nil {
					log.Logger.Error().Err(err).Msg("failed to update overdue tasks")
				}
				cancel()
				dw.doneWg.Done()
			case <-ctx.Done():
				dw.doneWg.Wait()
				ticker.Stop()
				log.Logger.Info().Msg("Date worker stopped")
				// Signal main goroutine that this worker is done
				dw.Exit <- 1
				return
			}
		}
	}()
}

func (dw *DateWorker) Wait() {
	<-dw.Exit
}

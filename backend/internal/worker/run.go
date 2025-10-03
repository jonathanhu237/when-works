package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
)

func (w *Worker) Run() error {
	// Create asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", w.config.Redis.Host, w.config.Redis.Port),
			Password: w.config.Redis.Password,
			DB:       w.config.Redis.DB,
		},
		asynq.Config{
			Concurrency: w.config.Asynq.Concurrency,
		},
	)

	// Register handlers
	mux := w.routes()

	// Run worker
	w.logger.Info("worker started")
	if err := srv.Run(mux); err != nil {
		return err
	}

	w.logger.Info("worker stopped")
	return nil
}

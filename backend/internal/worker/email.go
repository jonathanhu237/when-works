package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jonathanhu237/when-works/backend/internal/tasks"
)

func (w *Worker) HandleEmailNewUserTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.EmailNewUserPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	subject := "Welcome to WhenWorks!"

	if err := w.mailer.SendHTML(payload.Email, subject, "new_user.html", payload); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	w.logger.Info("new user email sent successfully", "email", payload.Email)
	return nil
}

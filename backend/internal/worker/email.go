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

	w.logger.Info("sending new user email",
		"email", payload.Email,
		"username", payload.Username,
	)

	subject := "Welcome to WhenWorks!"
	body := fmt.Sprintf(`Hello %s,

Welcome to WhenWorks!

Your account has been created successfully.

Username: %s
Password: %s

Please login and change your password as soon as possible.

Best regards,
WhenWorks Team`, payload.Username, payload.Username, payload.Password)

	if err := w.mailer.SendHTML(payload.Email, subject, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	w.logger.Info("new user email sent successfully", "email", payload.Email)
	return nil
}

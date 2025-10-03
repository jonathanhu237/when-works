package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jonathanhu237/when-works/backend/internal/config"
)

// Task type constants
const (
	TypeEmailNewUser = "email:new_user"
)

// ------------------------------
// Email New User Task
// ------------------------------
type EmailNewUserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewEmailNewUserTask(email, username, password string, cfg config.Config) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailNewUserPayload{
		Email:    email,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TypeEmailNewUser,
		payload,
		asynq.MaxRetry(cfg.Asynq.MaxRetry),
		asynq.Timeout(time.Duration(cfg.Asynq.Timeout)*time.Second)), nil
}

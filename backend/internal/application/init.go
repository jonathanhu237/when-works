package application

import (
	"github.com/jonathanhu237/when-works/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) Init() error {
	// Check if an admin user already exists
	adminExists, err := app.models.User.AdminExists()
	if err != nil {
		return err
	}

	if adminExists {
		app.logger.Info("admin user already exists, skipping initial admin creation")
		return nil
	}

	// Create the initial admin user
	app.logger.Info("creating initial admin user")
	initialAdmin := models.User{
		Username: app.config.InitialAdmin.Username,
		Email:    app.config.InitialAdmin.Email,
		Name:     "Admin",
		IsAdmin:  true,
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(app.config.InitialAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	initialAdmin.PasswordHash = string(passwordHash)

	// Insert the initial admin user into the database
	if err := app.models.User.Insert(&initialAdmin); err != nil {
		return err
	}

	return nil
}

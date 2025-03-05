package application

import (
	"github.com/Irurnnen/calc-master/internal/config"
	"github.com/Irurnnen/calc-master/internal/router"
)

type Application struct {
	Config config.Config
}

func (a *Application) Run() {
	// Setup router
	r := router.New()

	// Start server
	r.Run(a.Config.ServerURL)
}

func New() *Application {
	return &Application{
		Config: *config.NewConfigFromEnv(),
	}
}

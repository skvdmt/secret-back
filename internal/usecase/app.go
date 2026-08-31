package usecase

import (
	"context"

	"github.com/skvdmt/secret-back/internal/model"
)

const (
	SECRET = "This secret text require authentication."
)

// App Сервисный слой.
type App struct {
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("usecase layer creating")
	// Создание репозиторного слоя.
	return &App{}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	model.Log.Info.Info("usecase layer stopped")
	return nil
}

// Secret Секрет.
func (a *App) Secret() string {
	return SECRET
}

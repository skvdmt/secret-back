package usecase

import (
	"context"

	"github.com/skvdmt/secret-back/internal/model"
	"github.com/skvdmt/secret-back/internal/repository"
)

const (
	SECRET = "This secret text require authentication."
)

// App Сервисный слой.
type App struct {
	repository Repository
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("usecase layer creating")
	// Создание репозиторного слоя.
	rep, err := repository.NewApp()
	if err != nil {
		return nil, err
	}
	return &App{
		repository: rep,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("usecase layer stopped")
	return nil
}

// Secret Секрет.
func (a *App) Secret() string {
	return SECRET
}

// Documentation Документация.
func (a *App) Documentation(ctx context.Context) []byte {
	return a.repository.Documentation(ctx)
}

package repository

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/skvdmt/secret-back/internal/model"
)

// App Репозиторный слой.
type App struct {
	muDocumentation *sync.RWMutex
	documentation   []byte
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("repository layer creating")
	a := &App{
		muDocumentation: &sync.RWMutex{},
	}
	if err := a.updateDocumentation(); err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	model.Log.Info.Info("repository layer stopped")
	return nil
}

// Swagger Документация.
func (a *App) Documentation(ctx context.Context) []byte {
	a.muDocumentation.RLock()
	defer a.muDocumentation.RUnlock()
	return a.documentation
}

// updateSwagger Основление документации.
func (a *App) updateDocumentation() error {
	model.Log.Info.Info("update documentation starting")
	const (
		name = "openapi.yaml"
		dev  = "."
		prod = "/usr/local/share/doc"
	)
	m, o := os.LookupEnv(model.MODE)
	dir := prod
	if o && m == model.MODE_DEV {
		dir = dev
	}
	doc, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	a.documentation = doc
	return nil
}

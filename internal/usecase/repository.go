package usecase

import "context"

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Документация.
	Documentation(ctx context.Context) []byte
}

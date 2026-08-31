package internal

import (
	"context"
	"net/http"
)

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Секрет.
	Secret(w http.ResponseWriter, r *http.Request)
	// Авторизация.
	Authorization(next http.Handler) http.Handler
}

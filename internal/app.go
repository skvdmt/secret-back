package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/skvdmt/secret-back/internal/delivery"
	"github.com/skvdmt/secret-back/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов прерывания.
	interrupt chan os.Signal
	// Контекст.
	ctx context.Context
	// Отмена контекста.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Ошибки приложения.
	errors []error
	// Приложение закрывается.
	stopping bool
	// Корректное завершение горутин.
	wg *sync.WaitGroup
	// Роутер.
	router *http.ServeMux
	// Сервер.
	server *http.Server
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info(fmt.Sprintf("app %s creating", model.NAME))
	var err error
	// Загрузка конфигурации.
	if model.Conf, err = model.NewConfig(); err != nil {
		return nil, err
	}
	// Создание глобального канала ошибок.
	model.Errors = make(chan error)
	// Создание транспортного слоя.
	d, err := delivery.NewApp()
	if err != nil {
		return nil, err
	}
	r := http.NewServeMux()
	// Создание контекста.
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:       ctx,
		cancel:    cancel,
		interrupt: make(chan os.Signal),
		delivery:  d,
		wg:        &sync.WaitGroup{},
		router:    r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", model.Conf.Server.Port),
			Handler: r,
		},
	}, nil
}

// Start Запуск.
func (a *App) Start() []error {
	model.Log.Info.Info(fmt.Sprintf("app %s starting", model.NAME))

	a.wg.Add(1)
	go a.errorHandler()

	a.wg.Add(1)
	go func() {
		model.Log.Info.Info(fmt.Sprintf("http server starting on %d port",
			model.Conf.Server.Port))
		defer func() {
			model.Log.Info.Info("http server stopped")
			a.wg.Done()
		}()
		a.routing()
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			model.Errors <- err
		}
	}()

	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Log.Info.Info("error handler starting")
	defer func() {
		model.Log.Info.Info("error handler stopped")
		a.wg.Done()
	}()
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.errors = append(a.errors, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик канала сигналов прерывания.
func (a *App) interruptHandler() []error {
	model.Log.Info.Info("waiting for interrupt signals")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	a.stop()
	close(model.Errors)
	return a.errors
}

// stop Остановка.
func (a *App) stop() {
	defer model.Log.Info.Info(fmt.Sprintf("app %s stopped", model.NAME))
	a.stopping = true
	// Отмена контекста.
	a.cancel()
	model.Log.Info.Info("context canceled")
	// Остановка сервера.
	if err := a.server.Shutdown(context.Background()); err != nil {
		model.Errors <- err
	}
	// Остановка транспортного слоя.
	if err := a.delivery.Stop(a.ctx); err != nil {
		model.Errors <- err
	}
	// Закрытие канала сигналов прерывания.
	close(a.interrupt)
}

// routing Маршрутизация.
func (a *App) routing() {
	a.router.Handle("GET /text", a.delivery.Authorization(http.HandlerFunc(a.delivery.Secret)))
}

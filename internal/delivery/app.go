package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	valid_access_v1 "github.com/skvdmt/secret-back/gen/go/valid_access"
	"github.com/skvdmt/secret-back/internal/model"
	"github.com/skvdmt/secret-back/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	CONTENT_TYPE                 = "Content-type"
	APPLICATION_JSON             = "application/json"
	AUTHORIZATION_DIRECTORY_DEV  = "./certs/valid-access-back"
	AUTHORIZATION_DIRECTORY_PROD = "/etc/ssl/certs/valid-access-back"
	AUTHORIZATION_FILENAME_CRT   = "cert.crt"
)

// App Транспортный слой.
type App struct {
	// Сервисный слой.
	usecase Usecase
	// Клиент валидатора токена доступа
	client valid_access_v1.ValidatorClient
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("delivery layer creating")

	c, err := credentials.NewClientTLSFromFile(certFile(), "")
	if err != nil {
		return nil, err
	}
	con, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%d",
			model.Conf.ValidAccess.Host,
			model.Conf.ValidAccess.Port,
		),
		grpc.WithTransportCredentials(c),
	)
	if err != nil {
		return nil, err
	}
	// Создание сервисного слоя.
	u, err := usecase.NewApp()
	if err != nil {
		return nil, err
	}
	return &App{
		usecase: u,
		client:  valid_access_v1.NewValidatorClient(con),
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("delivery layer stopped")
	return nil
}

// Documentation Документация.
func (a *App) Documentation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(a.usecase.Documentation(r.Context()))
}

// Secret Секрет.
func (a *App) Secret(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(struct {
		Secret string `json:"secret"`
	}{
		Secret: a.usecase.Secret(),
	})
}

// Authorization Авторизация.
func (a *App) Authorization(next http.Handler) http.Handler {
	const BEARER = "Bearer "
	const AUTHORIZATION = "Authorization"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model.Log.Info.Info("Auth")
		w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
		auth := r.Header.Get(AUTHORIZATION)
		if !strings.HasPrefix(auth, BEARER) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{
				Message: fmt.Sprintf("%s %snot found", AUTHORIZATION, BEARER),
			})
			return
		}
		t := strings.TrimPrefix(auth, BEARER)
		_, err := a.client.AccessTokenValidation(r.Context(), &valid_access_v1.AccessToken{
			Token: t,
		})
		if err != nil {
			s, _ := status.FromError(err)
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{
				Message: s.Message(),
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// certFile Файл сертификата.
func certFile() string {
	m, o := os.LookupEnv(model.MODE)
	if o && m == model.MODE_DEV {
		return filepath.Join(AUTHORIZATION_DIRECTORY_DEV, AUTHORIZATION_FILENAME_CRT)
	}
	return filepath.Join(AUTHORIZATION_DIRECTORY_PROD, AUTHORIZATION_FILENAME_CRT)
}

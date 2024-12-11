package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/books"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans/repo"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/users"
	"golang.org/x/sync/errgroup"
)

type App struct {
	config         *Config
	router         *chi.Mux
	http           *http.Server
	routerInternal *chi.Mux
	httpInternal   *http.Server
}

func New(ctx context.Context, config *Config) (*App, error) {
	routerPub, httpPub := makeServer(config.PublicURL)
	routerInt, httpInt := makeServer(config.PrivateURL)

	return &App{
		config:         config,
		router:         routerPub,
		http:           httpPub,
		routerInternal: routerInt,
		httpInternal:   httpInt,
	}, nil
}

func makeServer(address string) (*chi.Mux, *http.Server) {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    0x10000,
	}
	return router, server
}

// Setup configures the application
func (a *App) Setup(ctx context.Context) error {
	userSvc := users.NewConn(a.config.UserServiceURL)
	bookSvc := books.NewConn(a.config.BookServiceURL)

	dsn := a.config.DSN
	var store loans.Repo
	switch {
	case strings.HasPrefix(dsn, "memory://"):
		store = repo.NewMemoryRepo(dsn)
	case strings.HasPrefix(dsn, "sqlite://"):
		var err error
		store, err = repo.NewSqliteRepo(dsn)
		if err != nil {
			return err
		}
	default:
		return fail.ErrInvalidDSN
	}

	service := loans.NewService(store, userSvc, bookSvc, a.config.BookReturnDeadline)
	handler := loans.NewHandler(a.router, a.routerInternal, service)
	handler.Register()

	return nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errs, ctx := errgroup.WithContext(ctx)

	log.Printf("starting web servers: public on %s, private on %s\n", a.config.PublicURL, a.config.PrivateURL)

	errs.Go(func() error {
		if err := a.http.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error (public api server): %w", err)
		}
		return nil
	})

	errs.Go(func() error {
		if err := a.httpInternal.ListenAndServe(); err != nil {
			return fmt.Errorf("listen and serve error (internal api server): %w", err)
		}
		return nil
	})

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully")

	// Perform application shutdown with a maximum timeout of 5 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.http.Shutdown(timeoutCtx); err != nil {
		log.Println(err.Error())
	}

	return nil
}

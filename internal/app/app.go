package app

import (
	"context"
	"errors"
	"net/http"
	"project-example/internal/modules/customer"
	customermemory "project-example/internal/modules/customer/infrastructure/memory"
	"project-example/internal/modules/order"
	orderdomain "project-example/internal/modules/order/domain"
	ordermemory "project-example/internal/modules/order/infrastructure/memory"
	orderpostgres "project-example/internal/modules/order/infrastructure/postgres"
	"project-example/internal/platform/config"
	"project-example/internal/platform/database"
	"project-example/internal/platform/httpserver"
	"project-example/internal/shared/httpx"
	"project-example/internal/shared/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	logger        *zap.Logger
	server        *httpserver.Server
	closePostgres func() error
}

func New(cfg config.Config, logger *zap.Logger) (*App, error) {
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	postgresDB, closePostgres, err := openPostgresDB(cfg, logger)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.AccessLog(cfg.Logging),
		gin.Recovery(),
	)
	router.GET("/healthz", func(c *gin.Context) {

		httpx.JSON(c, http.StatusOK, gin.H{
			"app":    cfg.AppName,
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")

	customerRepository := customermemory.NewRepository(customermemory.SeedCustomers())
	customerModule := customer.NewModule(customerRepository)
	customerModule.RegisterRoutes(api)

	orderRepository := newOrderRepository(postgresDB, logger)
	orderModule := order.NewModule(orderRepository)
	orderModule.RegisterRoutes(api)

	return &App{
		logger:        logger,
		server:        httpserver.New(cfg.HTTP.Addr(), router),
		closePostgres: closePostgres,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.closeResources()

	errCh := make(chan error, 1)

	a.logger.Info("starting http server", zap.String("addr", a.server.Addr()))
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		a.logger.Info("shutting down http server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return a.server.Shutdown(shutdownCtx)
	}
}

func openPostgresDB(cfg config.Config, logger *zap.Logger) (*gorm.DB, func() error, error) {
	if !cfg.Database.Enabled() {
		logger.Info("shared postgres connection is disabled")
		return nil, nil, nil
	}

	db, err := database.OpenPostgres(cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	logger.Info("shared postgres connection initialized")
	return db, sqlDB.Close, nil
}

func newOrderRepository(db *gorm.DB, logger *zap.Logger) orderdomain.Repository {
	if db == nil {
		logger.Info("order repository is using memory storage")
		return ordermemory.NewRepository(ordermemory.SeedOrders())
	}

	logger.Info("order repository is using postgres storage")
	return orderpostgres.NewRepository(db)
}

func (a *App) closeResources() {
	if a.closePostgres == nil {
		return
	}

	if err := a.closePostgres(); err != nil {
		a.logger.Error("failed to close postgres connection", zap.Error(err))
	}
}

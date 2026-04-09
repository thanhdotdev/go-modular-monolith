package app

import (
	"context"
	"errors"
	"net/http"
	"project-example/internal/modules/customer"
	customermemory "project-example/internal/modules/customer/infrastructure/memory"
	discount "project-example/internal/modules/discount"
	discountapplication "project-example/internal/modules/discount/application"
	discountmemory "project-example/internal/modules/discount/infrastructure/memory"
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
	setGinMode(cfg.GinMode)

	postgresDB, closePostgres, err := openPostgresDB(cfg, logger)
	if err != nil {
		return nil, err
	}

	router := newRouter(cfg)
	registerSharedRoutes(router, cfg)
	registerModules(router.Group("/api/v1"), postgresDB, logger)

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

func setGinMode(mode string) {
	if mode == "" {
		return
	}

	gin.SetMode(mode)
}

func newRouter(cfg config.Config) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.AccessLog(cfg.Logging),
		gin.Recovery(),
	)

	return router
}

func registerSharedRoutes(router *gin.Engine, cfg config.Config) {
	router.GET("/healthz", func(c *gin.Context) {
		httpx.JSON(c, http.StatusOK, gin.H{
			"app":    cfg.AppName,
			"status": "ok",
		})
	})
}

func registerModules(api *gin.RouterGroup, postgresDB *gorm.DB, logger *zap.Logger) {
	discounts := registerDiscountModule(api)
	registerCustomerModule(api)
	registerOrderModule(api, postgresDB, logger, discounts)
}

func registerDiscountModule(api *gin.RouterGroup) discountapplication.UseCase {
	module := discount.NewModule(discount.Dependencies{
		Repository: discountmemory.NewRepository(discountmemory.SeedDiscounts()),
	})
	module.RegisterRoutes(api)

	return module.UseCase()
}

func registerCustomerModule(api *gin.RouterGroup) {
	module := customer.NewModule(customer.Dependencies{
		Repository: customermemory.NewRepository(customermemory.SeedCustomers()),
	})
	module.RegisterRoutes(api)
}

func registerOrderModule(api *gin.RouterGroup, postgresDB *gorm.DB, logger *zap.Logger, discounts discountapplication.UseCase) {
	module := order.NewModule(order.Dependencies{
		Repository: newOrderRepository(postgresDB, logger),
		Discounts:  discounts,
	})
	module.RegisterRoutes(api)
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

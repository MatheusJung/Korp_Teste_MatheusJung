// @title           Faturamento API
// @version         1.0
// @description     Serviço de emissão de notas fiscais com Outbox Pattern e idempotência.
// @host            localhost:5002
// @BasePath        /api

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nf-system/servico-faturamento/internal/application/usecases"
	"github.com/nf-system/servico-faturamento/internal/handler"
	"github.com/nf-system/servico-faturamento/internal/handler/middleware"
	estoqueclient "github.com/nf-system/servico-faturamento/internal/infrastructure/http/estoque"
	pdfinfra "github.com/nf-system/servico-faturamento/internal/infrastructure/pdf"
	"github.com/nf-system/servico-faturamento/internal/infrastructure/persistence/migrations"
	"github.com/nf-system/servico-faturamento/internal/infrastructure/persistence/repositories"
	storageinfra "github.com/nf-system/servico-faturamento/internal/infrastructure/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/nf-system/servico-faturamento/docs"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbURL := getEnv("DATABASE_URL", "postgres://nf:nf@localhost:5432/faturamento?sslmode=disable")
	estoqueURL := getEnv("ESTOQUE_URL", "http://localhost:5001")
	port := getEnv("PORT", "8080")
	maxTentativas := 5
	pollInterval := 5 * time.Second

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Fatal("falha ao conectar ao PostgreSQL", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("PostgreSQL não respondeu ao ping", zap.Error(err))
	}
	logger.Info("conectado ao PostgreSQL")

	if err := migrations.Run(ctx, pool); err != nil {
		logger.Fatal("falha ao aplicar migrations", zap.Error(err))
	}
	logger.Info("migrations aplicadas")

	notaRepo := repositories.NewNotaFiscalRepository(pool)
	outboxRepo := repositories.NewOutboxRepository(pool)
	idempotencyRepo := repositories.NewIdempotencyRepository(pool)
	uow := repositories.NewUnitOfWork(pool)

	estoqueClient := estoqueclient.NewHTTPClient(estoqueURL, 10*time.Second)

	criarUC := usecases.NewCriarNotaUseCase(uow, notaRepo, logger)
	obterUC := usecases.NewObterNotaUseCase(notaRepo)
	listarUC := usecases.NewListarNotasUseCase(notaRepo)
	listarPaginadoUC := usecases.NewListarNotasPaginadoUseCase(notaRepo)
	imprimirUC := usecases.NewImprimirNotaUseCase(uow, notaRepo, idempotencyRepo, maxTentativas, logger)
	pdfService := pdfinfra.NewFPDFService()
	pdfStorage := storageinfra.NewLocalPDFStorage("/app/storage/notas")

	worker := usecases.NewOutboxWorker(
		outboxRepo, notaRepo, estoqueClient, uow, pdfService, pdfStorage,
		pollInterval, maxTentativas, logger,
	)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	defer cancelWorker()
	go worker.Start(workerCtx)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))

	r.GET("/health", healthHandler(pool, estoqueClient))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	notaHandler := handler.NewNotaFiscalHandler(
	criarUC,
	obterUC,
	listarUC,
	listarPaginadoUC,
	imprimirUC,
	pdfStorage,
)
	notaHandler.RegisterRoutes(r)

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info("serviço de faturamento iniciado", zap.String("porta", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("falha ao iniciar servidor", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("encerrando serviço...")
	cancelWorker()

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("erro no shutdown", zap.Error(err))
	}
	logger.Info("serviço encerrado")
}

// healthHandler godoc
// @Summary     Health check
// @Description Retorna status do serviço, banco de dados e conexão com o estoque
// @Tags        infra
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     503 {object} map[string]interface{}
// @Router      /health [get]
func healthHandler(pool *pgxpool.Pool, estoqueClient interface {
	HealthCheck(context.Context) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		estoqueStatus := "ok"
		if err := estoqueClient.HealthCheck(c.Request.Context()); err != nil {
			estoqueStatus = "down"
		}

		dbStatus := "ok"
		if err := pool.Ping(c.Request.Context()); err != nil {
			dbStatus = "down"
		}

		status := "ok"
		if dbStatus == "down" {
			status = "down"
		} else if estoqueStatus == "down" {
			status = "degraded"
		}

		httpStatus := http.StatusOK
		if status == "down" {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"status": status,
			"checks": gin.H{
				"database":        dbStatus,
				"estoque_service": estoqueStatus,
			},
			"timestamp": time.Now().UTC(),
		})
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

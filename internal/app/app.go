package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/homework3/notification/internal/config"
	"github.com/homework3/notification/internal/kafka"
	"github.com/homework3/notification/internal/metrics"
	"github.com/homework3/notification/internal/repository"
	"github.com/homework3/notification/internal/stmp_sender"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type App struct {
	repo repository.Repository
}

func New(repo repository.Repository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) Start(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := kafka.StartProcessMessages(ctx, a.repo, stmp_sender.NewSender(&cfg.Smtp), &cfg.Kafka); err != nil {
			log.Error().Err(err).Msg("Failed to start kafka consumer")
			cancel()
		}
	}()

	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	metricsServer := metrics.CreateMetricsServer(metricsAddr, cfg)

	go func() {
		log.Info().Msgf("Metrics http_server is running on %s", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running metrics http_server")
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info().Msgf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Info().Msgf("ctx.Done: %v", done)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("metricsServer.Shutdown")
	} else {
		log.Info().Msg("metricsServer shut down correctly")
	}

	return nil
}

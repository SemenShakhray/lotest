package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"lotest/internal/config"
	"lotest/internal/handlers"
	"lotest/internal/logger"
	"lotest/internal/repository"
	"lotest/internal/router"
	"lotest/internal/service"
)

const buferSize = 50

func main() {
	cfg := config.MustLoad()

	logCh := make(chan string, buferSize)

	ctxLog, cancelLog := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.InitLogger(ctxLog, logCh)
	}()

	repo := repository.NewRepo()

	s := service.NewService(repo, logCh)

	h := handlers.NewHandler(s, logCh)

	r := router.NewRouter(h)

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logCh <- logger.FormatLog(logger.Info, fmt.Sprintf("server started on %s:%s", cfg.Server.Host, cfg.Server.Port))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logCh <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to start server - %s", err.Error()))

			os.Exit(1)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-sigint

	ctxSrv, cancelSrv := context.WithTimeout(context.Background(), 10*time.Second)

	if err := srv.Shutdown(ctxSrv); err != nil {
		logCh <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to shutdown server - %s", err.Error()))
	}

	logCh <- logger.FormatLog(logger.Info, fmt.Sprintf("server shutdown with signal %s", sig.String()))

	cancelSrv()
	cancelLog()
	wg.Wait()
}

package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/darzox/test_task_iq_progress/app"
	"github.com/darzox/test_task_iq_progress/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

const (
	appName = "api-server"
)

var embedMigrations embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	loggerEntry := logger.WithField("app.name", appName)

	cfg, err := config.NewConfig()
	if err != nil {
		loggerEntry.Panicf("failed to create config: %v", err)
	}

	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDb, cfg.PostgresSslmode)
	dbPool, err := pgxpool.New(context.Background(), pgConnStr)
	if err != nil {
		loggerEntry.Panic("failed to create pool connections")
	}
	if err := dbPool.Ping(context.TODO()); err != nil {
		loggerEntry.Panicf("failed to ping db: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		loggerEntry.Panicf("failed to set dialect for goose: %v", err)
	}
	goose.SetBaseFS(embedMigrations)

	// db := stdlib.OpenDBFromPool(dbPool)
	// _, err = fs.Stat(embedMigrations, "migrations")
	// fmt.Println(err, "tut")
	// if err := goose.Up(db, "migrations"); err != nil {
	// 	loggerEntry.Panicf("failed to up migrations: %v", err)
	// }

	server := http.Server{}
	server.Addr = fmt.Sprintf(":%d", cfg.HttpPort)
	server.MaxHeaderBytes = 1 << 20

	if err := app.Run(&server, dbPool, loggerEntry); err != nil {
		loggerEntry.Panicf("failed to run app")
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			loggerEntry.Fatalf("failed to run server: %v", err)
		}
	}()
	loggerEntry.Infof("listening on %s", server.Addr)

	<-ctx.Done()
	loggerEntry.Info("shuting down server and closing connections")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorf("shutdown server err: %v", err)
	}

	loggerEntry.Info("closing db connections")
	dbPool.Close()

	loggerEntry.Info("application stopped")
}

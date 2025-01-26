package app

import (
	"net/http"

	"github.com/darzox/test_task_iq_progress/internal/handler"
	"github.com/darzox/test_task_iq_progress/internal/repository"
	"github.com/darzox/test_task_iq_progress/internal/routes"
	"github.com/darzox/test_task_iq_progress/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func Run(server *http.Server, dbPool *pgxpool.Pool, logger *logrus.Entry) error {
	repo := repository.NewRepo(dbPool, logger)
	service := service.NewService(repo, logger)
	handler := handler.NewHandler(service, logger)

	router := gin.Default()
	routes.RegisterRoutes(handler, router)
	server.Handler = router

	return nil
}

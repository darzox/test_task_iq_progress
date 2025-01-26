package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/darzox/test_task_iq_progress/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type handler struct {
	service Service
	logger  *logrus.Entry
}

func NewHandler(service Service, logger *logrus.Entry) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type Service interface {
	Deposit(ctx context.Context, userId int64, amount float64) error
	Transfer(ctx context.Context, fromUserId int64, toUserId int64, amount float64) error
	GetLast10Transactions(ctx context.Context, userId int64) ([]models.Transaction, error)
}

// Deposit - deposit balance
// @Description deposit to user balance
// @Tags User
// @Accept json
// @Param DepositRequest body handler.DepositRequest true "deposit info"
// @Success 200
// @Failure 400 "BadRequest"
// @Failure 500 "UnknownError"
// @Router /deposit [post]
func (h *handler) Deposit(ctx *gin.Context) {
	var req DepositRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := h.service.Deposit(ctx.Request.Context(), req.UserId, req.Amount)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// Transfer - from one balance to another
// @Description transfer money from one balance to another
// @Tags User
// @Accept json
// @Param TransferRequest body handler.TransferRequest true "transfer info"
// @Success 200
// @Failure 400 "BadRequest"
// @Failure 500 "UnknownError"
// @Router /transfer [post]
func (h *handler) Transfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := h.service.Transfer(ctx.Request.Context(), req.FromUserId, req.ToUserId, req.Amount)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// GetLast10Transactions - get last transactions of a user
// @Description get last transactions of a user
// @Tags User
// @Accept json
// @Param user_id query string true "User id"
// @Success 200 {array} models.Transaction
// @Failure 400 "BadRequest"
// @Failure 500 "UnknownError"
// @Router /transactions [get]
func (h *handler) GetLast10Transactions(ctx *gin.Context) {
	userIdStr := ctx.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	transactions, err := h.service.GetLast10Transactions(ctx.Request.Context(), userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, transactions)
}

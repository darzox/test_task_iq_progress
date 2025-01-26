package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/darzox/test_task_iq_progress/internal/models"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo   Repository
	logger *logrus.Entry
}

type Repository interface {
	Deposit(ctx context.Context, userId int64, amount float64, comment string) error
	Transfer(ctx context.Context, fromUserId int64, toUserId int64, amount float64, comment string) error
	GetUserBalance(ctx context.Context, userId int64) (float64, error)
	GetLast10Transactions(ctx context.Context, userId int64) ([]models.Transaction, error)
}

func NewService(repo Repository, logger *logrus.Entry) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Deposit(ctx context.Context, userId int64, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if err := validateUserId(userId); err != nil {
		return err
	}

	comment := fmt.Sprintf("someone deposit user_id=%d to amount=%f", userId, amount)

	return s.repo.Deposit(ctx, userId, amount, comment)
}

func (s *service) Transfer(ctx context.Context, fromUserId int64, toUserId int64, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if err := validateUserId(fromUserId); err != nil {
		return fmt.Errorf("from_user_id: %w", err)
	}
	if err := validateUserId(toUserId); err != nil {
		return fmt.Errorf("to_user_id: %w", err)
	}

	if fromUserId == toUserId {
		return errors.New("user_ids are equal")
	}

	userBalance, err := s.repo.GetUserBalance(ctx, fromUserId)
	if err != nil {
		return err
	}

	if userBalance-amount < 0 {
		return errors.New("balance cannot be below zero")
	}

	comment := fmt.Sprintf("user_id=%d transfer amount=%f to user_id=%d", fromUserId, amount, toUserId)

	return s.repo.Transfer(ctx, fromUserId, toUserId, amount, comment)
}

func (s *service) GetLast10Transactions(ctx context.Context, userId int64) ([]models.Transaction, error) {
	if err := validateUserId(userId); err != nil {
		return nil, err
	}

	return s.repo.GetLast10Transactions(ctx, userId)
}

func validateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount is negative or zero")
	}

	return nil
}

func validateUserId(userId int64) error {
	if userId <= 0 {
		return errors.New("user_id is zero")
	}

	return nil
}

package tests

import (
	"context"
	"testing"

	"github.com/darzox/test_task_iq_progress/internal/models"
	"github.com/darzox/test_task_iq_progress/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Deposit(ctx context.Context, userId int64, amount float64, comment string) error {
	args := m.Called(ctx, userId, amount, comment)
	return args.Error(0)
}

func (m *MockRepository) Transfer(ctx context.Context, fromUserId int64, toUserId int64, amount float64, comment string) error {
	args := m.Called(ctx, fromUserId, toUserId, amount, comment)
	return args.Error(0)
}

func (m *MockRepository) GetUserBalance(ctx context.Context, userId int64) (float64, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepository) GetLast10Transactions(ctx context.Context, userId int64) ([]models.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func TestDeposit(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	userId := int64(1)
	amount := 100.0
	comment := "someone deposit user_id=1 to amount=100.000000"

	mockRepo.On("Deposit", ctx, userId, amount, comment).Return(nil)

	err := svc.Deposit(ctx, userId, amount)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeposit_InvalidAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	userId := int64(1)
	amount := -50.0

	err := svc.Deposit(ctx, userId, amount)
	assert.Error(t, err)
	assert.Equal(t, "amount is negative or zero", err.Error())
}

func TestTransfer_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	fromUserId := int64(1)
	toUserId := int64(2)
	amount := 50.0
	comment := "user_id=1 transfer amount=50.000000 to user_id=2"

	mockRepo.On("GetUserBalance", ctx, fromUserId).Return(100.0, nil)
	mockRepo.On("Transfer", ctx, fromUserId, toUserId, amount, comment).Return(nil)

	err := svc.Transfer(ctx, fromUserId, toUserId, amount)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	fromUserId := int64(1)
	toUserId := int64(2)
	amount := 150.0

	mockRepo.On("GetUserBalance", ctx, fromUserId).Return(100.0, nil)

	err := svc.Transfer(ctx, fromUserId, toUserId, amount)
	assert.Error(t, err)
	assert.Equal(t, "balance cannot be below zero", err.Error())
}

func TestGetLast10Transactions(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	userId := int64(1)
	transactions := []models.Transaction{
		{Id: 1, Amount: 100},
		{Id: 2, Amount: 200},
	}

	mockRepo.On("GetLast10Transactions", ctx, userId).Return(transactions, nil)

	result, err := svc.GetLast10Transactions(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_InvalidUserIDs(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	err := svc.Transfer(ctx, -1, 2, 50.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "from_user_id")

	err = svc.Transfer(ctx, 1, -2, 50.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "to_user_id")
}

func TestTransfer_SameUserID(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	err := svc.Transfer(ctx, 1, 1, 50.0)
	assert.Error(t, err)
	assert.Equal(t, "user_ids are equal", err.Error())
}

func TestGetLast10Transactions_InvalidUserID(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	userId := int64(0)

	_, err := svc.GetLast10Transactions(ctx, userId)
	assert.Error(t, err)
	assert.Equal(t, "user_id is zero", err.Error())
}

func TestDeposit_InvalidUserId(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.NewEntry(logrus.New())
	svc := service.NewService(mockRepo, logger)

	ctx := context.Background()
	userId := int64(-1)
	amount := 100.0

	err := svc.Deposit(ctx, userId, amount)
	assert.Error(t, err)
	assert.Equal(t, "user_id is zero", err.Error())
}

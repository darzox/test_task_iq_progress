package handler

type DepositRequest struct {
	UserId int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	FromUserId int64   `json:"from_user_id"`
	ToUserId   int64   `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}


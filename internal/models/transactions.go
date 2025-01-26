package models

type Transaction struct {
	Id       int64   `json:"id" db:"id"`
	Amount   float64 `json:"amount" db:"amount"`
	Comment  string  `json:"comment" db:"comment"`
	TypeName string  `json:"type_name" db:"type_name"`
}

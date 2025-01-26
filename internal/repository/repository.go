package repository

import (
	"context"

	"github.com/darzox/test_task_iq_progress/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type repository struct {
	db     *pgxpool.Pool
	logger *logrus.Entry
}

func NewRepo(dbPool *pgxpool.Pool, logger *logrus.Entry) *repository {
	return &repository{
		db:     dbPool,
		logger: logger,
	}
}

func (r *repository) Deposit(ctx context.Context, userId int64, amount float64, comment string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "update users set balance = balance + $1 where id = $2", amount, userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert into transactions(transaction_type_id, user_id, amount, comment) values ((select id from transaction_types where name = 'deposit'), $1, $2, $3)", userId, amount, comment)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Transfer(ctx context.Context, fromUserId int64, toUserId int64, amount float64, comment string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "update users set balance = balance - $1 where id = $2", amount, fromUserId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "update users set balance = balance + $1 where id = $2", amount, toUserId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert into transactions(transaction_type_id, user_id, amount, comment) values ((select id from transaction_types where name = 'transfer'), $1, $2, $3)", fromUserId, amount, comment)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUserBalance(ctx context.Context, userId int64) (float64, error) {
	var result float64
	row := r.db.QueryRow(ctx, "select blance from users where id = $1", userId)
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, err
}

func (r *repository) GetLast10Transactions(ctx context.Context, userId int64) ([]models.Transaction, error) {
	var result []models.Transaction
	rows, err := r.db.Query(ctx, `
	select
		t.id as id,
		t.amount as amount,
		t.comment as comment,
		tt.name as type_name
	from transactions t
	left join transaction_types tt on t.transaction_type_id = tt.id
	where user_id = $1
	order by t.created_at desc
	limit 10`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tr models.Transaction
		if err := rows.Scan(&tr.Id, &tr.Amount, &tr.Comment, &tr.TypeName); err != nil {
			return nil, err
		}
		result = append(result, tr)
	}

	return result, nil
}

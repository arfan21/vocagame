package transactionrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/pkg/constant"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db    dbpostgres.Queryer
	rawDb dbpostgres.Raw
}

func New(raw dbpostgres.Raw, queryer dbpostgres.Queryer) *Repository {
	return &Repository{
		db:    queryer,
		rawDb: raw,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.rawDb.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Transaction) (id uuid.UUID, err error) {
	query := `
		INSERT INTO transactions (user_id, transaction_type_id, status, total_amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = r.db.QueryRow(ctx, query,
		data.UserID,
		data.TransactionTypeID,
		data.Status,
		data.TotalAmount,
	).Scan(&data.ID)

	if err != nil {
		err = fmt.Errorf("transaction.repository.Create: failed to create transaction: %w", err)
		return
	}

	return data.ID, nil
}

func (r Repository) CreateDetail(ctx context.Context, data []entity.TransactionDetail) (err error) {
	columns := []string{"transaction_id", "product_id", "qty"}

	rows := make([][]interface{}, len(data))
	for i, item := range data {
		rows[i] = []interface{}{item.TransactionID, item.ProductID, item.Qty}
	}

	rowsAffected, err := r.db.CopyFrom(ctx,
		pgx.Identifier{entity.TransactionDetail{}.TableName()},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		err = fmt.Errorf("transaction.repository.CreateDetail: failed to create transaction detail: %w", err)
		return
	}

	if rowsAffected != int64(len(data)) {
		err = fmt.Errorf("transaction.repository.CreateDetail: failed to create transaction detail: %w", constant.ErrTxDetailInsertedNotEqual)
		return
	}

	return
}

func (r Repository) GetHistoryWalletByUserID(ctx context.Context, userID uuid.UUID) (res []entity.Transaction, err error) {
	query := `
		SELECT 
			t.id, 
			t.user_id, 
			t.transaction_type_id,
			tt.name AS transaction_type_name, 
			t.status, 
			t.total_amount, 
			t.created_at, 
			t.updated_at
		FROM transactions t
		JOIN transaction_types tt ON t.transaction_type_id = tt.id
		WHERE t.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		err = fmt.Errorf("transaction.repository.GetHistoryWalletByUserID: failed to get history wallet by user id: %w", err)
		return
	}

	for rows.Next() {
		var data entity.Transaction

		err = rows.Scan(
			&data.ID,
			&data.UserID,
			&data.TransactionTypeID,
			&data.TransactionType.Name,
			&data.Status,
			&data.TotalAmount,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			err = fmt.Errorf("transaction.repository.GetHistoryWalletByUserID: failed to scan data: %w", err)
			return
		}

		res = append(res, data)
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (res entity.Transaction, err error) {
	query := `
		SELECT 
			t.id, 
			t.user_id, 
			t.transaction_type_id,
			tt.name AS transaction_type_name, 
			t.status, 
			t.total_amount, 
			t.created_at, 
			t.updated_at,
			td.id AS transaction_detail_id,
			td.product_id,
			td.qty,
			p.name AS product_name,
			p.price AS product_price
		FROM transactions t
		LEFT JOIN transaction_types tt ON t.transaction_type_id = tt.id
		LEFT JOIN transaction_details td ON t.id = td.transaction_id
		LEFT JOIN products p ON td.product_id = p.id
		WHERE t.id = $1 AND t.user_id = $2
	`

	rows, err := r.db.Query(ctx, query, id, userID)
	if err != nil {
		err = fmt.Errorf("transaction.repository.GetByID: failed to get transaction by id: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var detail entity.TransactionDetail

		err = rows.Scan(
			&res.ID,
			&res.UserID,
			&res.TransactionTypeID,
			&res.TransactionType.Name,
			&res.Status,
			&res.TotalAmount,
			&res.CreatedAt,
			&res.UpdatedAt,
			&detail.ID,
			&detail.ProductID,
			&detail.Qty,
			&detail.Product.Name,
			&detail.Product.Price,
		)

		if err != nil {
			err = fmt.Errorf("transaction.repository.GetByID: failed to scan data: %w", err)
			return
		}

		if detail.ID.Valid {
			res.TransactionDetail = append(res.TransactionDetail, detail)
		}

	}

	if rows.Err() != nil {
		err = fmt.Errorf("transaction.repository.GetByID: failed after scan data: %w", err)
		return
	}

	if rows.CommandTag().RowsAffected() == 0 {
		err = constant.ErrTransactionNotFound
		return
	}

	return
}

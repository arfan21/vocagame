package productrepo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func (r Repository) Create(ctx context.Context, data entity.Product) (err error) {
	query := `
		INSERT INTO products (user_id, name, description, stok, price)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query,
		data.UserID,
		data.Name,
		data.Description,
		data.Stok,
		data.Price,
	)

	if err != nil {
		err = fmt.Errorf("product.repository.Create: failed to create product: %w", err)
		return
	}

	return
}

func (r Repository) queryRowsProductWithFilter(ctx context.Context, query string, filter entity.ListProductFilter) (rows pgx.Rows, err error) {
	var filterArgs []any
	var whereQuery string

	if len(filter.Name) != 0 {
		filterName := "%" + strings.ToLower(filter.Name) + "%"
		filterArgs = append(filterArgs, filterName)
		whereQuery += "LOWER(p.name) LIKE $" + strconv.Itoa(len(filterArgs)) + " AND "
	}

	if filter.UserID.Valid {
		filterArgs = append(filterArgs, filter.UserID.UUID)
		whereQuery += "p.user_id = $" + strconv.Itoa(len(filterArgs)) + " AND "
	}

	if filter.ID.Valid {
		filterArgs = append(filterArgs, filter.ID.UUID)
		whereQuery += "p.id = $" + strconv.Itoa(len(filterArgs)) + " AND "
	}

	// if filterArgsLen  > 0, add WHERE statement and remove last AND
	if filterArgsLen := len(filterArgs); filterArgsLen > 0 {
		whereQuery = "WHERE " + whereQuery[:len(whereQuery)-len(" AND ")] + " "
	}

	query += whereQuery

	if !filter.DisableOffset {
		filterArgs = append(filterArgs, filter.Limit)
		query += "LIMIT $" + strconv.Itoa(len(filterArgs)) + " "

		offset := (filter.Page - 1) * filter.Limit
		filterArgs = append(filterArgs, offset)
		query += "OFFSET $" + strconv.Itoa(len(filterArgs)) + " "
	}

	return r.db.Query(ctx, query, filterArgs...)
}

func (r Repository) GetProducts(ctx context.Context, filter entity.ListProductFilter) (result []entity.Product, err error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.stok,
			p.price,
			p.description,
			u.id AS owner_id,
			u.fullname AS owner_name
		FROM
			products p
			JOIN users u ON u.id = p.user_id
	`

	rows, err := r.queryRowsProductWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("product.repository.GetProducts: failed to get products: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var product entity.Product

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Stok,
			&product.Price,
			&product.Description,
			&product.User.ID,
			&product.User.Fullname,
		)

		if err != nil {
			err = fmt.Errorf("product.repository.GetProducts: failed to scan product: %w", err)
			return
		}

		result = append(result, product)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("product.repository.GetProducts: failed after scan products: %w", err)
		return
	}

	return
}

func (r Repository) GetTotalProduct(ctx context.Context, filter entity.ListProductFilter) (result int, err error) {
	query := `
		SELECT
			COUNT(p.id)
		FROM
			products p
			JOIN users u ON u.id = p.user_id
	`
	filter.DisableOffset = true
	rows, err := r.queryRowsProductWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("product.repository.GetTotalProduct: failed to get total product: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&result)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("product.repository.GetProducts: failed after scan total product: %w", err)
		return
	}

	return
}

func (r Repository) Update(ctx context.Context, data entity.Product) (err error) {
	query := `
		UPDATE products
		SET
			name = $1,
			description = $2,
			stok = $3,
			price = $4
		WHERE
			id = $5
	`

	_, err = r.db.Exec(ctx, query,
		data.Name,
		data.Description,
		data.Stok,
		data.Price,
		data.ID,
	)

	if err != nil {
		err = fmt.Errorf("product.repository.Update: failed to update product: %w", err)
		return
	}

	return
}

func (r Repository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("product.repository.Delete: failed to delete product: %w", err)
		return
	}

	return
}

func (r Repository) ReduceStok(ctx context.Context, id uuid.UUID, reduceBy int) (err error) {
	query := `
		UPDATE products
		SET stok = stok - $1
		WHERE id = $2 AND (stok - $1) >= 0
	`

	cmd, err := r.db.Exec(ctx, query, reduceBy, id)
	if err != nil {
		err = fmt.Errorf("product.repository.BatchUpdateStok: failed to reduce stok: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = fmt.Errorf("product.repository.BatchUpdateStok: nothing updated: %w", constant.ErrProductNotFoundOrStok)
		return err
	}

	return
}

func (r Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) (result map[uuid.UUID]entity.Product, err error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.stok,
			p.price,
			p.user_id
		FROM
			products p
		WHERE p.id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		err = fmt.Errorf("product.repository.GetByIDs: failed to get products by ids: %w", err)
		return
	}

	defer rows.Close()

	result = make(map[uuid.UUID]entity.Product)

	for rows.Next() {
		var product entity.Product

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Stok,
			&product.Price,
			&product.UserID,
		)

		if err != nil {
			err = fmt.Errorf("product.repository.GetByIDs: failed to scan product: %w", err)
			return
		}

		result[product.ID] = product
	}

	if rows.Err() != nil {
		err = fmt.Errorf("product.repository.GetByIDs: failed after scan products: %w", err)
		return
	}

	return
}

package productrepo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/arfan21/vocagame/internal/entity"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db    dbpostgres.Queryer
	rawDb *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:    db,
		rawDb: db,
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

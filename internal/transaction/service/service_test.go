package transactionsvc

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	walletrepo "github.com/arfan21/vocagame/internal/wallet/repository"
	walletsvc "github.com/arfan21/vocagame/internal/wallet/service"
	"github.com/arfan21/vocagame/migration"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var db *pgxpool.Pool
var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

func initDocker(t *testing.T) (*dockertest.Pool, *dockertest.Resource) {
	ctx := context.Background()
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Errorf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.Errorf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.0-alpine3.18",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=postgres-test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	},
	)
	if err != nil {
		t.Errorf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		connString := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres-test?sslmode=disable", resource.GetPort("5432/tcp"))
		db, err = pgxpool.New(ctx, connString)
		if err != nil {
			return err
		}
		err = db.Ping(ctx)
		return err
	}); err != nil {
		t.Errorf("Could not connect to database: %s", err)
		db = nil
	}

	resource.Expire(30)
	fmt.Println("Database connected")

	dbSql := stdlib.OpenDBFromPool(db)

	fmt.Println("migrate main")
	mig, err := migration.New(dbSql)
	if err != nil {
		t.Errorf("Could not get migrations: %s", err)
	}

	mig.Up(ctx)

	return pool, resource
}

func initDep(t *testing.T) (svc *Service) {
	dockerPool, dockerResource = initDocker(t)

	walletRepo := walletrepo.New(db, db)
	walletSvc := walletsvc.New(walletRepo)

	transactionRepo := transactionrepo.New(db, db)
	svc = New(transactionRepo, walletSvc)

	return
}

func initUser(t *testing.T) (id uuid.UUID) {
	query := `
		INSERT INTO users (fullname, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.QueryRow(context.Background(), query, "test", "test@email.com", "password").
		Scan(&id)

	assert.NoError(t, err)

	return
}

var (
	initialBalance = decimal.NewFromInt(5000)
)

func initWallet(t *testing.T, userID uuid.UUID) {
	query := `
		INSERT INTO wallets (user_id, balance)
		VALUES ($1, $2)
	`

	_, err := db.Exec(context.Background(), query, userID, initialBalance)
	assert.NoError(t, err)
}

func truncateAllTable(t *testing.T) {
	tables := []string{
		"transactions",
		"wallets",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s CASCADE", table))
		assert.NoError(t, err)
	}
}

func getWalletBalance(t *testing.T, userID uuid.UUID) (balance decimal.Decimal) {
	query := `
		SELECT balance
		FROM wallets
		WHERE user_id = $1
	`

	err := db.QueryRow(context.Background(), query, userID).Scan(&balance)
	assert.NoError(t, err)

	return
}

func getTransactionTotalAmount(t *testing.T, id uuid.UUID) (totalAmount decimal.Decimal) {
	query := `
		SELECT total_amount
		FROM transactions
		WHERE id = $1
	`

	err := db.QueryRow(context.Background(), query, id).Scan(&totalAmount)
	assert.NoError(t, err)

	return
}

func TestCreateDepositTransactionConcurrent(t *testing.T) {
	svc := initDep(t)

	assert.NotNil(t, db)

	req := model.CreateDepositTransactionRequest{
		Amount: decimal.NewFromInt(50000),
	}
	truncateAllTable(t)
	userID := initUser(t)
	initWallet(t, userID)

	req.UserID = userID

	totalConcurrent := 10

	wg := &sync.WaitGroup{}
	// concurrent
	for i := 0; i < totalConcurrent; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()

			id, err := svc.CreateDepositTransaction(context.Background(), req)
			assert.NoError(t, err)
			assert.NotEqual(t, id.TransactionID, uuid.Nil)
			fmt.Println("concurrent", i, "done with id", id.TransactionID)
		}(wg, i)
	}

	wg.Wait()
	fmt.Println("all concurrent done")
	balance := getWalletBalance(t, userID)
	fmt.Println(balance)
	assert.True(t, balance.Equal(req.Amount.Mul(decimal.NewFromInt(int64(totalConcurrent))).Add(initialBalance)))
}

func initPgMock(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	return mock
}

func initDepMock(db pgxmock.PgxPoolIface) (svc *Service) {

	walletRepo := walletrepo.New(db, db)
	walletSvc := walletsvc.New(walletRepo)

	transactionRepo := transactionrepo.New(db, db)
	svc = New(transactionRepo, walletSvc)

	return
}

func TestCreateDepositTransactionSuccess(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	transactionID := uuid.New()

	req := model.CreateDepositTransactionRequest{
		Amount: decimal.NewFromInt(50000),
		UserID: userID,
	}

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Add(req.Amount), walletID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// insert transaction
	dbMock.ExpectQuery("INSERT INTO transactions (.+) VALUES (.+) RETURNING id").
		WithArgs(userID, constant.TransactionTypeDepositID, entity.TransactionStatusCompleted, req.Amount).
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).AddRow(transactionID),
		)

	dbMock.ExpectCommit()

	id, err := svc.CreateDepositTransaction(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEqual(t, id.TransactionID, transactionID)

}

func TestCreateDepositTransactionFailedWalletNotFound(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	req := model.CreateDepositTransactionRequest{
		Amount: decimal.NewFromInt(50000),
		UserID: userID,
	}

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	dbMock.ExpectRollback()

	id, err := svc.CreateDepositTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrWalletNotFound)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateDepositTransactionFailedUpdateBalance(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	req := model.CreateDepositTransactionRequest{
		Amount: decimal.NewFromInt(50000),
		UserID: userID,
	}

	errUnexpected := fmt.Errorf("unexpected error")

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Add(req.Amount), walletID).
		WillReturnError(errUnexpected)

	dbMock.ExpectRollback()

	id, err := svc.CreateDepositTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errUnexpected)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateDepositTransactionFailedInsertTransaction(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	req := model.CreateDepositTransactionRequest{
		Amount: decimal.NewFromInt(50000),
		UserID: userID,
	}

	errUnexpected := fmt.Errorf("unexpected error")

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Add(req.Amount), walletID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// insert transaction
	dbMock.ExpectQuery("INSERT INTO transactions (.+) VALUES (.+) RETURNING id").
		WithArgs(userID, constant.TransactionTypeDepositID, entity.TransactionStatusCompleted, req.Amount).
		WillReturnError(errUnexpected)

	dbMock.ExpectRollback()

	id, err := svc.CreateDepositTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errUnexpected)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateWithdrawTransactionConcurrent(t *testing.T) {
	svc := initDep(t)

	assert.NotNil(t, db)

	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(5000),
	}
	truncateAllTable(t)
	userID := initUser(t)
	initWallet(t, userID)

	req.UserID = userID

	totalConcurrent := 10
	arrError := make([]error, totalConcurrent)
	arrId := make([]string, totalConcurrent)

	wg := &sync.WaitGroup{}
	// concurrent
	for i := 0; i < totalConcurrent; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()

			id, err := svc.CreateWithdrawTransaction(context.Background(), req)
			arrError[i] = err
			arrId[i] = id.TransactionID
		}(wg, i)
	}

	wg.Wait()
	fmt.Println("all concurrent done")

	for i := 0; i < totalConcurrent; i++ {
		if arrError[i] != nil {
			assert.ErrorIs(t, arrError[i], constant.ErrInsufficientBalance)
		} else {
			assert.NotEqual(t, arrId[i], "")
		}
	}

	balance := getWalletBalance(t, userID)
	fmt.Println("balance ->  ", balance)
	assert.True(t, balance.Equal(initialBalance.Sub(req.Amount)))
}

func TestCreateWithdrawTransactionSuccess(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	transactionID := uuid.New()

	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(3000),
		UserID: userID,
	}

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Sub(req.Amount), walletID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// insert transaction
	dbMock.ExpectQuery("INSERT INTO transactions (.+) VALUES (.+) RETURNING id").
		WithArgs(userID, constant.TransactionTypeWithdrawID, entity.TransactionStatusCompleted, req.Amount).
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).AddRow(transactionID),
		)

	dbMock.ExpectCommit()

	id, err := svc.CreateWithdrawTransaction(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEqual(t, id.TransactionID, transactionID)
}

func TestCreateWithdrawTransactionFailedWalletNotFound(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(3000),
		UserID: userID,
	}

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	dbMock.ExpectRollback()

	id, err := svc.CreateWithdrawTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrWalletNotFound)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateWithdrawTransactionFailedInsufficienBalance(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(6000),
		UserID: userID,
	}

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	dbMock.ExpectRollback()

	id, err := svc.CreateWithdrawTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrInsufficientBalance)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateWithdrawTransactionFailedUpdateBalance(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(3000),
		UserID: userID,
	}

	errUnexpected := fmt.Errorf("unexpected error")

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Sub(req.Amount), walletID).
		WillReturnError(errUnexpected)

	dbMock.ExpectRollback()

	id, err := svc.CreateWithdrawTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errUnexpected)
	assert.Equal(t, "", id.TransactionID)
}

func TestCreateWithdrawTransactionFailedInsertTransaction(t *testing.T) {
	dbMock := initPgMock(t)
	svc := initDepMock(dbMock)

	assert.NotNil(t, dbMock)

	userID := uuid.New()
	walletID := uuid.New()
	req := model.CreateWithdrawTransactionRequest{
		Amount: decimal.NewFromInt(3000),
		UserID: userID,
	}

	errUnexpected := fmt.Errorf("unexpected error")

	dbMock.ExpectBegin()
	// get wallet
	dbMock.ExpectQuery("SELECT (.+) FROM wallets (.+) FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, userID, initialBalance, nil, nil),
		)

	// update balance
	dbMock.ExpectExec("UPDATE wallets SET balance = (.+) WHERE id (.+)  ").
		WithArgs(initialBalance.Sub(req.Amount), walletID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// insert transaction
	dbMock.ExpectQuery("INSERT INTO transactions (.+) VALUES (.+) RETURNING id").
		WithArgs(userID, constant.TransactionTypeWithdrawID, entity.TransactionStatusCompleted, req.Amount).
		WillReturnError(errUnexpected)

	dbMock.ExpectRollback()

	id, err := svc.CreateWithdrawTransaction(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errUnexpected)
	assert.Equal(t, "", id.TransactionID)
}

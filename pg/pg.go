package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required
)

// Repository is the application's data layer functionality.
type Repository interface {
	// transfer queries
	GetTransferByTransaction(ctx context.Context, transaction string) (Transfer, error)
	ListTransfersBySender(ctx context.Context, sender string) ([]Transfer, error)
	ListTransfersByReceiver(ctx context.Context, receiver string) ([]Transfer, error)
	ListTransfersByTokenID(ctx context.Context, tokenID int32) ([]Transfer, error)
	ListUnreadTransfers(ctx context.Context) ([]Transfer, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error)
	ReadTransfer(ctx context.Context, transaction string) (*Transfer, error)
}

type repoSvc struct {
	*Queries
	db *sql.DB
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

func (r *repoSvc) CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error) {
	transfer := new(Transfer)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}
		transfer = &res
		return nil
	})
	return transfer, err
}

func (r *repoSvc) ReadTransfer(ctx context.Context, transactionArg string) (*Transfer, error) {
	transfer := new(Transfer)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.ReadTransfer(ctx, transactionArg)
		if err != nil {
			return err
		}
		transfer = &res
		return nil
	})
	return transfer, err
}

// NewRepository returns an implementation of the Repository interface.
func NewRepository(db *sql.DB) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

// Open opens a database specified by the data source name.
// Format: host=foo port=5432 user=bar password=baz dbname=qux sslmode=disable"
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

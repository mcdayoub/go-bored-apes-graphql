package gqlgen

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"strconv"

	"github.com/mcdayoub/go-bored-apes-graphql/pg"
)

type Resolver struct {
	Repository pg.Repository
}

func (r *mutationResolver) CreateTransfer(ctx context.Context, input TranferInput) (*pg.Transfer, error) {
	transfer, err := r.Repository.CreateTransfer(ctx, pg.CreateTransferParams{
		Transaction: input.Transaction,
		Sender:      input.Sender,
		Receiver:    input.Receiver,
		TokenID:     int32(input.TokenID),
		Read:        input.Read,
	})
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *mutationResolver) ReadTransfer(ctx context.Context, transaction string) (*pg.Transfer, error) {
	transfer, err := r.Repository.ReadTransfer(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *queryResolver) TransfersByTransaction(ctx context.Context, transaction string) ([]pg.Transfer, error) {
	transfers, err := r.Repository.ListTransfersByTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *queryResolver) TransfersBySender(ctx context.Context, sender string) ([]pg.Transfer, error) {
	transfers, err := r.Repository.ListTransfersBySender(ctx, sender)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *queryResolver) TransfersByReceiver(ctx context.Context, receiver string) ([]pg.Transfer, error) {
	transfers, err := r.Repository.ListTransfersByReceiver(ctx, receiver)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *queryResolver) TransfersByTokenID(ctx context.Context, tokenID string) ([]pg.Transfer, error) {
	t, err := strconv.Atoi(tokenID)
	if err != nil {
		return nil, err
	}
	transfers, err := r.Repository.ListTransfersByTokenID(ctx, int32(t))
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *queryResolver) UnreadTransfers(ctx context.Context) ([]pg.Transfer, error) {
	transfers, err := r.Repository.ListUnreadTransfers(ctx)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

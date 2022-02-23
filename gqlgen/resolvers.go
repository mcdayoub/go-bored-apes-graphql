package gqlgen

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

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
	panic("not implemented")
}

func (r *queryResolver) TransferByTransaction(ctx context.Context, transaction string) (*pg.Transfer, error) {
	transfer, err := r.Repository.GetTransferByTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *queryResolver) TransfersBySender(ctx context.Context, sender string) ([]pg.Transfer, error) {
	panic("not implemented")
}

func (r *queryResolver) TransfersByReceiver(ctx context.Context, receiver string) ([]pg.Transfer, error) {
	panic("not implemented")
}

func (r *queryResolver) TransfersByTokenID(ctx context.Context, tokenID string) ([]pg.Transfer, error) {
	panic("not implemented")
}

func (r *queryResolver) UnreadTransfers(ctx context.Context) ([]pg.Transfer, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

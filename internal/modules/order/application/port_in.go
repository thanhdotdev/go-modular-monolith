package orderapplication

import "context"

type UseCase interface {
	GetOrder(ctx context.Context, id string) (*OrderDTO, error)
}

package appapi

import (
	"context"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/crud"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type nestsService struct {
	client requester.Requester
}

func (s *nestsService) Eggs(nestID int) EggsService {
	return NewEggsService(s.client, nestID)
}

func NewNestsService(client requester.Requester) *nestsService {
	return &nestsService{client: client}
}

func (s *nestsService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.Nest, *api.Meta, error) {
	return crud.List[api.Nest](ctx, s.client, "/api/application/nests", options)
}

func (s *nestsService) ListAll(ctx context.Context) ([]*api.Nest, error) {
	return crud.ListAll[api.Nest](ctx, s.client, "/api/application/nests", 100)
}

// Get fetches a single nest by its ID.
func (s *nestsService) Get(ctx context.Context, nestID int) (*api.Nest, error) {
	return crud.Get[api.Nest](ctx, s.client, "/api/application/nests", nestID)
}

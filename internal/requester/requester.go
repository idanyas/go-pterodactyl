package requester

import (
	"context"
	"io"
	"net/http"

	"github.com/idanyas/go-pterodactyl/api"
)

type Requester interface {
	NewRequest(ctx context.Context, method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error)
	Do(ctx context.Context, req *http.Request, v any) (*http.Response, error)
}

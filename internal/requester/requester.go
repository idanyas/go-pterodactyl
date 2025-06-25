package requester

import (
	"github.com/davidarkless/go-pterodactyl/api"
	"io"
	"net/http"
)

type Requester interface {
	NewRequest(method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error)
	Do(req *http.Request, v any) (*http.Response, error)
}

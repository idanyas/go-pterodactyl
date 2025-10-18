package application

import (
	"context"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// PaginatorClient defines the interface required by the Paginator.
type PaginatorClient interface {
	Do(ctx context.Context, method, path string, body, v interface{}) (*http.Response, error)
}

// ApplicationClient is the client for the Application API.
type ApplicationClient interface {
	// User Management
	ListUsers(ctx context.Context, options pagination.ListOptions) ([]*models.User, *pagination.Paginator[*models.User], error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserExternal(ctx context.Context, externalID string) (*models.User, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error

	// Server Management
	ListServers(ctx context.Context, options pagination.ListOptions) ([]*models.Server, *pagination.Paginator[*models.Server], error)
	GetServer(ctx context.Context, id int) (*models.Server, error)
	GetServerExternal(ctx context.Context, externalID string) (*models.Server, error)
	CreateServer(ctx context.Context, req CreateServerRequest) (*models.Server, error)
	UpdateServerDetails(ctx context.Context, serverID int, req UpdateServerDetailsRequest) (*models.Server, error)
	UpdateServerBuild(ctx context.Context, serverID int, req UpdateServerBuildRequest) (*models.Server, error)
	UpdateServerStartup(ctx context.Context, serverID int, req UpdateServerStartupRequest) (*models.Server, error)
	SuspendServer(ctx context.Context, serverID int) error
	UnsuspendServer(ctx context.Context, serverID int) error
	ReinstallServer(ctx context.Context, serverID int) error
	DeleteServer(ctx context.Context, serverID int, force bool) error

	// Database Management
	ListServerDatabases(ctx context.Context, serverID int, options pagination.ListOptions) ([]*models.ApplicationDatabase, *pagination.Paginator[*models.ApplicationDatabase], error)
	GetServerDatabase(ctx context.Context, serverID, databaseID int) (*models.ApplicationDatabase, error)
	CreateServerDatabase(ctx context.Context, serverID int, req CreateServerDatabaseRequest) (*models.ApplicationDatabase, error)
	UpdateServerDatabase(ctx context.Context, serverID, databaseID int, req UpdateServerDatabaseRequest) (*models.ApplicationDatabase, error)
	ResetServerDatabasePassword(ctx context.Context, serverID, databaseID int) error
	DeleteServerDatabase(ctx context.Context, serverID, databaseID int) error

	// Node Management
	ListNodes(ctx context.Context, options pagination.ListOptions) ([]*models.Node, *pagination.Paginator[*models.Node], error)
	GetNode(ctx context.Context, id int) (*models.Node, error)
	CreateNode(ctx context.Context, req CreateNodeRequest) (*models.Node, error)
	UpdateNode(ctx context.Context, id int, req UpdateNodeRequest) (*models.Node, error)
	DeleteNode(ctx context.Context, id int) error
	GetNodeConfiguration(ctx context.Context, id int) (*models.NodeConfiguration, error)
	ListNodeAllocations(ctx context.Context, nodeID int, options pagination.ListOptions) ([]*models.Allocation, *pagination.Paginator[*models.Allocation], error)
	CreateNodeAllocations(ctx context.Context, nodeID int, req CreateNodeAllocationRequest) error
	DeleteNodeAllocation(ctx context.Context, nodeID, allocationID int) error
	GetDeployableNodes(ctx context.Context, memory, disk int64) ([]*models.Node, error)

	// Location Management
	ListLocations(ctx context.Context, options pagination.ListOptions) ([]*models.Location, *pagination.Paginator[*models.Location], error)
	GetLocation(ctx context.Context, id int) (*models.Location, error)
	CreateLocation(ctx context.Context, req CreateLocationRequest) (*models.Location, error)
	UpdateLocation(ctx context.Context, id int, req UpdateLocationRequest) (*models.Location, error)
	DeleteLocation(ctx context.Context, id int) error

	// Nest & Egg Management
	ListNests(ctx context.Context, options pagination.ListOptions) ([]*models.Nest, *pagination.Paginator[*models.Nest], error)
	GetNest(ctx context.Context, id int) (*models.Nest, error)
	ListNestEggs(ctx context.Context, nestID int, options pagination.ListOptions) ([]*models.Egg, *pagination.Paginator[*models.Egg], error)
	GetEgg(ctx context.Context, nestID, eggID int) (*models.Egg, error)
}

type client struct {
	client PaginatorClient
}

// New creates a new Application API client.
func New(c PaginatorClient) ApplicationClient {
	return &client{client: c}
}

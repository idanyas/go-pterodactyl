package application

import "github.com/davidarkless/go-pterodactyl/api"

type AllocationsService interface {
	List(options *api.PaginationOptions) ([]*api.Allocation, *api.Meta, error)
	ListAll() ([]*api.Allocation, error)
	Create(options api.AllocationCreateOptions) error //TODO Include int / allocation on return?
	Delete(allocationID int) error
}

type DatabaseService interface {
	List(options api.PaginationOptions) ([]*api.Database, *api.Meta, error)
	Get(databaseID int) (*api.Database, error)
	Create(options api.DatabaseCreateOptions) (*api.Database, error)
	ResetPassword(databaseID int) error
	Delete(databaseID int) error
}

type NodesService interface {
	List(options *api.PaginationOptions) ([]*api.Node, *api.Meta, error)
	ListAll() ([]*api.Node, error)
	Get(id int) (*api.Node, error)
	GetConfiguration(nodeID int) (*api.NodeConfiguration, error)
	Create(options api.NodeCreateOptions) (*api.Node, error)
	Update(nodeID int, options api.NodeUpdateOptions) (*api.Node, error)
	Delete(nodeID int) error
	Allocations(nodeID int) AllocationsService
}

type EggsService interface {
	List(options *api.PaginationOptions) ([]*api.Egg, *api.Meta, error)
	ListAll() ([]*api.Egg, error)
	Get(eggID int) (*api.Egg, error)
}

// NestsService defines the actions for nests and provides access to egg management.
type NestsService interface {
	List(options *api.PaginationOptions) ([]*api.Nest, *api.Meta, error)
	ListAll() ([]*api.Nest, error)
	Get(id int) (*api.Nest, error)
	Eggs(nestID int) EggsService // Returns the EggsService interface
}

// UsersService defines the actions for users.
type UsersService interface {
	List(options *api.PaginationOptions) ([]*api.User, *api.Meta, error)
	ListAll() ([]*api.User, error)
	Get(id int) (*api.User, error)
	GetExternalID(externalId string) (*api.User, error)
	Create(options api.UserCreateOptions) (*api.User, error)
	Update(id int, options api.UserUpdateOptions) (*api.User, error)
	Delete(id int) error
}

type ServersService interface {
	List(options api.PaginationOptions) ([]*api.Server, *api.Meta, error)
	ListAll() ([]*api.Server, error)
	Get(id int) (*api.Server, error)
	GetExternal(externalID string) (*api.Server, error)
	Create(options api.ServerCreateOptions) (*api.Server, error)
	UpdateDetails(serverID int, options api.ServerUpdateDetailsOptions) (*api.Server, error)
	UpdateBuild(serverID int, options api.ServerUpdateBuildOptions) (*api.Server, error)
	UpdateStartup(serverID int, options api.ServerUpdateStartupOptions) (*api.Server, error)
	Suspend(serverID int) error
	Unsuspend(serverID int) error
	Reinstall(serverID int) error
	Delete(serverID int, force bool) error
	Databases(serverID int) DatabaseService
}

// LocationsService defines the actions for locations.
type LocationsService interface {
	List(options *api.PaginationOptions) ([]*api.Location, *api.Meta, error)
	ListAll() ([]*api.Location, error)
	Get(id int) (*api.Location, error)
	Create(options api.LocationCreateOptions) (*api.Location, error)
	Update(id int, options api.LocationUpdateOptions) (*api.Location, error)
	Delete(id int) error
}

// Application is the container for all top-level API services.
type Application struct {
	Users     UsersService
	Nodes     NodesService
	Locations LocationsService
	Servers   ServersService
	Nests     NestsService
}

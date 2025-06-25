package clientapi

import (
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
	"io"
)

type APIKeysService interface {
	List(options api.PaginationOptions) ([]*api.APIKey, *api.Meta, error)
	Create(options api.APIKeyCreateOptions) (*api.APIKey, error)
	Delete(identifier string) error
}

type DatabasesService interface {
	List(options api.PaginationOptions) ([]*api.ClientDatabase, *api.Meta, error)
	Create(options api.ClientDatabaseCreateOptions) (*api.ClientDatabase, error)
	RotatePassword(databaseID string) (*api.ClientDatabase, error)
	Delete(databaseID string) error
}

type FileService interface {
	List(directory string) ([]*api.FileObject, error)
	GetContents(filePath string) (string, error)
	Download(filePath string) (*api.SignedURL, error)
	Rename(options api.RenameFilesOptions) error
	Copy(options api.CopyFileOptions) error
	Write(filePath string, content io.Reader) error
	Compress(options api.CompressFilesOptions) (*api.FileObject, error)
	Decompress(options api.DecompressFileOptions) error
	Delete(options api.DeleteFilesOptions) error
	CreateFolder(options api.CreateFolderOptions) error
	GetUploadURL() (*api.SignedURL, error)
}

type ScheduleService interface {
	List(options api.PaginationOptions) ([]*api.Schedule, *api.Meta, error)
	Create(options api.ScheduleCreateOptions) (*api.Schedule, error)
	Details(scheduleID int) (*api.Schedule, error)
	Update(scheduleID int, options api.ScheduleUpdateOptions) (*api.Schedule, error)
	Delete(scheduleID int) error
	CreateTask(scheduleID int, options api.TaskCreateOptions) (*api.Task, error)
	UpdateTask(scheduleID, taskID int, options api.TaskUpdateOptions) (*api.Task, error)
	DeleteTask(scheduleID, taskID int) error
}

type NetworkService interface {
	ListAllocations(options api.PaginationOptions) ([]*api.Allocation, *api.Meta, error)
	AssignAllocation() (*api.Allocation, error)
	SetAllocationNote(allocationID int, options api.AllocationNoteOptions) (*api.Allocation, error)
	SetPrimaryAllocation(allocationID int) (*api.Allocation, error)
	UnassignAllocation(allocationID int) error
}

type UsersService interface {
	List(options api.PaginationOptions) ([]*api.Subuser, *api.Meta, error)
	Create(options api.SubuserCreateOptions) (*api.Subuser, error)
	Details(uuid string) (*api.Subuser, error)
	Update(uuid string, options api.SubuserUpdateOptions) (*api.Subuser, error)
	Delete(uuid string) error
}

type BackupService interface {
	List(options api.PaginationOptions) ([]*api.Backup, *api.Meta, error)
	Create(options api.BackupCreateOptions) (*api.Backup, error)
	Details(uuid string) (*api.Backup, error)
	Download(uuid string) (*api.BackupDownload, error)
	Delete(uuid string) error
}

type StartupService interface {
	ListVariables(options api.PaginationOptions) ([]*api.StartupVariable, *api.Meta, error)
	UpdateVariable(options api.UpdateVariableOptions) (*api.StartupVariable, error)
}

type SettingsService interface {
	Rename(options api.RenameOptions) error
	Reinstall() error
}

type ServersService interface {
	GetDetails() (*api.ClientServer, error)
	GetWebsocket() (*api.WebsocketDetails, error)
	GetResources() (*api.Resources, error)
	SendCommand(command string) error
	SetPowerState(signal string) error

	Databases() DatabasesService
	Files() FileService
	Schedules() ScheduleService
	Network() NetworkService
	Users() UsersService
	Backups() BackupService
	Startup() StartupService
	Settings() SettingsService
}

type AccountService interface {
	GetDetails() (*api.Account, error)
	GetTwoFactorDetails() (*api.TwoFactorDetails, error)
	EnableTwoFactor(options api.TwoFactorEnableOptions) error
	DisableTwoFactor(options api.TwoFactorDisableOptions) error
	UpdateEmail(options api.UpdateEmailOptions) error
	UpdatePassword(options api.UpdatePasswordOptions) error
	APIKeys() APIKeysService
}

type ClientAPI interface {
	ListServers(options api.PaginationOptions) ([]*api.ClientServer, *api.Meta, error)
	ListPermissions() (*api.Permission, error)

	Servers(identifier string) ServersService
	Account() AccountService
}

type ClientAPIService struct {
	client requester.Requester
}

func NewClientAPI(client requester.Requester) *ClientAPIService {
	return &ClientAPIService{client: client}
}

func (s *ClientAPIService) ListServers(options api.PaginationOptions) ([]*api.ClientServer, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/client", nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client server list request: %w", err)
	}

	response := &api.PaginatedResponse[api.ClientServer]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.ClientServer, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *ClientAPIService) ListPermissions() (*api.Permission, error) {
	req, err := s.client.NewRequest("GET", "/api/client/permissions", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create permissions request: %w", err)
	}

	response := &api.Permission{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *ClientAPIService) Servers(identifier string) ServersService {
	return newServerService(c.client, identifier)
}

func (c *ClientAPIService) Account() AccountService {
	return newAccountService(c.client)
}

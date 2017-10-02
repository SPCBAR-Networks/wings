package control

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pterodactyl/wings/config"
	"github.com/Pterodactyl/wings/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"unicode"
)

// ErrServerExists is returned when a server already exists on creation.
type ErrServerExists struct {
	id string
}

func (e ErrServerExists) Error() string {
	return "server " + e.id + " already exists"
}

// Server is a Server
type Server interface {
	Start() error
	Stop() error
	Restart() error
	Kill() error
	Exec(command string) error
	Rebuild() error

	Save() error

	Environment() (Environment, error)

	HasPermission(string, string) bool
}

// ServerStruct is a single instance of a Service managed by the panel
type ServerStruct struct {
	// ID is the unique identifier of the server
	ID string `yaml:"uuid"`

	// ServiceName is the name of the service. It is mainly used to allow storing the service
	// in the config
	ServiceName string `yaml:"serviceName"`
	service     *Service
	environment Environment

	// StartupCommand is the command executed in the environment to start the server
	StartupCommand string `yaml:"startupCommand"`

	// DockerContainer holds information regarding the docker container when the server
	// is running in a docker environment
	DockerContainer Container `yaml:"dockerContainer"`

	// EnvironmentVariables are set in the Environment the server is running in
	EnvironmentVariables map[string]string `yaml:"env"`

	// Allocations contains the ports and ip addresses assigned to the server
	Allocations Allocations `yaml:"allocation"`

	// Settings are the environment settings and limitations for the server
	Settings Settings `yaml:"settings"`

	// Keys are some auth keys we will hopefully replace by something better.
	Keys map[string][]string `yaml:"keys"`
}

type Allocations struct {
	Ports       []int16 `yaml:"ports"`
	PrimaryIp   string  `yaml:"ip"`
	PrimaryPort int16   `yaml:"port"`
}

type Settings struct {
	Memory int64  `yaml:"memory"`
	Swap   int64  `yaml:"swap"`
	IO     int64  `yaml:"io"`
	CPU    int16  `yaml:"cpu"`
	Disk   int64  `yaml:"disk"`
	Image  string `yaml:"image"`
	User   string `yaml:"user"`
	UserID int16  `yaml:"userID"`
}

type Container struct {
	ID    string `yaml:"id"`
	Image string `yaml:"image"`
}

// ensure server implements Server
var _ Server = &ServerStruct{}

type serversMap map[string]*ServerStruct

var servers = make(serversMap)

// LoadServerConfigurations loads the configured servers from a specified path
func LoadServerConfigurations(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	// Can this be removed? Isn't it set in the global scope for this package above?
	// servers = make(serversMap)

	for _, f := range files {
		if f.IsDir() {
			server, err := LoadServerFromDisk(filepath.Join(path, f.Name(), constants.ServerConfigFile))
			if err != nil {
				return err
			}

			servers[server.ID] = server
		}
	}

	return nil
}

func LoadServerFromDisk(path string) (*ServerStruct, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := &ServerStruct{}
	if err := yaml.Unmarshal(f, s); err != nil {
		return nil, err
	}

	return s, nil
}

// Writes a server configuration to the disk in YAML format. This will use the defined
// location from the configuration to save to.
func WriteServerToDisk(server *ServerStruct) error {
	y, err := yaml.Marshal(server)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(server.Path(), constants.DefaultFolderPerms); err != nil {
		return err
	}

	if err := ioutil.WriteFile(server.ConfigPath(), y, constants.DefaultFilePerms); err != nil {
		return err
	}

	return nil
}

// Loops through all configured servers and saves their configuration to the disk
// using the WriteServerToDisk function.
func WriteConfigurationsToDisk() error {
	for _, s := range servers {
		if err := WriteServerToDisk(s); err != nil {
			return err
		}
	}
	return nil
}

// GetServers returns an array of all servers the daemon manages
func GetServers() []Server {
	serverArray := make([]Server, len(servers))
	i := 0
	for _, s := range servers {
		serverArray[i] = s
		i++
	}
	return serverArray
}

// GetServer returns the server identified by the provided uuid
func GetServer(id string) Server {
	server := servers[id]
	if server == nil {
		return nil // https://golang.org/doc/faq#nil_error
	}
	return server
}

// CreateServer creates a new server
func CreateServer(server *ServerStruct) (Server, error) {
	if servers[server.ID] != nil {
		return nil, ErrServerExists{server.ID}
	}
	servers[server.ID] = server
	if err := server.Save(); err != nil {
		return nil, err
	}
	return server, nil
}

// Delete a server and all associated files and pointers from Wings.
func DeleteServer(id string) error {
	if err := DeleteServerData(id); err != nil {
		log.WithField("server", id).WithError(err).Error("Failed to delete a server data directory.")
	}

	if err := DeleteServerConfig(id); err != nil {
		log.WithField("server", id).WithError(err).Error("Failed to delete a server configuration directory.")
	}

	delete(servers, id)
	return nil
}

// Deletes a server configuration file from Wings.
// TODO: should these actually error if there is no directory? Can we not just move on?
func DeleteServerConfig(id string) error {
	p := filepath.Join(viper.GetString(config.DataPath), id)
	f, err := os.Stat(p)
	if os.IsNotExist(err) || !f.IsDir() {
		return err
	}

	return os.RemoveAll(p)
}

// Delete the data folder for a specific server. This is a permanent
// deletion and will remove all of the data stored (by default) in
// /srv/wings/<server uuid> as well as the folder.
func DeleteServerData(id string) error {
	p := filepath.Join(viper.GetString(config.ServerDataPath), id)
	f, err := os.Stat(p)
	if os.IsNotExist(err) || !f.IsDir() {
		return err
	}

	return os.RemoveAll(p)
}

// Start a new environment process for a server. If an environment
// does not currently exist for the server a new one will be created.
func (s *ServerStruct) Start() error {
	e, err := s.Environment()
	if err != nil {
		return err
	}

	if !e.Exists() {
		if err := e.Create(); err != nil {
			return err
		}
	}

	return e.Start()
}

// Stop a running server environment.
func (s *ServerStruct) Stop() error {
	e, err := s.Environment()
	if err != nil {
		return err
	}

	return e.Stop()
}

// Restart a running server environment.
func (s *ServerStruct) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}

	return s.Start()
}

// Kill a running server environment.
func (s *ServerStruct) Kill() error {
	e, err := s.Environment()
	if err != nil {
		return err
	}

	return e.Kill()
}

// Execute a command in a running server environment.
func (s *ServerStruct) Exec(c string) error {
	e, err := s.Environment()
	if err != nil {
		return err
	}

	return e.Exec(c)
}

// Rebuild a server.
func (s *ServerStruct) Rebuild() error {
	e, err := s.Environment()
	if err != nil {
		return err
	}

	return e.ReCreate()
}

// Return the service configuration for a server. This will include the
// environment name, and if relevant, the docker image being used.
func (s *ServerStruct) Service() *Service {
	if s.service == nil {
		// TODO: use the correct service, mock for now.
		s.service = &Service{
			DockerImage:     "quay.io/pterodactyl/core:java",
			EnvironmentName: "docker",
		}
	}

	return s.service
}

// Return the first chunk of the server UUID.
func (s *ServerStruct) UUIDShort() string {
	return s.ID[0:strings.Index(s.ID, "-")]
}

// Return information about the server environment.
func (s *ServerStruct) Environment() (Environment, error) {
	var err error
	if s.environment == nil {
		switch s.Service().EnvironmentName {
		case "docker":
			s.environment, err = NewDockerEnvironment(s)
		default:
			log.WithField("service", s.ServiceName).Error("An invalid environment name was provided for this server. The currently enabled environments are: docker")
			return nil, errors.New("invalid environment argument was provided, currently accepted are: docker")
		}
	}

	return s.environment, err
}

// HasPermission checks wether a provided token has a specific permission
// TODO: fix this function
func (s *ServerStruct) HasPermission(token string, permission string) bool {
	return true
	//for key, perms := range s.Keys {
	//	if key == token {
	//		for _, perm := range perms {
	//			if perm == permission || perm == "s:*" {
	//				return true
	//			}
	//		}
	//		return false
	//	}
	//}
	//return false
}

// Save a server configuration by writing it to the disk.
func (s *ServerStruct) Save() error {
	if err := WriteServerToDisk(s); err != nil {
		log.WithField("server", s.ID).WithError(err).Error("Unable to write a server configuration to the disk.")
		return err
	}

	return nil
}

// Generate a path to the server folder. In most instances this will not return the correct path
// for writing or modifying files. One should use DataPath(path string) to perform those actions.
func (s *ServerStruct) Path() string {
	return filepath.Join(viper.GetString(config.ServerDataPath), s.ID)
}

// Generate a path to the server data folder. By default with no argument this is the path
// that is mounted into the docker container. If a path is passed it will be normalized
// to the root server data folder. This prevents a user from providing a path like ../../etc and
// escaping thier server data directory.
func (s *ServerStruct) DataPath(p string) string {
	b := filepath.Join(s.Path(), constants.ServerDataFolder)
	strip := func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, s)
	}

	f := filepath.Join(b, strip(p))
	if !strings.HasPrefix(f, s.Path()) {
		return b
	}

	return f
}

// Returns the path to the server configuration YAML file. This path is different than Path() as
// it refers to a completely different part of the filesystem and is not mounted into a docker container.
func (s *ServerStruct) ConfigPath() string {
	return filepath.Join(viper.GetString(config.DataPath), s.ID, constants.ServerConfigFile)
}

package control

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Server is a Server
type Server interface {
	Start() error
	Stop() error
	Exec(command string) error
	Rebuild() error

	HasPermission(string, string) bool
}

// server is a single instance of a Service managed by the panel
type ServerStruct struct {
	// ID is the unique identifier of the server
	ID string `json:"uuid"`

	// ServiceName is the name of the service. It is mainly used to allow storing the service
	// in the config
	ServiceName string `json:"serviceName"`
	service     *service
	environment Environment

	// StartupCommand is the command executed in the environment to start the server
	StartupCommand string `json:"startupCommand"`

	// DockerContainer holds information regarding the docker container when the server
	// is running in a docker environment
	DockerContainer dockerContainer `json:"dockerContainer"`

	// EnvironmentVariables are set in the Environment the server is running in
	EnvironmentVariables map[string]string `json:"env"`

	// Allocations contains the ports and ip addresses assigned to the server
	Allocations allocations `json:"allocation"`

	// Settings are the environment settings and limitations for the server
	Settings settings `json:"settings"`

	// Keys are some auth keys we will hopefully replace by something better.
	Keys map[string][]string `json:"keys"`
}

type allocations struct {
	Ports    []int16 `json:"ports"`
	MainIP   string  `json:"ip"`
	MainPort int16   `json:"port"`
}

type settings struct {
	Memory int64  `json:"memory"`
	Swap   int64  `json:"swap"`
	IO     int64  `json:"io"`
	CPU    int16  `json:"cpu"`
	Disk   int64  `json:"disk"`
	Image  string `json:"image"`
	User   string `json:"user"`
	UserID int16  `json:"userID"`
}

type dockerContainer struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

// ensure server implements Server
var _ Server = &ServerStruct{}

type serversMap map[string]*ServerStruct

var servers = make(serversMap)

// LoadServerConfigurations loads the configured servers from a specified path
func LoadServerConfigurations(path string) error {
	serverFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	servers = make(serversMap)

	for _, file := range serverFiles {
		if !file.IsDir() {
			server, err := loadServerConfiguration(path + file.Name())
			if err != nil {
				return err
			}
			servers[server.ID] = server
		}
	}

	return nil
}

func loadServerConfiguration(path string) (*ServerStruct, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	server := &ServerStruct{}
	if err := json.Unmarshal(file, server); err != nil {
		return nil, err
	}
	return server, nil
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
	servers[server.ID] = server
	return server, nil
}

// DeleteServer deletes a server and all related files
// NOTE: This is not reversible.
func DeleteServer(uuid string) error {
	delete(servers, uuid)
	return nil
}

func (s *ServerStruct) Start() error {
	/*if err := s.Environment().Create(); err != nil {
		return err
	}
	if err := s.Environment().Start(); err != nil {
		return err
	}*/
	return nil
}

func (s *ServerStruct) Stop() error {
	/*if err := s.Environment().Stop(); err != nil {
		return err
	}*/
	return nil
}

func (s *ServerStruct) Exec(command string) error {
	/*if err := s.Environment().Exec(command); err != nil {
		return err
	}*/
	return nil
}

func (s *ServerStruct) Rebuild() error {
	/*if err := s.Environment().ReCreate(); err != nil {
		return err
	}*/
	return nil
}

// Service returns the server's service configuration
func (s *ServerStruct) Service() *service {
	if s.service == nil {
		// TODO: Properly use the correct service, mock for now.
		s.service = &service{
			DockerImage:     "quay.io/pterodactyl/core:java",
			EnvironmentName: "docker",
		}
	}
	return s.service
}

// UUIDShort returns the first block of the UUID
func (s *ServerStruct) UUIDShort() string {
	return s.ID[0:strings.Index(s.ID, "-")]
}

// Environment returns the servers environment
func (s *ServerStruct) Environment() (Environment, error) {
	var err error
	if s.environment == nil {
		switch s.Service().EnvironmentName {
		case "docker":
			s.environment, err = NewDockerEnvironment(s)
		default:
			log.WithField("service", s.ServiceName).Error("Invalid environment name")
			return nil, errors.New("Invalid environment name")
		}
	}
	return s.environment, err
}

// HasPermission checks wether a provided token has a specific permission
// @todo fix this function
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

func (s *ServerStruct) save() {

}

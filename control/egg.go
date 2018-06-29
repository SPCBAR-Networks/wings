package control

type Egg struct {
	server *Server

	// EnvironmentName is the name of the environment used by the egg
	EnvironmentName string `json:"environmentName" jsonapi:"primary,service"`

	DockerImage string `json:"dockerImage" jsonapi:"attr,docker_image"`
}

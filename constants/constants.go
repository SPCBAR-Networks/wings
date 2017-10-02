package constants

import "os"

// Version is the current wings version.
const Version = "0.0.0-canary"

// DefaultFilePerms are the file perms used for created files.
const DefaultFilePerms os.FileMode = 0644

// DefaultFolderPerms are the file perms used for created folders.
const DefaultFolderPerms os.FileMode = 0744

// ServersPath is the path of the servers within the configured DataPath.
const ServersPath string = "servers"

// ServerConfigFile is the filename of the server config file.
const ServerConfigFile string = "server.yaml"

// ServerDataFolder is the name of the folder containing all of the server
// files for a given server.
const ServerDataFolder string = "data"

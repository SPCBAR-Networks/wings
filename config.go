package main

import (
	"encoding/json"
	"os"
)

//daemonConfig struct
type daemonConfig struct {
	//web settings
	web     string `json:"web"`
	webip   string `json:"host"`
	webport int    `json:"listen"`
	//ssl settings
	ssl     string `json:"ssl"`
	sslenbl string `json:"enabled"`
	sslcert string `json:"certificate"`
	sslkey  string `json:"key"`
	//Docker config struct
	dckr     string `json:"docker"`
	dckrsckt string `json:"socket"`
	dckritfc string `json:"interface"`
	dckrupdt string `json:"autoupdate_images"`
	dckrtzon string `json:"timezone_path"`
	//sftp server struct
	sftp     string `json:"sftp"`
	sftppath string `json:"path"`
	sftpport int    `json:"port"`
	//container query struct
	query    string `json:"query"`
	killfail string `json:"kill_on_fail"`
	faillimt string `json:"fail_limit"`
	//logger struct
	logger  string `json:"logger"`
	logpath string `json:"path"`
	logsrc  string `json:"src"`
	loglevl string `json:"level"`
	logperd string `json:"period"`
	logcnt  int    `json:"count"`
	//remote panel struct
	remote    string `json:"remote"`
	remtbase  string `json:"base"`
	remtdnld  string `json:"download"`
	remtinstl string `json:"installed"`
	//uploads struct
	upld      string `json:"uploads"`
	upldmax   string `json:"maximumSize"`
	upldlimit string `json:"size_limit"`
	//keys struct
	keys []string `json:"keys"`
}

//serverConfig struct
type serverConfig struct {
}

func getDaemonConfig(in string) string {
	//Opens config.json and returns values
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	daemon := daemonConfig{}

	err := decoder.Decode(&daemon)
	if err != nil {
		log.Fatal("issue decoding the config file", err)
	}

	out := daemon.dckrsckt
	return out
}

func getServerConfig() {

}

package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type Project struct {
	Folder     string `json:"folder"`   //项目文件夹
	Hostname   string `json:"hostname"` //机器域名
	SshPort    int    `json:"ssh_port"` //ssh 端口号
	RunCommand string `json:"run_command"`
}
type configs struct {
	Projects map[string]Project
	Servers  map[string]Server
}

var configfile = "config.json"
var serverfile = "server.json"

func init() {
	LoadConfigFile()
}
func NewConfigs() *configs {
	return &configs{
		Projects: make(map[string]Project),
		Servers:  make(map[string]Server),
	}
}

func GetRootPath() string {
	file, _ := exec.LookPath(os.Args[0])
	paths, _ := filepath.Abs(path.Dir(file))
	return paths
}

var conf = NewConfigs()

func LoadConfigFile() {
	content, err := ioutil.ReadFile(GetRootPath() + "/" + configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf.Projects)
	if err != nil {
		log.Fatal(err)
	}
	content, err = ioutil.ReadFile(GetRootPath() + "/" + serverfile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf.Servers)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *configs {
	return conf
}

func GetProjectConfig(projectname string) (Project, error) {
	for k, v := range conf.Projects {
		if k != projectname {
			continue
		}
		return v, nil
	}
	return Project{}, errors.New("project not found")
}

func GetServerConfig(serverName string) (Server, error) {
	for k, v := range conf.Servers {
		if k != serverName {
			continue
		}
		return v, nil
	}
	return Server{}, errors.New("server not found")
}

func ProjectExist(key string, projects map[string]Project) bool {
	for name := range projects {
		if name == key {
			return true
		}
	}
	return false
}

type Server struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

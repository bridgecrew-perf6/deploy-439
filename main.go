package main

import (
	"deploy/parser"
	"deploy/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var project = flag.String("p", "", "project name")

func init() {
	flag.Parse()
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
}
func main() {
	config, err := parser.GetProjectConfig(*project)
	checkError(err)
	//gzip压缩文件夹
	filename, err := utils.GzipFile(config.Folder)
	checkError(err)
	//ssh 上传压缩包文件
	remoteFileName, err := utils.SshUpload(filename, config.Hostname, config.SshPort)
	checkError(err)
	fmt.Println(remoteFileName)
	//解压缩包
	//utils.UnGzipFile(remoteFileName)
	remoteUnGzip(config, remoteFileName)
	//运行镜像
	StopImage(config)
	RunImage(config)
}

func StopImage(config parser.Project) {
	session, err := utils.NewSSHSession(config.Hostname, config.SshPort)
	checkError(err)
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	session.Run(fmt.Sprintf("docker ps | grep %s:test | awk '{print $1}' | xargs docker stop", utils.BaseName(config.Folder)))
}

func RunImage(config parser.Project) {
	session, err := utils.NewSSHSession(config.Hostname, config.SshPort)
	checkError(err)
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	log.Println(config.RunCommand)
	err = session.Run(config.RunCommand)
	checkError(err)
}

//远程解压缩
func remoteUnGzip(config parser.Project, remoteFileName string) {
	session, err := utils.NewSSHSession(config.Hostname, config.SshPort)
	checkError(err)
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	remoteFile := strings.TrimSuffix(utils.BaseName(remoteFileName), ".tar.gz")
	projectName := utils.BaseName(config.Folder)
	log.Println(fmt.Sprintf("bash /vagrant/launch.sh -p %s -d %s", remoteFile, projectName))
	err = session.Run(fmt.Sprintf("bash /vagrant/launch.sh -p %s -d %s", remoteFile, projectName))
	checkError(err)

}

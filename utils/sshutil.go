package utils

import (
	"deploy/parser"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func GetSshConfig(hostname string) *ssh.ClientConfig {
	server, err := parser.GetServerConfig(hostname)
	if err != nil {
		log.Println(hostname, "server account not found")
		return nil
	}
	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}
	return sshConfig
}
func sshAddress(hostname string, port int) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func SshUpload(filename string, hostname string, port int) (string, error) {
	sshClient, err := ssh.Dial("tcp", sshAddress(hostname, port), GetSshConfig(hostname))
	if err != nil {
		return "", err
	}
	defer sshClient.Close()
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()

	remoteFileName := fmt.Sprintf("/vagrant/%s", BaseName(filename))
	remoteFile, err := sftpClient.Create(remoteFileName)
	if err != nil {
		return "", err
	}
	defer remoteFile.Close()

	localFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer localFile.Close()
	n, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return "", err
	}
	localFileInfo, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	log.Printf("文件上传成功[%s->%s]本地文件大小：%s，上传文件大小： %s", filename, remoteFileName, formatFileSize(localFileInfo.Size()), formatFileSize(n))
	return remoteFileName, nil
}

func formatFileSize(s int64) (size string) {
	if s < 1024 {
		return fmt.Sprintf("%.2fB", float64(s)/float64(1))
	} else if s < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(1024))
	} else if s < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(1024*1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(1024*1024*1024*1024))
	} else { //if s < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(s)/float64(1024*1024*1024*1024*1024))
	}
}

func NewClient(hostname string, port int) (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", sshAddress(hostname, port), GetSshConfig(hostname))
	return sshClient, err
}

func NewSSHSession(hostname string, port int) (*ssh.Session, error) {
	sshClient, err := NewClient(hostname, port)
	if err != nil {
		return nil, err
	}
	if sshSession, err := sshClient.NewSession(); err != nil {
		return nil, err
	} else {
		return sshSession, nil
	}
}

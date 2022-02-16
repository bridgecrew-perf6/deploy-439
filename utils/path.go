package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetTempFolder() string {
	return os.TempDir()
}

func Getmd5(content string) string {
	h := md5.New()
	io.WriteString(h, content)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

func BaseName(filename string) string {
	projectName := strings.TrimSuffix(filename, "/")
	projectName = strings.TrimSuffix(filename, "\\")
	for i := len(projectName) - 1; i >= 0; i-- {
		if projectName[i] == '/' || projectName[i] == '\\' {
			projectName = projectName[i+1:]
			break
		}
	}
	return projectName
}

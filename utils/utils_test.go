package utils

import (
	"os"
	"testing"
)

func TestPath(t *testing.T) {
	t.Log(GetTempFolder())
}

func TestGzip(t *testing.T) {
	t.Log(GzipFile("D:/programming/web"))
}

func TestUnGzip(t *testing.T) {
	t.Log(UnGzipFile("fun.tar.gz"))
}

func TestSshUpload(t *testing.T) {
	t.Log(SshUpload("fun.tar.gz", "192.168.33.12", 22))
}

func TestFindLast(t *testing.T) {
	t.Log(findLast("d:/programming/web"))
	t.Log(findLast("d:/programming/web12/"))
}

func TestNewSession(t *testing.T) {
	session, err := NewSSHSession("192.168.33.12", 22)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	err = session.Run("ls")
	if err != nil {
		t.Fatal(err)
	}
}

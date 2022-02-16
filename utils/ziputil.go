package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GzipFile(dirPath string) (string, error) {
	var filename = GetTempFolder() + GetNowTimeStamp() + Getmd5(dirPath) + ".tar.gz"
	//var filename = "fun.tar.gz"
	var buf bytes.Buffer
	if err := compress(dirPath, &buf); err != nil {
		return filename, err
	}
	var fileToWrite, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if err != nil {
		return filename, err
	}
	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		return filename, err
	}
	return filename, nil
}

var dstfolder string

func init() {
	if runtime.GOOS == "windows" {
		dstfolder = "d:/temp"
	} else {
		dstfolder = "/vagrant"
	}
}

var errWrongFormat = errors.New("wrong file format, cannot be gzip decoded")

func UnGzipFile(filename string) error {
	if !strings.HasSuffix(filename, ".tar.gz") {
		return errWrongFormat
	}
	var f, err = os.Open(filename)
	if err != nil {
		return err
	}
	err = uncompress(f, dstfolder)
	if err != nil {
		return err
	}
	return nil
}

func stripPrefix(filename string, prefix string) string {
	if runtime.GOOS == "windows" {
		filename = strings.ReplaceAll(filename, "/", "\\")
		prefix = strings.ReplaceAll(prefix, "/", "\\")
	} else {
		filename = strings.ReplaceAll(filename, "\\", "/")
		prefix = strings.ReplaceAll(prefix, "\\", "/")
	}
	filename = strings.TrimPrefix(filename, findLast(prefix))
	return filename
}

func findLast(src string) string {
	src = strings.TrimSuffix(src, string('\\'))
	src = strings.TrimSuffix(src, string('/'))
	var i int
	for i = len(src) - 1; i != -1; i-- {
		if src[i] == '/' || src[i] == '\\' {
			break
		}
	}
	var name = src[:i+1]
	return name
}
func compress(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if file == src {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(stripPrefix(file, src))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err := tw.Close(); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

func uncompress(src io.Reader, dst string) error {
	//ungzip
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)

	//uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := header.Name

		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", target)
		}

		target = filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			fileToWrite.Close()
		}
	}
	return nil
}

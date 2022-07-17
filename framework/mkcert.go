package framework

import (
	_ "embed"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed mkcert/mkcert-windows-amd64.exe
var mkcert []byte

func GenerateSSL(mkcertPath string) (err error) {
	err = exec.Command(mkcertPath, "-install").Run()
	if err != nil {
		return err
	}
	return exec.Command(mkcertPath, "localhost", "127.0.0.1", "::1").Run()
}

func CreateMkcertIfNotExists() (mkcertPath string, err error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	mkcertPath = filepath.Join(hd, ".mkcert", "mkcert-windows-amd64.exe")
	if _, err = os.Stat(mkcertPath); errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir(filepath.Join(hd, ".mkcert"), 0755)
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile(mkcertPath, mkcert, 0644)
	}

	return mkcertPath, err
}

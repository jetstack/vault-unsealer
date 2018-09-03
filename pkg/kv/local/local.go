package local

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type Local struct {
	keyDir string
}

func New(keyDir string) (*Local, error) {
	dir, err := homedir.Expand(keyDir)
	if err != nil {
		return nil, err
	}

	return &Local{
		keyDir: dir,
	}, nil
}

func (l *Local) Set(key string, value []byte) error {
	path := filepath.Join(l.keyDir, key)
	return ioutil.WriteFile(path, value, os.FileMode(0600))
}

func (l *Local) Get(key string) ([]byte, error) {
	path := filepath.Join(l.keyDir, key)
	return ioutil.ReadFile(path)
}

func (l *Local) Test(key string) error {
	path := filepath.Join(l.keyDir, key)
	_, err := os.Stat(path)
	return err
}

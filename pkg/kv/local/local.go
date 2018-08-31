package local

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Local struct {
	keyPath string
}

func New(keyPath string) (*Local, error) {
	path, err := homedir.Expand(keyPath)
	if err != nil {
		return nil, err
	}

	return &Local{
		keyPath: path,
	}, nil
}

func (l *Local) Set(key string, value []byte) error {
	return ioutil.WriteFile(l.keyPath, value, os.FileMode(0600))
}

func (l *Local) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(l.keyPath)
}

func (l *Local) Test(key string) error {
	_, err := os.Stat(l.keyPath)
	return err
}

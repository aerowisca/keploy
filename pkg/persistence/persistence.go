package persistence

import (
	"fmt"
	"go.keploy.io/server/pkg"
	"os"
	"path/filepath"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Filesystem
type Filesystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	CreateYamlFile(folderLocation string, fileNameWithoutExtension string) (bool, error)
}

// Native filesystem is a wrapper over the os functions for local storage.
type Native struct {
}

func NewNativeFilesystem() Filesystem {
	return &Native{}
}

func (n *Native) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (n *Native) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (n *Native) CreateYamlFile(folderLocation string, fileNameWithoutExtension string) (bool, error) {
	if !pkg.IsValidPath(folderLocation) {
		return false, fmt.Errorf("file path should be absolute. got path: %s",
			pkg.SanitiseInput(folderLocation))
	}
	if _, err := os.Stat(filepath.Join(folderLocation, fileNameWithoutExtension+".yaml")); err != nil {
		err := os.MkdirAll(filepath.Join(folderLocation), os.ModePerm)
		if err != nil {
			return false, fmt.Errorf("failed to create a mock dir. error: %v", err.Error())
		}
		_, err = os.Create(filepath.Join(folderLocation, fileNameWithoutExtension+".yaml"))
		if err != nil {
			return false, fmt.Errorf("failed to create a yaml file. error: %v", err.Error())
		}
		return true, nil
	} else {
		return false, err
	}
}

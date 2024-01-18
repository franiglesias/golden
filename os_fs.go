package golden

import (
	"errors"
	"os"
	"path"
)

type OsFs struct {
}

func NewOsFs() OsFs {
	return OsFs{}
}

func (o OsFs) Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func (o OsFs) WriteFile(name string, data []byte) error {
	p := path.Dir(name)
	_, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return os.WriteFile(name, data, os.ModePerm)
}

func (o OsFs) ReadFile(name string) ([]byte, error) {
	content, err := os.ReadFile(name)
	if errors.Is(err, os.ErrNotExist) {
		return content, SnapshotNotFound
	}
	return content, err
}

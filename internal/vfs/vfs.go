package vfs

import "errors"

type Vfs interface {
	Exists(name string) (bool, error)
	WriteFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)
}

var SnapshotNotFound = errors.New("snapshot not found")

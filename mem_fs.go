package golden

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MemFs struct {
	files map[string][]byte
}

func NewMemFs() *MemFs {
	return &MemFs{
		files: make(map[string][]byte),
	}
}

func (fs *MemFs) Exists(name string) (bool, error) {
	_, ok := fs.files[name]
	if ok {
		return true, nil
	}
	return false, nil
}

func (fs *MemFs) WriteFile(name string, data []byte) error {
	fs.files[name] = data
	return nil
}

func (fs *MemFs) ReadFile(name string) ([]byte, error) {
	content, ok := fs.files[name]
	if ok {
		return content, nil
	}
	return []byte{}, SnapshotNotFound
}

func AssertContentWasStored(t *testing.T, fs *MemFs, path string, expected []byte) {
	content, ok := fs.files[path]
	assert.Truef(t, ok, "path not found '%s'", path)
	assert.Equal(t, expected, content, "content doesn't match")
}

func AssertSnapshotWasCreated(t *testing.T, fs *MemFs, path string) {
	_, ok := fs.files[path]
	assert.Truef(t, ok, "path not found '%s'", path)
}

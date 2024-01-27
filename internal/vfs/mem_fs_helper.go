package vfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertContentWasStored(t *testing.T, fs *MemFs, path string, expected []byte) {
	content, ok := fs.files[path]
	assert.Truef(t, ok, "path not found '%s'", path)
	assert.Equal(t, expected, content, "content doesn't match")
}

func AssertSnapshotWasCreated(t *testing.T, fs *MemFs, path string) {
	_, ok := fs.files[path]
	assert.Truef(t, ok, "path not found '%s'", path)
}

func AssertSnapShotContains(t *testing.T, fs *MemFs, path string, expected string) {
	snapshot, ok := fs.files[path]
	assert.Truef(t, ok, "path not found '%s'", path)
	assert.Contains(t, string(snapshot), expected)
}

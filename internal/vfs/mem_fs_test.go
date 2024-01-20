package vfs

import (
	"errors"
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemFs(t *testing.T) {
	memFs := NewMemFs()

	// Use memFs to write a file in the path ensuring all things work
	writeFile := func(t *testing.T, path string, content []byte) {
		err := memFs.WriteFile(path, content)
		assert.NoError(t, err)
	}

	t.Run("should write file", func(t *testing.T) {
		filePath := "file.snap"

		content := []byte("some content")
		writeFile(t, filePath, content)

		AssertContentWasStored(t, memFs, filePath, content)
	})

	t.Run("should override existing file", func(t *testing.T) {
		filePath := "override.snap"

		content := []byte("original content")
		writeFile(t, filePath, content)

		lastContent := []byte("new content")
		writeFile(t, filePath, lastContent)

		AssertContentWasStored(t, memFs, filePath, lastContent)
	})

	t.Run("should allow full paths", func(t *testing.T) {
		filePath := "__snapshots/file.snap"

		content := []byte("some content")
		writeFile(t, filePath, content)

		AssertContentWasStored(t, memFs, filePath, content)
	})

	t.Run("should read existing files", func(t *testing.T) {
		filePath := "file_to_read.snap"

		writeFile(t, filePath, []byte("The content we wanted."))

		content, err := memFs.ReadFile(filePath)
		assert.NoError(t, err)

		assert.Equal(t, "The content we wanted.", string(content))
	})

	t.Run("should return error if not found", func(t *testing.T) {
		filePath := "no_existent.snap"

		_, err := memFs.ReadFile(filePath)
		assert.Error(t, err)

		assert.True(t, errors.Is(err, golden.SnapshotNotFound))
	})

	t.Run("should know if file exists", func(t *testing.T) {
		filePath := "file_to_read.snap"

		writeFile(t, filePath, []byte("The content we wanted."))

		exists, err := memFs.Exists(filePath)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should know if file does not exist", func(t *testing.T) {
		exists, err := memFs.Exists("some/file.snap")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

}

package golden_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golden"
	"testing"
)

func TestMemFs(t *testing.T) {
	memFs := golden.NewMemFs()

	t.Run("should write file", func(t *testing.T) {
		filePath := "file.snap"
		content := []byte("some content")
		err := memFs.WriteFile(filePath, content)
		assert.NoError(t, err)
		golden.AssertContentWasStored(t, memFs, filePath, content)
	})

	t.Run("should allow full paths", func(t *testing.T) {
		filePath := "__snapshots/file.snap"
		content := []byte("some content")
		err := memFs.WriteFile(filePath, content)
		assert.NoError(t, err)
		golden.AssertContentWasStored(t, memFs, filePath, content)
	})

	t.Run("should read existing files", func(t *testing.T) {
		filePath := "file_to_read.snap"
		err := memFs.WriteFile(filePath, []byte("The content we wanted."))
		assert.NoError(t, err)
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
}

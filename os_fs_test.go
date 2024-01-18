package golden_test

import (
	"errors"
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestOsFs(t *testing.T) {
	osFs := golden.NewOsFs()

	t.Run("should write file", func(t *testing.T) {
		filePath := "file.snap"

		content := []byte("some content")
		err := osFs.WriteFile(filePath, content)
		assert.NoError(t, err)

		err = os.Remove(filePath)
		assert.NoError(t, err)
	})

	t.Run("should override existing file", func(t *testing.T) {
		filePath := "override.snap"

		content := []byte("original content")
		err := osFs.WriteFile(filePath, content)
		assert.NoError(t, err)

		lastContent := []byte("new content")
		err = osFs.WriteFile(filePath, lastContent)
		assert.NoError(t, err)

		got, err := os.ReadFile(filePath)
		assert.Equal(t, lastContent, got)

		err = os.Remove(filePath)
		assert.NoError(t, err)
	})

	t.Run("should allow full paths", func(t *testing.T) {
		filePath := "__snapshots/file.snap"

		content := []byte("some content")
		err := osFs.WriteFile(filePath, content)
		assert.NoError(t, err)

		err = os.Remove(filePath)
		assert.NoError(t, err)
	})

	t.Run("should read existing files", func(t *testing.T) {
		filePath := "file_to_read.snap"

		err := osFs.WriteFile(filePath, []byte("The content we wanted."))
		assert.NoError(t, err)

		content, err := osFs.ReadFile(filePath)
		assert.NoError(t, err)

		assert.Equal(t, "The content we wanted.", string(content))
		err = os.Remove(filePath)
		assert.NoError(t, err)
	})

	t.Run("should return error if not found", func(t *testing.T) {
		filePath := "no_existent.snap"

		_, err := osFs.ReadFile(filePath)
		assert.Error(t, err)

		assert.True(t, errors.Is(err, golden.SnapshotNotFound))
	})

	t.Run("should know if file exists", func(t *testing.T) {
		filePath := "file_to_read.snap"

		err := osFs.WriteFile(filePath, []byte("The content we wanted."))
		assert.NoError(t, err)

		exists, err := osFs.Exists(filePath)
		assert.NoError(t, err)
		assert.True(t, exists)

		err = os.Remove(filePath)
		assert.NoError(t, err)
	})

	t.Run("should know if file does not exist", func(t *testing.T) {
		exists, err := osFs.Exists("some/file.snap")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

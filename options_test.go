package golden

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions(t *testing.T) {
	t.Run("should configure snapshot", func(t *testing.T) {
		c := Config{name: ""}
		option := Snapshot("a_name")
		option(&c)
		assert.Equal(t, "a_name", c.name)
	})

	t.Run("should configure approval mode", func(t *testing.T) {
		c := Config{approve: false}
		option := WaitApproval()
		option(&c)
		assert.True(t, c.approve)
	})

	t.Run("should configure snapshot folder", func(t *testing.T) {
		c := Config{folder: "testdata"}
		option := Folder("a_folder")
		option(&c)
		assert.Equal(t, "a_folder", c.folder)
	})

	t.Run("should configure extension for snapshot", func(t *testing.T) {
		c := Config{ext: ".snap"}
		option := Extension(".other")
		option(&c)
		assert.Equal(t, ".other", c.ext)
	})
}

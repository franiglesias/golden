package golden

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigMerge(t *testing.T) {
	base := Config{
		folder:  "default_folder",
		name:    "default_name",
		ext:     "default_ext",
		approve: false,
	}

	t.Run("should override folder if present", func(t *testing.T) {
		other := Config{folder: "override"}
		merged := base.merge(other)
		assert.Equal(t, "override", merged.folder)
	})

	t.Run("should not override folder if empty", func(t *testing.T) {
		other := Config{folder: ""}
		merged := base.merge(other)
		assert.Equal(t, "default_folder", merged.folder)
	})

	t.Run("should override name if present", func(t *testing.T) {
		other := Config{name: "override"}
		merged := base.merge(other)
		assert.Equal(t, "override", merged.name)
	})

	t.Run("should not override name if empty", func(t *testing.T) {
		other := Config{name: ""}
		merged := base.merge(other)
		assert.Equal(t, "default_name", merged.name)
	})

	t.Run("should override ext if present", func(t *testing.T) {
		other := Config{ext: "override"}
		merged := base.merge(other)
		assert.Equal(t, "override", merged.ext)
	})

	t.Run("should not override ext if empty", func(t *testing.T) {
		other := Config{ext: ""}
		merged := base.merge(other)
		assert.Equal(t, "default_ext", merged.ext)
	})

}

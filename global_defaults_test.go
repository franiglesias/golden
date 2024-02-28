package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/vfs"
	"testing"
)

func TestGlobalDefaults(t *testing.T) {
	// fs is kind of an in-memory filesystem. This allows us to test the library
	// without polluting the local file system with temporary files. This also allows
	// us to inspect the generated paths and files

	var fs *vfs.MemFs

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the
		// snapshot, otherwise the name won't be accurate

		// Inits a new instance of Golden using an in-memory filesystem that will be
		// empty on each test run

		fs = vfs.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
	}

	t.Run("should use defaults defined folder in all tests", func(t *testing.T) {
		setUp(t)
		golden.Defaults(golden.Folder("__snapshots"))
		golden.Verify(t, "first subject.", golden.Snapshot("example-1"))
		golden.Verify(t, "second subject.", golden.Snapshot("example-2"))
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/example-1.snap")
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/example-2.snap")
	})

	t.Run("should use defaults defined extension in all tests", func(t *testing.T) {
		setUp(t)
		golden.Defaults(golden.Extension(".snapshot"))
		golden.Verify(t, "first subject.", golden.Snapshot("example-1"))
		golden.Verify(t, "second subject.", golden.Snapshot("example-2"))
		vfs.AssertSnapshotWasCreated(t, fs, "testdata/example-1.snapshot")
		vfs.AssertSnapshotWasCreated(t, fs, "testdata/example-2.snapshot")
	})

	t.Run("should not allow set default snapshot name", func(t *testing.T) {
		setUp(t)
		golden.Defaults(golden.Snapshot("example"))
		golden.Verify(t, "example subject.")
		vfs.AssertSnapshotWasCreated(t, fs, "testdata/TestGlobalDefaults/should_not_allow_set_default_snapshot_name.snap")
	})

}

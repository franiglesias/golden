package golden_test

import (
	"fmt"
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	var gld golden.Golden

	// fs is kind of an in-memory filesystem. This allows us to test the library
	// without polluting the local file system with temporary files. This also allows
	// us to inspect the generated paths and files

	var fs *vfs.MemFs

	// tSpy holds a replacement of the standard testing.T. This allows us to separate
	// the simulated test from the test itself.

	var tSpy helper.TSpy

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the
		// snapshot, otherwise the name won't be accurate

		// Inits a new instance of Golden using an in-memory filesystem that will be
		// empty on each test run

		fs = vfs.NewMemFs()
		gld = *golden.NewUsingFs(fs)

		// Replace testing.T with this double to allow spying results

		tSpy = helper.TSpy{
			T: t,
		}
	}

	t.Run("should create snapshot if not exists", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some subject.")
		vfs.AssertSnapshotWasCreated(t, fs, "testdata/TestVerify/should_create_snapshot_if_not_exists.snap")
	})

	t.Run("should write subject as snapshot content", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some output.")
		expected := []byte("some output.")
		vfs.AssertContentWasStored(t, fs, "testdata/TestVerify/should_write_subject_as_snapshot_content.snap", expected)
	})

	t.Run("should not alter snapshot when it already exists", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "some output.")
		gld.Verify(&tSpy, "different output.")

		want := []byte(("some output."))

		vfs.AssertContentWasStored(t, fs, "testdata/TestVerify/should_not_alter_snapshot_when_it_already_exists.snap", want)
	})

	t.Run("should detect and report differences by line", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot for first time
		gld.Verify(&tSpy, "original output.")
		// Changes happened. Verify against existing snapshot
		gld.Verify(&tSpy, "different output.")

		helper.AssertFailedTest(t, &tSpy)
		helper.AssertReportContains(t, &tSpy, "-original output.\n+different output.\n")
	})

	t.Run("should use custom name for snapshot", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "original output", golden.Snapshot("custom_snapshot"))

		vfs.AssertSnapshotWasCreated(t, fs, "testdata/custom_snapshot.snap")
	})

	t.Run("should use default name after spend customized", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "original output", golden.Snapshot("custom_snapshot"))
		gld.Verify(&tSpy, "original output")

		vfs.AssertSnapshotWasCreated(t, fs, "testdata/custom_snapshot.snap")
		vfs.AssertSnapshotWasCreated(t, fs, "testdata/TestVerify/should_use_default_name_after_spend_customized.snap")
	})

	t.Run("should allow external file via custom name", func(t *testing.T) {
		setUp(t)

		// Creates a file in the path, simulating that we put our own
		err := fs.WriteFile("testdata/external_snapshot.snap", []byte("external output"))
		assert.NoError(t, err)

		// By default, golden would create a snapshot. But given that we have a file in
		// the expected path, Golden will use it as criteria, so test should fail given
		// that subject and snapshot doesn't match

		gld.Verify(&tSpy, "generated output", golden.Snapshot("external_snapshot"))
		helper.AssertFailedTest(t, &tSpy)
	})

	t.Run("should scrub data", func(t *testing.T) {
		setUp(t)

		scrubber := golden.NewScrubber("\\d{2}:\\d{2}:\\d{2}.\\d{3}", "<Current Time>")

		// Here we have a non-deterministic subject
		subject := fmt.Sprintf("Current time is: %s", time.Now().Format("15:04:05.000"))

		gld.Verify(&tSpy, subject, golden.WithScrubbers(scrubber))
		helper.AssertPassTest(t, &tSpy)
		vfs.AssertSnapShotContains(t, fs, "testdata/TestVerify/should_scrub_data.snap", "<Current Time>")
	})
}

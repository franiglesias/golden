package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	var fs *golden.MemFs
	var gld golden.Golden
	var tSpy golden.TSpy

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the snapshot

		// Inits a new instance of Golden
		// Avoid using real filesystem in test
		// Inits the fs, so it's empty on each test

		fs = golden.NewMemFs()
		gld = *golden.NewUsingFs(fs)

		// Replace testing.T with this double to allow spying results

		tSpy = golden.TSpy{
			T: t,
		}
	}

	t.Run("should create snapshot if not exists", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some subject.")
		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_create_snapshot_if_not_exists.snap")
	})

	t.Run("should write subject as snapshot content", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some output.")
		expected := []byte(("some output."))
		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_write_subject_as_snapshot_content.snap", expected)
	})

	t.Run("should not alter snapshot when it exists", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "some output.")
		gld.Verify(&tSpy, "different output.")

		want := []byte(("some output."))

		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_not_alter_snapshot_when_it_exists.snap", want)
	})

	t.Run("should detect and report differences by line", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot
		gld.Verify(&tSpy, "original output.")
		// Changes happened. Verify against existing snapshot
		gld.Verify(&tSpy, "different output.")

		golden.AssertFailedTest(t, &tSpy)
		golden.AssertReportContains(t, &tSpy, "-original output.\n+different output.\n")
	})

	t.Run("should use custom name for snapshot", func(t *testing.T) {
		setUp(t)

		gld.UseSnapshot("custom_snapshot").Verify(&tSpy, "original output")

		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/custom_snapshot.snap")
	})

	t.Run("should use default name after spend customized", func(t *testing.T) {
		setUp(t)

		gld.UseSnapshot("custom_snapshot").Verify(&tSpy, "original output")
		gld.Verify(&tSpy, "original output")

		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_use_default_name_after_spend_customized.snap")
	})

	t.Run("should allow external file via custom name", func(t *testing.T) {
		setUp(t)

		// Creates a file in the path, simulating that we put our own
		err := fs.WriteFile("__snapshots/external_snapshot.snap", []byte("external output"))
		assert.NoError(t, err)

		gld.UseSnapshot("external_snapshot").Verify(&tSpy, "generated output")

		// By default, golden would create a snapshot. But given that we have a file in
		// the expected path, Golden will use it as criteria, so test should fail given
		// that subject and snapshot doesn't match

		golden.AssertFailedTest(t, &tSpy)
	})
}

func TestToApprove(t *testing.T) {
	var fs *golden.MemFs
	var gld golden.Golden
	var tSpy golden.TSpy

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the snapshot

		// Inits a new instance of Golden
		// Avoid using real filesystem in test
		// Inits the fs, so it's empty on each test

		fs = golden.NewMemFs()
		gld = *golden.NewUsingFs(fs)

		// Replace testing.T with this double to allow spying results

		tSpy = golden.TSpy{
			T: t,
		}
	}

	t.Run("should create snapshot and fail", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some subject.")
		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/TestToApprove/should_create_snapshot_and_fail.snap")
		golden.AssertFailedTest(t, &tSpy)
	})
}

package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	var fs *vfs.MemFs
	var gld golden.Golden
	var tSpy helper.TSpy

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the snapshot

		// Inits a new instance of Golden
		// Avoid using real filesystem in test
		// Inits the fs, so it's empty on each test

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
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_create_snapshot_if_not_exists.snap")
	})

	t.Run("should write subject as snapshot content", func(t *testing.T) {
		setUp(t)

		gld.Verify(t, "some output.")
		expected := []byte(("some output."))
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_write_subject_as_snapshot_content.snap", expected)
	})

	t.Run("should not alter snapshot when it exists", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "some output.")
		gld.Verify(&tSpy, "different output.")

		want := []byte(("some output."))

		vfs.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_not_alter_snapshot_when_it_exists.snap", want)
	})

	t.Run("should detect and report differences by line", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot
		gld.Verify(&tSpy, "original output.")
		// Changes happened. Verify against existing snapshot
		gld.Verify(&tSpy, "different output.")

		helper.AssertFailedTest(t, &tSpy)
		helper.AssertReportContains(t, &tSpy, "-original output.\n+different output.\n")
	})

	t.Run("should use custom name for snapshot", func(t *testing.T) {
		setUp(t)

		gld.UseSnapshot("custom_snapshot").Verify(&tSpy, "original output")

		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/custom_snapshot.snap")
	})

	t.Run("should use default name after spend customized", func(t *testing.T) {
		setUp(t)

		gld.UseSnapshot("custom_snapshot").Verify(&tSpy, "original output")
		gld.Verify(&tSpy, "original output")

		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_use_default_name_after_spend_customized.snap")
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

		helper.AssertFailedTest(t, &tSpy)
	})
}

func TestToApprove(t *testing.T) {
	var fs *vfs.MemFs
	var gld golden.Golden
	var tSpy helper.TSpy

	setUp := func(t *testing.T) {
		fs = vfs.NewMemFs()
		gld = *golden.NewUsingFs(fs)
		tSpy = helper.TSpy{
			T: t,
		}
	}

	t.Run("should create snapshot and fail", func(t *testing.T) {
		setUp(t)

		gld.ToApprove(&tSpy, "some subject.")
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestToApprove/should_create_snapshot_and_fail.snap")
		helper.AssertFailedTest(t, &tSpy)
	})

	/*
		Simulates the process of running approval tests so snapshot is never taken as
		criteria for matching, meaning that you are waiting for human approval before
		changing test type to Verify
	*/
	t.Run("should update snapshot and fail in second run", func(t *testing.T) {
		setUp(t)

		gld.ToApprove(&tSpy, "starting subject.")
		helper.AssertFailedTest(t, &tSpy)
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestToApprove/should_update_snapshot_and_fail_in_second_run.snap", []byte("starting subject."))
		tSpy.Reset()

		gld.ToApprove(&tSpy, "updated subject.")
		helper.AssertFailedTest(t, &tSpy)
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestToApprove/should_update_snapshot_and_fail_in_second_run.snap", []byte("updated subject."))
		tSpy.Reset()
	})

	/*
		Simulates the process of running approval tests until you obtain approval for
		the generated snapshot

	*/
	t.Run("should accept snapshot at Verify", func(t *testing.T) {
		setUp(t)

		gld.ToApprove(&tSpy, "starting subject.")
		tSpy.Reset()

		// After this run the snapshot will be approved by an expert
		gld.ToApprove(&tSpy, "updated subject.")
		tSpy.Reset()

		// Last snapshot was approved, so we can change the test to Verification
		gld.Verify(&tSpy, "updated subject.")
		helper.AssertPassTest(t, &tSpy)
	})
}

package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"testing"
)

/*
TestToApprove needs the same setup as TestVerify. Check it for documentation.
*/
func TestToApprove(t *testing.T) {
	var gld golden.Golden
	var fs *vfs.MemFs
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

		gld.Verify(&tSpy, "some subject.", golden.WaitApproval())
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestToApprove/should_create_snapshot_and_fail.snap")
		helper.AssertFailedTest(t, &tSpy)
	})

	/*
		Simulates the process of running approval tests so snapshot is never taken as
		criteria for matching, meaning that you are waiting for human approval before
		changing test type to Verify
	*/
	t.Run("should keep test failing while approval mode", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "starting subject.", golden.WaitApproval())
		helper.AssertFailedTest(t, &tSpy)
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestToApprove/should_keep_test_failing_while_approval_mode.snap", []byte("starting subject."))
		tSpy.Reset()

		gld.Verify(&tSpy, "updated subject.", golden.WaitApproval())
		helper.AssertFailedTest(t, &tSpy)
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestToApprove/should_keep_test_failing_while_approval_mode.snap", []byte("updated subject."))
		tSpy.Reset()
	})

	/*
		Simulates the process of running approval tests until you obtain approval for
		the generated snapshot
	*/
	t.Run("should accept snapshot at Verify", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "starting subject.", golden.WaitApproval())
		tSpy.Reset()

		// After this run the snapshot will be approved by an expert

		gld.Verify(&tSpy, "updated subject.", golden.WaitApproval())
		tSpy.Reset()

		// At this point, the snapshot was approved, so we can change the test back to
		// Verification mode, removing the golden.WaitApproval() option

		gld.Verify(&tSpy, "updated subject.")
		helper.AssertPassTest(t, &tSpy)
	})

	/*
		Simulates the process of running approval tests until you obtain approval for
		the generated snapshot, but with custom snapshot file name
	*/
	t.Run("should work with custom snapshot", func(t *testing.T) {
		setUp(t)

		gld.Verify(&tSpy, "starting subject.", golden.Snapshot("approval_snapshot"), golden.WaitApproval())
		tSpy.Reset()

		gld.Verify(&tSpy, "updated subject.", golden.Snapshot("approval_snapshot"), golden.WaitApproval())
		tSpy.Reset()

		gld.Verify(&tSpy, "updated subject.", golden.Snapshot("approval_snapshot"))
		helper.AssertPassTest(t, &tSpy)
	})

	t.Run("should detect and report differences first run", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot for first time
		gld.Verify(&tSpy, "original output.", golden.WaitApproval())

		// Report should show original content as differences
		helper.AssertFailedTest(t, &tSpy)
		helper.AssertReportContains(t, &tSpy, "+original output.\n")
	})

	t.Run("should detect and report differences subsequent run", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot for first time
		gld.Verify(&tSpy, "original output.", golden.WaitApproval())

		// Changes happened. Verify against existing snapshot
		gld.Verify(&tSpy, "different output.", golden.WaitApproval())

		helper.AssertFailedTest(t, &tSpy)
		helper.AssertReportContains(t, &tSpy, "-original output.\n+different output.\n")
	})
}

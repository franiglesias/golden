package golden_test

import (
	"fmt"
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"testing"
	"time"
)

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

		gld.Verify(&tSpy, "some subject.", golden.WaitApproval())
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

		gld.Verify(&tSpy, "starting subject.", golden.WaitApproval())
		helper.AssertFailedTest(t, &tSpy)
		vfs.AssertContentWasStored(t, fs, "__snapshots/TestToApprove/should_update_snapshot_and_fail_in_second_run.snap", []byte("starting subject."))
		tSpy.Reset()

		gld.Verify(&tSpy, "updated subject.", golden.WaitApproval())
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

		gld.Verify(&tSpy, "starting subject.", golden.WaitApproval())
		tSpy.Reset()

		// After this run the snapshot will be approved by an expert
		gld.Verify(&tSpy, "updated subject.", golden.WaitApproval())
		tSpy.Reset()

		// Last snapshot was approved, so we can change the test to Verification
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

		// After this run the snapshot will be approved by an expert
		gld.Verify(&tSpy, "updated subject.", golden.Snapshot("approval_snapshot"), golden.WaitApproval())
		tSpy.Reset()

		// Last snapshot was approved, so we can change the test to Verification
		gld.Verify(&tSpy, "updated subject.", golden.Snapshot("approval_snapshot"))
		helper.AssertPassTest(t, &tSpy)
	})

	t.Run("should scrub data", func(t *testing.T) {
		setUp(t)

		scrubber := golden.NewScrubber("\\d{2}:\\d{2}:\\d{2}.\\d{3}", "<Current Time>")

		// Here we have a non-deterministic subject
		subject := fmt.Sprintf("Current time is: %s", time.Now().Format("15:04:05.000"))

		gld.Verify(&tSpy, subject, golden.WithScrubbers(scrubber))
		helper.AssertPassTest(t, &tSpy)
		vfs.AssertSnapShotContains(t, fs, "__snapshots/TestToApprove/should_scrub_data.snap", "<Current Time>")
	})
}

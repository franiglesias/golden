package golden_test

import (
	"golden"
	"testing"
)

func TestVerify(t *testing.T) {
	var fs *golden.MemFs
	var gt golden.TSpy

	setUp := func(t *testing.T) {
		// Passing t in each setup guarantees that we are using the right name for the snapshot

		// Avoid using real filesystem in test
		// Inits the fs, so it's empty on each test

		fs = golden.NewMemFs()

		// Inits global G for using test filesystem

		golden.G = golden.NewUsingFs(fs)

		// Replace testing.T with this double to allow spying results

		gt = golden.TSpy{
			T: t,
		}
	}

	t.Run("should create snapshot if not exists", func(t *testing.T) {
		setUp(t)

		subject := "some subject."
		golden.Verify(t, subject)
		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_create_snapshot_if_not_exists.snap")
	})

	t.Run("should write subject as snapshot content", func(t *testing.T) {
		setUp(t)

		subject := "some output."
		golden.Verify(t, subject)
		expected := []byte(subject)
		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_write_subject_as_snapshot_content.snap", expected)
	})

	t.Run("should not alter snapshot when it exists", func(t *testing.T) {
		setUp(t)

		subject := "some output."
		golden.Verify(&gt, subject)
		modified := "different output."
		golden.Verify(&gt, modified)

		want := []byte(subject)

		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_not_alter_snapshot_when_it_exists.snap", want)
	})

	t.Run("should detect and report differences by line", func(t *testing.T) {
		setUp(t)

		// Sets the snapshot
		golden.Verify(&gt, "original output.")
		// Changes happened. Verify against existing snapshot
		golden.Verify(&gt, "different output.")

		golden.AssertFailedTest(t, &gt)
		golden.AssertReportContains(t, &gt, "-original output.\n+different output.\n")
	})
}

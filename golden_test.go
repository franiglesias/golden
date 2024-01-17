package golden_test

import (
	"golden"
	"testing"
)

func TestVerify(t *testing.T) {

	// Replace testing.T with a double to be able to spy results
	gt := golden.TSpy{
		T: t,
	}

	t.Run("should create snapshot if not exists", func(t *testing.T) {
		fs := golden.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
		subject := "some subject."
		golden.Verify(t, subject)
		golden.AssertSnapshotWasCreated(t, fs, "__snapshots/TestVerify/should_create_snapshot_if_not_exists.snap")
	})

	t.Run("should write subject as snapshot content", func(t *testing.T) {
		fs := golden.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
		subject := "some output."
		golden.Verify(t, subject)
		expected := []byte(subject)
		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify/should_write_subject_as_snapshot_content.snap", expected)
	})

	t.Run("should not alter snapshot when it exists", func(t *testing.T) {
		fs := golden.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
		subject := "some output."
		golden.Verify(&gt, subject)
		modified := "different output."
		golden.Verify(&gt, modified)
		expected := []byte(subject)
		// When used this way, t.Name() is not used
		golden.AssertContentWasStored(t, fs, "__snapshots/TestVerify.snap", expected)
	})

	t.Run("should detect difference", func(t *testing.T) {
		fs := golden.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
		subject := "original output."
		golden.Verify(&gt, subject)
		modified := "different output."
		golden.Verify(&gt, modified)
		golden.AssertFailedTest(t, &gt)
	})

	t.Run("should report differences", func(t *testing.T) {
		fs := golden.NewMemFs()
		golden.G = golden.NewUsingFs(fs)
		subject := "original output."
		golden.Verify(&gt, subject)
		modified := "different output."
		golden.Verify(&gt, modified)
		golden.AssertReportContains(t, &gt, "-original output.\n+different output.\n")
	})
}

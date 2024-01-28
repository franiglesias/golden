package golden_test

import (
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"testing"
)

/*
TestGlobalMaster needs the same setup as TestVerify. Check it for documentation.
*/
func TestGlobalMaster(t *testing.T) {
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

	t.Run("should create a golden master snapshot", func(t *testing.T) {
		setUp(t)
		f := func(args ...any) any {
			title := args[0].(string)
			part := args[1].(string)
			span := args[2].(int)
			return border(title, part, span)
		}

		titles := []any{"Title", "Subtitle"}
		parts := []any{"*", "#"}
		times := []any{1, 2}

		gld.Master(&tSpy, f, golden.Combine(titles, parts, times))
	})

	t.Run("should manage the error", func(t *testing.T) {
		setUp(t)
		f := func(args ...any) any {
			result, err := division(args[0].(float64), args[1].(float64))
			if err != nil {
				return err.Error()
			}
			return result
		}

		dividend := []any{1.0, 2.0}
		divisor := []any{0.0, -1.0, 1.0, 2.0}

		values := golden.Combine(dividend, divisor)

		gld.Master(&tSpy, f, values)
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestGlobalMaster/should_manage_the_error.snap.json")
		vfs.AssertSnapShotContains(t, fs, "__snapshots/TestGlobalMaster/should_manage_the_error.snap.json", "division by 0")
	})

	t.Run("should support custom name", func(t *testing.T) {
		setUp(t)
		f := func(args ...any) any {
			result, err := division(args[0].(float64), args[1].(float64))
			if err != nil {
				return err.Error()
			}
			return result
		}

		dividend := []any{1.0, 2.0}
		divisor := []any{0.0, -1.0, 1.0, 2.0}

		gld.Master(&tSpy, f, golden.Combine(dividend, divisor), golden.Snapshot("combinations"))
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/combinations.snap.json")
		vfs.AssertSnapShotContains(t, fs, "__snapshots/combinations.snap.json", "division by 0")
	})

	t.Run("should support approval", func(t *testing.T) {
		setUp(t)
		f := func(args ...any) any {
			result, err := division(args[0].(float64), args[1].(float64))
			if err != nil {
				return err.Error()
			}
			return result
		}

		dividend := []any{1.0, 2.0}
		divisor := []any{0.0, -1.0, 1.0, 2.0}

		gld.Master(&tSpy, f, golden.Combine(dividend, divisor), golden.WaitApproval())
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestGlobalMaster/should_support_approval.snap.json")
		vfs.AssertSnapShotContains(t, fs, "__snapshots/TestGlobalMaster/should_support_approval.snap.json", "division by 0")
		helper.AssertFailedTest(t, &tSpy)
	})
}

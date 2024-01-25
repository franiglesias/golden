package golden_test

import (
	"errors"
	"github.com/franiglesias/golden"
	"github.com/franiglesias/golden/internal/helper"
	"github.com/franiglesias/golden/internal/vfs"
	"strings"
	"testing"
)

func TestGoldenMaster(t *testing.T) {
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

		gld.Master(&tSpy, f, golden.Values(titles, parts, times))
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

		values := golden.Values(dividend, divisor)

		gld.Master(&tSpy, f, values)
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/TestGoldenMaster/should_manage_the_error.snap.json")
		vfs.AssertSnapShotContains(t, fs, "__snapshots/TestGoldenMaster/should_manage_the_error.snap.json", "division by 0")
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

		gld.Master(&tSpy, f, golden.Values(dividend, divisor), golden.Snapshot("combinations"))
		vfs.AssertSnapshotWasCreated(t, fs, "__snapshots/combinations.snap.json")
		vfs.AssertSnapShotContains(t, fs, "__snapshots/combinations.snap.json", "division by 0")
	})
}

func border(title string, part string, span int) string {
	width := span*2 + len(title) + 2
	top := strings.Repeat(part, width)
	body := part + strings.Repeat(" ", span) + title + strings.Repeat(" ", span) + part
	return top + "\n" + body + "\n" + top + "\n"
}

func division(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by 0")
	}

	return a / b, nil
}

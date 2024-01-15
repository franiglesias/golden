package golden

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"path"
	"sync"
	"testing"
)

type Golden struct {
	sync.RWMutex
	fs     Vfs
	folder string
	ext    string
	name   string
}

func (g *Golden) Verify(t Failable, s any) {
	g.Lock()
	t.Helper()

	g.name = t.Name()
	name := g.snapshotPath()

	snapshotExists, err := g.fs.Exists(name)
	if err != nil {
		return
	}

	n, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("could not normalize subject %s: %s", n, err)
	}

	if !snapshotExists {
		err = g.fs.WriteFile(name, n)
		if err != nil {
			log.Fatalf("could not create snapshot %s: %s", name, err)
		}
	}

	snapshot, err := g.fs.ReadFile(name)
	if err != nil {
		log.Fatalf("could not read snapshot %s: %s", name, err)
	}

	if string(snapshot) != string(n) {
		t.Errorf("There are differences")
	}

	g.Unlock()
}

func (g *Golden) snapshotPath() string {
	return path.Join(g.folder, g.name+g.ext)
}

var G = New()

func Verify(t Failable, subject any) {
	G.Verify(t, subject)
}

/*
New initializes a new Golden object with defaults. Usually you don't need to
invoke it directly, because it is used to initialize the G var. You may invoke
it when you want to be sure that default settings will be used or to reset G
after using other settings.
*/
func New() *Golden {
	return &Golden{
		folder: "__snapshots",
		ext:    ".snap",
	}
}

func NewUsingFs(fs Vfs) *Golden {
	return &Golden{
		folder: "__snapshots",
		ext:    ".snap",
		fs:     fs,
	}
}

type Failable interface {
	Errorf(format string, args ...any)
	Helper()
	Name() string
}

type TSpy struct {
	*testing.T
	failed bool
}

func (t *TSpy) Errorf(format string, args ...any) {
	t.failed = true
}

func AssertFailedTest(t *testing.T, gt *TSpy) {
	assert.True(t, gt.failed)
}

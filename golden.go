package golden

import (
	"github.com/franiglesias/golden/internal/vfs"
	"log"
	"path"
	"sync"
)

type Golden struct {
	sync.RWMutex
	fs         vfs.Vfs
	normalizer Normalizer
	reporter   DiffReporter
	folder     string
	ext        string
	name       string
	approve    bool
}

type config struct {
	folder  string
	name    string
	ext     string
	approve bool
}

func (c config) snapshotPath(t Failable) string {
	if c.name == "" {
		c.name = t.Name()
	}

	return path.Join(c.folder, c.name+c.ext)
}

func (c config) toApprove() bool {
	return c.approve
}

/*
Verify takes the subject and tries to compare with the content of the snapshot
file. If this file doesn't exist, it creates it.

If the contents of the snapshot and the subject are different, the test fails
and a report of the differences are showed.
*/
func (g *Golden) Verify(t Failable, s any) {
	g.Lock()
	t.Helper()

	// We should separate global configuration and test configuration
	// This way we could start fresh on every run and reset after
	// Also, this could be helpful to have separated global and per test config

	// Global (defaults): path, reporter, ext, normalizer
	// Per test: same as Global, approve mode, name

	subject := g.normalize(s)
	conf := g.testConf()

	name := conf.snapshotPath(t)

	// toApprove mode is like when snapshot doesn't exist, so we have to write it always

	snapshotExists := g.snapshotExists(name)
	if !snapshotExists || conf.toApprove() {
		g.writeSnapshot(name, subject)
	}
	// We should reset the mode, so other test work as expected and you have to explicitly mark as toApprove
	// But not here if we want to do something during reporting

	snapshot := g.readSnapshot(name)
	if snapshot != subject || conf.toApprove() {
		// If Approval add some reminder in the header
		t.Errorf("%s", g.reportDiff(snapshot, subject))
	}

	g.Unlock()
}

func (g *Golden) ToApprove(t Failable, subject any) {
	g.Lock()
	g.approve = true
	g.Unlock()

	g.Verify(t, subject)
}

func (g *Golden) reportDiff(snapshot string, subject string) string {
	return g.reporter.Differences(snapshot, subject)
}

func (g *Golden) normalize(s any) string {
	n, err := g.normalizer.Normalize(s)
	if err != nil {
		log.Fatalf("could not normalize subject %s: %s", n, err)
	}
	return n
}

func (g *Golden) snapshotExists(name string) bool {
	snapshotExists, err := g.fs.Exists(name)
	if err != nil {
		log.Fatalf("could not determine if snahpshot %s exists: %s", name, err)
	}
	return snapshotExists
}

func (g *Golden) writeSnapshot(name string, n string) {
	err := g.fs.WriteFile(name, []byte(n))
	if err != nil {
		log.Fatalf("could not create snapshot %s: %s", name, err)
	}
}

func (g *Golden) readSnapshot(name string) string {
	snapshot, err := g.fs.ReadFile(name)
	if err != nil {
		log.Fatalf("could not read snapshot %s: %s", name, err)
	}
	return string(snapshot)
}

func (g *Golden) snapshotPath(t Failable) string {
	if g.name == "" {
		g.name = t.Name()
	}

	snapshotName := path.Join(g.folder, g.name+g.ext)

	// resets g.name after using it
	g.name = ""

	return snapshotName
}

/*
UseSnapshot allows you custom the name of the snapshot. This can be useful when
you want several snapshot in the same test. Also, it allows you to bring
external files to use as snapshot.

If you don't indicate any name, the snapshot will be named after the test.
*/
func (g *Golden) UseSnapshot(name string) *Golden {
	g.name = name
	return g
}

func (g *Golden) testConf() config {
	c := config{
		folder:  g.folder,
		name:    g.name,
		ext:     g.ext,
		approve: g.approve,
	}
	g.name = ""
	g.approve = false
	return c
}

/*

Global vars and functions

*/

/*
G is a singleton instance of the Golden object. Usually you will not need to instantiate it.
*/
var G = New()

/*
Verify see Golden.Verify

This is a tiny wrapper around the Golden.Verify method.
*/
func Verify(t Failable, subject any) {
	G.Verify(t, subject)
}

/*
UseSnapshot see Golden.UseSnapshot

This is a tiny wrapper around the Golden.UseSnapshot method
*/
func UseSnapshot(name string) *Golden {
	return G.UseSnapshot(name)
}

/*
New initializes a new Golden object with defaults. Usually you don't need to
invoke it directly, because it is used to initialize the G var. You may invoke
it when you want to be sure that default settings will be used or to reset G
after using other settings.
*/
func New() *Golden {
	return &Golden{
		folder:     "__snapshots",
		ext:        ".snap",
		fs:         vfs.NewOsFs(),
		normalizer: JsonNormalizer{},
		reporter:   LineDiffReporter{},
	}
}

/*
NewUsingFs initializes a new Golden object allowing us to change some defaults
from the beginning. Usually for testing purposes only
*/
func NewUsingFs(fs vfs.Vfs) *Golden {
	return &Golden{
		folder:     "__snapshots",
		ext:        ".snap",
		fs:         fs,
		normalizer: JsonNormalizer{},
		reporter:   LineDiffReporter{},
	}
}

/*

Interfaces

*/

/*
Failable interface allows us to replace *testing.T in the own library tests.
*/
type Failable interface {
	Errorf(format string, args ...any)
	Helper()
	Name() string
}

/*
Normalizer normalizes the subject to a string representation that can be compared
*/
type Normalizer interface {
	Normalize(subject any) (string, error)
}

/*
DiffReporter is an interface to represent an object that can show differences
between expected snapshot and subject
*/
type DiffReporter interface {
	Differences(want, got string) string
}

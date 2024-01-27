package golden

import (
	"github.com/franiglesias/golden/internal/combinatory"
	"github.com/franiglesias/golden/internal/vfs"
	"log"
	"sync"
)

/*
Golden is the type that manages snapshotting and test evaluation
*/
type Golden struct {
	sync.RWMutex
	fs         vfs.Vfs
	normalizer Normalizer
	reporter   DiffReporter
	test       Config
	global     Config
}

/*
Verify takes the subject and tries to compare it with the content of the snapshot
file. If this file doesn't exist, it creates it.

If the contents of the snapshot and the subject are different, the test fails
and a report with the differences is showed.
*/
func (g *Golden) Verify(t Failable, s any, options ...Option) {
	g.Lock()
	t.Helper()

	// We should separate global configuration and test configuration
	// This way we could start fresh on every run and reset after
	// Also, this could be helpful to have separated global and per test config

	// Global (defaults): path, reporter, ext, normalizer
	// Per test: same as Global, approve mode, name

	for _, option := range options {
		option(&g.test)
	}

	conf := g.testConfig()
	subject := g.normalize(s, conf)

	name := conf.snapshotPath(t)

	// approval mode works as if the snapshot doesn't exist, so we have to write it always

	snapshotExists := g.snapshotExists(name)
	if !snapshotExists || conf.approvalMode() {
		g.writeSnapshot(name, subject)
	}

	snapshot := g.readSnapshot(name)
	if snapshot != subject || conf.approvalMode() {
		t.Errorf(conf.header(), g.reportDiff(snapshot, subject))
	}

	g.Unlock()
}

/*
ToApprove acts exactly as Verify except that the test never passes waiting for
human approval. This is intentional and the purpose is to remind that you should
review and approve the current snapshot.

When you are totally ok with the snapshot, replace ToApprove with Verify in the test.

Deprecated: use golden.Verify(t, subject, WaitApproval()) instead
*/
func (g *Golden) ToApprove(t Failable, subject any, options ...Option) {
	options = append(options, WaitApproval())
	g.Verify(t, subject, options...)
}

/*
Master generates all combinations of possible values for the parameters of
the subject under test, executes the SUT with all those combinations,
accumulates the outputs, and creates a snapshot of that using Verify internally.

You need to pass a wrapper function that executes the subject under test and
returns a string representation of its output. This wrapper function receives any
number of parameters of any type. It's up to you to cast or convert these
parameters in something that can be managed by the subject under test.

Also, is up to you to capture the output of the SUT as a string.

The parameters received by the wrapper function are the result of combining all
the possible values for each parameter that you would pass to the SUT. This will
create a lot of tests (tenths or hundredths).
*/
func (g *Golden) Master(t Failable, f combinatory.Wrapper, values [][]any, options ...Option) {
	g.global.ext = ".snap.json"
	subject := combinatory.Master(f, values...)
	g.Verify(t, subject, options...)
}

func (g *Golden) reportDiff(snapshot string, subject string) string {
	return g.reporter.Differences(snapshot, subject)
}

func (g *Golden) normalize(s any, conf Config) string {
	n, err := g.normalizer.Normalize(s)
	if err != nil {
		log.Fatalf("could not normalize subject %s: %s", n, err)
	}
	for _, scrubber := range conf.scrubbers {
		n = scrubber.Clean(n)
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

/*
UseSnapshot allows you to customize the name of the snapshot. This can be useful
when you want several snapshots in the same test. Also, it allows you to bring
your own files to use them as snapshots.

If you don't set any name, the snapshot will be named after the test.

Deprecated: use golden.Verify(t, subject, golden.Snapshot("snapshot_name")) instead
*/
func (g *Golden) UseSnapshot(name string) *Golden {
	Snapshot(name)
	return g
}

func (g *Golden) testConfig() Config {
	c := g.global.merge(g.test)
	g.test = Config{}
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

TL;DR Verify the subject against a snapshot

This is a tiny wrapper around the Golden.Verify method.
*/
func Verify(t Failable, subject any, options ...Option) {
	G.Verify(t, subject, options...)
}

/*
ToApprove see Golden.ToApprove

TL;DR Updates a snapshot until someone approves it

This is a tiny wrapper around the Golden.ToApprove method.

Deprecated: use golden.Verify(t, subject, golden.WaitApproval()) instead
*/
func ToApprove(t Failable, subject any) {
	G.Verify(t, subject, WaitApproval())
}

/*
Master see Golden.Master

TL;DR Generates and executes SUT with all possible combinations of values

This is a tiny wrapper around the Golden.Master method.
*/
func Master(t Failable, f combinatory.Wrapper, values [][]any, options ...Option) {
	G.Master(t, f, values, options...)
}

/*
UseSnapshot see Golden.UseSnapshot

# This is a tiny wrapper around the Golden.UseSnapshot method

Deprecated: use golden.Verify(t, subject , golden.Snapshot("snapshot_name") instead
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
	return NewUsingFs(vfs.NewOsFs())
}

/*
NewUsingFs initializes a new Golden object allowing us to change some defaults
from the beginning. Usually for testing purposes only
*/
func NewUsingFs(fs vfs.Vfs) *Golden {
	g := Config{
		folder:  "__snapshots",
		name:    "",
		ext:     ".snap",
		approve: false,
	}
	return &Golden{
		global:     g,
		test:       Config{},
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

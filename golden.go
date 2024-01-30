package golden

import (
	"github.com/franiglesias/golden/internal/combinatory"
	"github.com/franiglesias/golden/internal/vfs"
	"log"
	"sync"
)

const approvalHeader = "**Approval mode**: Remove WaitApproval() when you are happy with this snapshot.\n%s"
const verifyHeader = "**Verify mode**\n%s"

/*
Golden is the type that manages snapshotting and test evaluation
*/
type Golden struct {
	sync.RWMutex
	fs         vfs.Vfs
	normalizer Normalizer
	reporter   DiffReporter
	global     Config
}

/*
Verify takes the subject and tries to compare it with the content of the snapshot
file. If this file doesn't exist, it creates it.

If the contents of the snapshot and the subject are different, the test fails
and a report with the differences is showed.
*/func (g *Golden) Verify(t Failable, s any, options ...Option) {
	g.Lock()
	t.Helper()

	conf := g.global
	for _, option := range options {
		option(&conf)
	}

	subject := g.normalize(s, conf.scrubbers)

	name := conf.snapshotPath(t)

	if conf.approvalMode() {
		g.approvalFlow(t, name, subject)
	} else {
		g.verifyFlow(t, name, subject)
	}

	g.Unlock()
}

func (g *Golden) approvalFlow(t Failable, name string, subject string) {
	var previous string
	if g.snapshotExists(name) {
		previous = g.readSnapshot(name)
	}

	g.writeSnapshot(name, subject)

	t.Errorf(approvalHeader, g.reportDiff(previous, subject))
}

func (g *Golden) verifyFlow(t Failable, name string, subject string) {
	if !g.snapshotExists(name) {
		g.writeSnapshot(name, subject)
	}

	snapshot := g.readSnapshot(name)

	if snapshot != subject {
		t.Errorf(verifyHeader, g.reportDiff(snapshot, subject))
	}
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

func (g *Golden) normalize(s any, scrubbers []Scrubber) string {
	n, err := g.normalizer.Normalize(s)
	if err != nil {
		log.Fatalf("could not normalize subject %s: %s", n, err)
	}
	for _, scrubber := range scrubbers {
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
Master see Golden.Master

TL;DR Generates and executes SUT with all possible combinations of values

This is a tiny wrapper around the Golden.Master method.
*/
func Master(t Failable, f combinatory.Wrapper, values [][]any, options ...Option) {
	G.Master(t, f, values, options...)
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
	return &Golden{
		global: Config{
			folder:  "__snapshots",
			name:    "",
			ext:     ".snap",
			approve: false,
		},
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

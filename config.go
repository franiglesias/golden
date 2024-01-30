package golden

import "path"

type Config struct {
	folder    string
	name      string
	ext       string
	approve   bool
	scrubbers []Scrubber
}

func (c Config) snapshotPath(t Failable) string {
	if c.name == "" {
		c.name = t.Name()
	}

	return path.Join(c.folder, c.name+c.ext)
}

func (c Config) approvalMode() bool {
	return c.approve
}

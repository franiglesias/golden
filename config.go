package golden

import "path"

type Config struct {
	folder  string
	name    string
	ext     string
	approve bool
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

func (c Config) merge(other Config) Config {
	if len(other.name) != 0 {
		c.name = other.name
	}
	if len(other.ext) != 0 {
		c.ext = other.ext
	}

	if len(other.folder) != 0 {
		c.folder = other.folder
	}
	if other.approve == true {
		c.approve = other.approve
	}

	return c
}

func (c Config) header() string {
	if c.approvalMode() {
		return "**Approval mode**: Remove WaitApproval() when you are happy with this snapshot.\n%s"
	}
	return "**Verify mode**\n%s"
}

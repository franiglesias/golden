package golden

import "path"

type Config struct {
	folder  string
	name    string
	ext     string
	approve bool
}

func Options() Config {
	return Config{}
}

func (c Config) UseSnapshot(name string) Config {
	c.name = name
	return c
}

func (c Config) snapshotPath(t Failable) string {
	if c.name == "" {
		c.name = t.Name()
	}

	return path.Join(c.folder, c.name+c.ext)
}

func (c Config) toApprove() bool {
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

package golden

type Option func(g *Config) Option

/*
Snapshot allows to pass a string to be used as file name of the current snapshot
*/
func Snapshot(name string) Option {
	return func(c *Config) Option {
		previous := c.name
		c.name = name
		return func(c *Config) Option {
			return Snapshot(previous)
		}
	}
}

/*
WaitApproval will execute this test in Approval Mode, so the snapshot will be
updated but the test will not pass. To make the test pass, remove this option
*/
func WaitApproval() Option {
	return func(c *Config) Option {
		c.approve = true
		return verifyMode()
	}
}

/*
verifyMode will return to verify Mode. It is used only internally to reset the WaitApproval option
*/
func verifyMode() Option {
	return func(c *Config) Option {
		c.approve = false
		return WaitApproval()
	}
}

func WithScrubbers(scrubbers ...Scrubber) Option {
	return func(c *Config) Option {
		c.scrubbers = scrubbers
		return WithScrubbers()
	}
}

/*
Folder configure a folder to store the snapshot

Can be useful if you want to join several snapshot under some criteria
*/
func Folder(folder string) Option {
	return func(c *Config) Option {
		previous := c.folder
		c.folder = folder
		return func(c *Config) Option {
			return Folder(previous)
		}
	}
}

/*
Extension configure and extension for the snapshot file

Can be useful if you want to manipulate the snapshot file based on type
*/
func Extension(extension string) Option {
	return func(c *Config) Option {
		previous := c.ext
		c.ext = extension
		return func(c *Config) Option {
			return Extension(previous)
		}
	}
}

/*
Combine is a convenience function that wraps the values you pass to golden.Master() tests.

This way, the values can be received as a single parameter allowing Master to use options
*/
func Combine(v ...[]any) [][]any {
	return v
}

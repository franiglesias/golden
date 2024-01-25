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

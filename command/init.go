package command

type InitCommand struct {
	Meta
}

func (c *InitCommand) Run(args []string) int {
	return 1
}

func (c *InitCommand) Synopsis() string {
	return "Initializes the Appfile"
}

func (c *InitCommand) Help() string {
	return "I need somebody..."
}

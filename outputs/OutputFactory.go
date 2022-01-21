package outputs

func GetOutput(name string) OutputInterface {
	switch name {
	case "stdout":
		return &StdoutOutput{}
	default:
		return &StdoutOutput{}
	}
}

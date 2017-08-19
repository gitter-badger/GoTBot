package fail

import "fmt"

type CommandNotFound struct {
	Name string
}

func (e *CommandNotFound) Error() string {
	return fmt.Sprintf("No command defined with name '%s'", e.Name)
}

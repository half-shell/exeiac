package infra

import "fmt"

type RoomError struct {
	id     string
	path   string
	reason string
	trace  error
}

func (e RoomError) Error() string {
	return fmt.Sprintf("! Error%s:room: %s: %s\n< %s", e.id,
		e.reason, e.path, e.trace.Error())
}

type ErrBrickNotFound struct {
	brick string
}

func (e ErrBrickNotFound) Error() string {
	return fmt.Sprintf("Brick not found: %s", e.brick)
}

type ActionNotImplementedError struct {
	Action string
	Module *Module
}

func (err ActionNotImplementedError) Error() string {
	return fmt.Sprintf("Module %s does not implement action %s", err.Module.Name, err.Action)
}

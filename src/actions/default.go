package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

// Triggers a module execution for very single brick in `bricksToExecute`
// Ignores errors and calls the action in `args.Action` for every single brick,
// then prints out a summary of it all.
// Exit code matches 3 if an error occured, 0 otherwise.
func Default(
	infra *exinfra.Infra,
	args *exargs.Arguments,
	bricksToExecute exinfra.Bricks) (statusCode int, err error) {
	if infra == nil && args == nil {
		err = exargs.ErrBadArg{Reason: "Error: infra and args are not set"}

		return
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		report := ExecReport{Brick: b}
		statusCode, err = b.Module.Exec(b, args.Action, args.OtherOptions, []string{})

		if err != nil {
			if _, is := err.(exinfra.ActionNotImplementedError); is {
				// NOTE(half-shell): if action if not implemented, we don't take it as an error
				// and move on with the execution
				fmt.Printf("%v ; assume there is nothing to do.\n", err)
				err = nil
				report.Status = "OK"
			} else {
				statusCode = 3
				report.Status = "ERR"
				report.Error = err
			}
		} else {
			report.Status = "DONE"
		}

		execSummary[i] = report
	}

	execSummary.Display()

	return
}

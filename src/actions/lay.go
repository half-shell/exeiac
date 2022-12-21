package actions

import (
	"bytes"
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Lay(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for lay action"}

		return 3, err
	}

	if conf.Interactive {
		fmt.Println("Here, the bricks list to lay :")
		fmt.Print(bricksToExecute)

		// NOTE(half-shell): We might change this behavior to only ask for a "\n" input
		// instead of a Y/N choice.
		confirm, err := extools.AskConfirmation("\nDo you want to continue ?")

		if err != nil {
			return 3, err
		} else if !confirm {
			return 0, nil
		}
	}

	err = enrichDatas(bricksToExecute, infra)
	if err != nil {
		return 3, err
	}

	skipFollowing := false
	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

		// skip if an error was encounter before
		if skipFollowing {
			report.Status = TAG_SKIP
			execSummary[i] = report
			continue
		}

		// write env file if needed
		envs, err := writeEnvFilesAndGetEnvs(b)
		if err != nil {

			return 3, err
		}

		// TODO(half-shell): Work around to avoid polluting conf's OtherOptions.
		// Ideally, we would have a flexible way of providing a "non-interactive" flag
		// to a module.
		args := conf.OtherOptions
		if !conf.Interactive {
			args = append(args, "--non-interactive")
		}

		exitStatus, err := b.Module.Exec(b, "lay", args, envs)
		if err != nil {
			skipFollowing = true
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = 3
		} else if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("lay return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}

		// check if outputs has changed
		stdout := exinfra.StoreStdout{}
		exitStatus, err = b.Module.Exec(b, "output", []string{}, envs, &stdout)
		if err != nil {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output failed : %v", err)
			report.Status = TAG_ERROR
			statusCode = 3
		}
		if exitStatus != 0 {
			skipFollowing = true
			report.Error = fmt.Errorf("layed apparently success but output return : %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}
		if bytes.Compare(stdout.Output, b.Output) == 0 {
			report.Status = TAG_NO_CHANGE
		} else {
			b.Output = stdout.Output
			report.Status = TAG_DONE
		}

		execSummary[i] = report
	}

	execSummary.Display()
	return
}

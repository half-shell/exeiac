package actions

import (
	"fmt"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func Plan(
	infra *exinfra.Infra,
	conf *exargs.Configuration,
	bricksToExecute exinfra.Bricks,
) (
	statusCode int,
	err error,
) {
	if len(bricksToExecute) == 0 {
		err = exinfra.ErrBadArg{Reason: "Error: you should specify at least a brick for plan action"}

		return 3, err
	}

	err = enrichDatas(bricksToExecute, infra)
	if err != nil {
		return 3, err
	}

	execSummary := make(ExecSummary, len(bricksToExecute))

	for i, b := range bricksToExecute {
		extools.DisplaySeparator(b.Name)
		report := ExecReport{Brick: b}

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

		exitStatus, err := b.Module.Exec(b, "plan", args, envs)
		if err != nil {
			report.Error = err
			report.Status = TAG_ERROR
			statusCode = 3
		} else if exitStatus == 0 {
			report.Status = TAG_NO_CHANGE
		} else if exitStatus == 1 {
			report.Status = TAG_DRIFT
			if statusCode == 0 {
				statusCode = 1
			}
		} else {
			report.Error = fmt.Errorf("plan return: %d", exitStatus)
			report.Status = TAG_ERROR
			statusCode = 3
		}

		execSummary[i] = report

	}

	execSummary.Display()

	return
}

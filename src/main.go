package main

import (
	"fmt"
	"log"
	"os"
	exaction "src/exeiac/actions"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"
)

func main() {
	var statusCode int
	// get arguments
	args, err := exargs.GetArguments()
	if err != nil {
		fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get arguments\n",
			err)
		os.Exit(2)
	}

	// build infra representation
	infra, err := exinfra.CreateInfra(args.Rooms, args.Modules)
	if err != nil {
		fmt.Printf("%v\n> Error636f6894:main/main: "+
			"unable to get an infra representation\n", err)
		os.Exit(1)
	}

	// valid arguments (arg.brickNames are in infra.Bricks...)
	err = validArgBricksAreInInfra(&infra, &args)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(arthur91f): func getBricksToExecute(args.BricksNames args.Specifier)
	var bricksToExecute []string
	bricksToExecute = args.BricksNames

	// enrich bricks that we will execute
	enrichBricks(&infra)

	// executeAction
	// if args.action is in the list do that else use otherAction
	if behaviour, ok := exaction.BehaviourMap[args.Action]; ok {
		statusCode, err = behaviour(&infra, &args, bricksToExecute)
	} else {
		statusCode, err = exaction.BehaviourMap["default"](&infra, &args, bricksToExecute)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	os.Exit(statusCode)
}

var availableBricksSpecifiers = []string{
	"linked_previous", "all_previous", "lp", "ap",
	"direct_previous", "dp",
	"selected", "s",
	"direct_next", "dn",
	"linked_next", "all_next", "ln", "an"}

func validArgBricksAreInInfra(infra *exinfra.Infra, args *exargs.Arguments) error {
	// valid that args.BricksNames items are valid
	for _, arg := range args.BricksNames {
		if _, ok := infra.Bricks[arg]; !ok {
			return exargs.ErrBadArg{Reason: "Brick doesn't exist:", Value: arg}
		}
	}

	// valid BricksSpecifiers
	for _, specifier := range args.BricksSpecifiers {
		if !extools.ContainsString(availableBricksSpecifiers, specifier) {
			return exargs.ErrBadArg{Reason: "Brick's specifier doesn't exist:",
				Value: specifier}
		}
	}
	return nil
}

func enrichBricks(infra *exinfra.Infra) {
	// TODO(arthur91f): remove log.Fatal
	for _, b := range infra.Bricks {
		if b.IsElementary {
			conf, err := exinfra.BrickConfYaml{}.New(b.ConfigurationFilePath)
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}

			err = b.Enrich(conf, infra)
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}
			err = b.Module.LoadAvailableActions()
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}
		}
	}
}

/*func getBricksToExecute(infra *exinfra.Infra, args *exargs.Arguments) (
	bricksToExecute []string, err error) {

	var bricks []*exinfra.Brick
	for _, brickName := range args.BricksNames {
		bricks = append(bricks, infra.Bricks[brickName])
	}

	var bricksSpecified []*exinfra.Brick
	var bricksToAdd *exinfra.Brick
	for _, specifier := range args.BricksSpecifiers {
		switch specifier {
		case "linked_previous", "all_previous", "lp", "ap":
			for _, brick := range bricks {
				bricksToAdd = exinfra.GetLinkedPrevious(infra, brickName)
			}
		case "direct_previous", "dp":
			bricksToAdd = exinfra.GetDirectPrevious(infra, brickName)
		case "selected", "s":
			bricksToAdd = exinfra.GetLinkedPrevious()
		case "direct_next", "dn":
			bricksToAdd = exinfra.GetLinkedPrevious()
		case "linked_next", "all_next", "ln", "an":
			bricksToAdd = exinfra.GetLinkedPrevious()
		default:
			err = exargs.ErrBadArg{Reason: "Brick's specifier doesn't exist:",
				Value: specifier}
			return
		}
		bricksSpecified = append(bricksSpecified, bricksToAdd)
	}

	sort.Sort(bricks)

	return
}*/

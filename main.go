package main

import (
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/executil"
	"github.com/kovetskiy/lorg"
	"github.com/seletskiy/hierr"
)

const version = "1.0"

const usage = `orgalorg-bottlebreaker - based on gunter sync tool for orgalorg.

orgalorg-bottlebreaker will run three tools in sequence:

    * gunter to sync files from current working directory into specified
      destination root.
      gunter will generate changed files log, that will be passed further;

    * treetrunks to remove files from destination root that are not presents
      in the current working directory.
      Removed files will be appended to the changed files log.

    * guntalina will receive changed files log and run specified actions;

Tool expects valid orgalorg sync protocol to be received on stdin.

Usage:
    orgalorg-bottlebreaker -h | --help
    orgalorg-bottlebreaker --check-deps
    orgalorg-bottlebreaker [options] --sync

Options:
    -h --help                   Show this help.
    -r --root <root>            Destination root directory to synchronize with.
                                 [default: /]
    -n --dry-run                Run all tools in dry-run mode.
    -f --force                  Run all tools with force mode on.
    -b --backup <dir>           Store backups in the specified directory.
    -c --actions-config <path>  Specify config path for guntalina (on the
                                 remote host).
    --sync                      Run sync tools as described.
    --check-deps                Check dependencies and exit.
`

var (
	exit = os.Exit

	logger = lorg.NewLog()
)

func main() {
	args, err := docopt.Parse(
		usage,
		nil,
		true,
		"orgalorg-bottlebreaker "+version,
		false,
	)
	if err != nil {
		panic(err)
	}

	dependencies := []string{`gunter`, `guntalina`, `treetrunks`}

	err = checkDependencies(dependencies)
	if err != nil {
		logger.Fatalf(
			"%s",
			hierr.Errorf(
				err,
				`can't find all dependencies: %v`,
				dependencies,
			),
		)
	}

	switch {
	case args["--check-deps"].(bool):
		return

	}
}

func checkDependencies(dependencies []string) error {
	for _, bin := range dependencies {
		_, _, err := executil.Run(exec.Command(bin, `--version`))
		if err != nil {
			if !executil.IsExitError(err) {
				return hierr.Errorf(
					err.(*executil.Error).RunErr,
					`can't find: '%s'`,
					bin,
				)
			}

			return hierr.Errorf(
				err,
				`unexpected error while checking dependency: '%s'`,
				bin,
			)
		}
	}

	return nil
}

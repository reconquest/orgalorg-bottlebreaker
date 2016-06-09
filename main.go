package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/executil"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/go-prefixwriter"
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
    orgalorg-bottlebreaker [options] [--sync]

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

	logger.SetLevel(lorg.LevelDebug)

	switch {
	case args["--check-deps"].(bool):
		return

	default:
		sync(args)
	}
}

func checkDependencies(dependencies []string) error {
	for _, bin := range dependencies {
		_, _, err := evaluateCommandWithStdin(nil, bin, `--version`)
		if err != nil {

			return hierr.Errorf(
				err,
				`unexpected error while checking dependency: '%s'`,
				bin,
			)
		}
	}

	return nil
}

func sync(args map[string]interface{}) error {
	var (
		rootDir = args["--root"].(string)
		dryRun  = args["--dry-run"].(bool)
		force   = args["--force"].(bool)

		actionsConfig, _ = args["--actions-config"].(string)
	)

	changedFiles, err := runGunter(rootDir, dryRun)
	if err != nil {
		return hierr.Errorf(
			err,
			`failure during running gunter`,
		)
	}

	err = runGuntalina(changedFiles, dryRun, force, actionsConfig)
	if err != nil {
		return hierr.Errorf(
			err,
			`failure during running guntalina`,
		)
	}

	return nil
}

func runGunter(
	rootDir string,
	dryRun bool,
) ([]string, error) {
	command := []string{`gunter`, `-d`, rootDir, `-l`, `/dev/stdout`}

	if dryRun {
		command = append(command, `-r`)
	}

	stdout, stderr, err := evaluateCommandWithStdin(nil, command...)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`can't run gunter`,
		)
	}

	if stderr != "" {
		fmt.Fprint(prefixwriter.New(os.Stderr, `{gunter} `), stderr)
	}

	return strings.Split(stdout, "\n"), nil
}

func runGuntalina(
	changedFiles []string,
	dryRun bool,
	force bool,
	actionsConfig string,
) error {
	command := []string{`guntalina`, `-s`, `/dev/stdin`}

	if dryRun {
		command = append(command, `-r`)
	}

	if force {
		command = append(command, `-f`)
	}

	if actionsConfig != "" {
		command = append(command, `-c`, actionsConfig)
	}

	stdout, stderr, err := evaluateCommandWithStdin(
		bytes.NewBufferString(strings.Join(changedFiles, "\n")),
		command...,
	)
	if err != nil {
		return hierr.Errorf(
			err,
			`can't run guntalina`,
		)
	}

	if stdout != "" {
		fmt.Fprint(prefixwriter.New(os.Stdout, `{guntalina} `), stdout)
	}

	if stderr != "" {
		fmt.Fprint(prefixwriter.New(os.Stderr, `{guntalina} `), stderr)
	}

	return nil
}

func evaluateCommandWithStdin(
	stdin io.Reader,
	args ...string,
) (string, string, error) {
	command := exec.Command(args[0], args[1:]...)

	if stdin != nil {
		command.Stdin = stdin
	}

	logger.Debugf(`running command: %v`, args)

	stdout, stderr, err := executil.Run(command)
	if err != nil {
		if executil.IsExitError(err) {
			return string(stdout), string(stderr), hierr.Errorf(
				err,
				`command exited with non-zero exit code: %d`,
				executil.GetExitStatus(err),
			)
		}

		return string(stdout), string(stderr), hierr.Errorf(
			err.(*executil.Error).RunErr,
			`unexpected error while running command`,
		)
	}

	return string(stdout), string(stderr), nil
}

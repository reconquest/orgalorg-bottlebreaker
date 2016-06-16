package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kovetskiy/executil"
	"github.com/reconquest/go-prefixwriter"
	"github.com/seletskiy/hierr"
)

func runGunter(
	rootDir string,
	backupDir string,
	dryRun bool,
) ([]string, error) {
	command := []string{`gunter`, `-d`, rootDir, `-l`, `/dev/stdout`}

	if dryRun {
		command = append(command, `-r`)
	}

	if backupDir != "" {
		command = append(command, `-b`, backupDir)
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

	stdout = strings.Trim(stdout, "\n")
	if stdout == "" {
		return []string{}, nil
	}

	return strings.Split(stdout, "\n"), nil
}

func runTreetrunks(rootDir string, dryRun bool) ([]string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`can't get working directory`,
		)
	}

	command := []string{`treetrunks`, workDir, rootDir}

	if dryRun {
		command = append(command, `-n`)
	}

	stdout, stderr, err := evaluateCommandWithStdin(nil, command...)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`cant't run treetrunks`,
		)
	}

	if stderr != "" {
		fmt.Fprint(prefixwriter.New(os.Stderr, `{treetrunks} `), stderr)
	}

	stdout = strings.Trim(stdout, "\n")
	if stdout == "" {
		return []string{}, nil
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

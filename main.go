package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
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
    -m --sync-mode <mode>       Specify number of received SYNCs should be
                                 received prior continuing to next step.
                                 * all - all nodes should acknowledge;
                                 [default: all]
    -a --alone                  Do not expect any SYNC protocol on stdin.
    -v --verbose                Print debug information.
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
		logger.Errorf(
			"%s",
			hierr.Errorf(
				err,
				`can't find all dependencies: %v`,
				dependencies,
			),
		)

		exit(1)
	}

	logger.SetLevel(lorg.LevelInfo)
	if args["--verbose"].(bool) {
		logger.SetLevel(lorg.LevelDebug)
	}

	switch {
	case args["--check-deps"].(bool):
		return

	default:
		err = handleSync(args)
	}

	if err != nil {
		logger.Error(err)

		exit(1)
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

func parseSyncMode(line string) (syncMode, error) {
	switch line {
	case "all":
		return syncModeAll, nil
	}

	return syncModeUnknown, fmt.Errorf(
		`unknown sync mode given: '%s'`,
		line,
	)
}

func handleSync(args map[string]interface{}) error {
	var (
		rootDir      = args["--root"].(string)
		dryRun       = args["--dry-run"].(bool)
		force        = args["--force"].(bool)
		rawSyncMode  = args["--sync-mode"].(string)
		backupDir, _ = args["--backup"].(string)

		alone = args["--alone"].(bool)

		actionsConfig, _ = args["--actions-config"].(string)
	)

	syncMode, err := parseSyncMode(rawSyncMode)
	if err != nil {
		return hierr.Errorf(
			err,
			`can't parse sync mode`,
		)
	}

	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	prefix, err := receiveProtocolPrefix(reader)
	if err != nil && !alone {
		return hierr.Errorf(
			err,
			`can't receive protocol prefix`,
		)
	}

	nodes, err := receiveNodesList(reader)
	if err != nil && !alone {
		return hierr.Errorf(
			err,
			`can't receive nodes list`,
		)
	}

	logger.Infof(`{protocol} syncing with %d nodes`, len(nodes))

	changedFiles, err := runGunter(rootDir, backupDir, dryRun)
	if err != nil {
		return hierr.Errorf(
			err,
			`failure during running gunter`,
		)
	}

	logger.Infof(`{gunter} %d files changed`, len(changedFiles))

	err = synchronize(reader, writer, "gunter", prefix, nodes, syncMode)
	if err != nil && !alone {
		return hierr.Errorf(
			err,
			`can't synchronize with other nodes [after gunter]`,
		)
	}

	removedFiles, err := runTreetrunks(rootDir, dryRun)
	if err != nil {
		return hierr.Errorf(
			err,
			`failure during running treetrunks`,
		)
	}

	logger.Infof(`{treetrunks} %d files removed`, len(removedFiles))

	err = synchronize(reader, writer, "treetrunks", prefix, nodes, syncMode)
	if err != nil && !alone {
		return hierr.Errorf(
			err,
			`can't synchronize with other nodes [after treetrunks]`,
		)
	}

	err = runGuntalina(
		append(changedFiles, removedFiles...),
		dryRun,
		force,
		actionsConfig,
	)
	if err != nil {
		return hierr.Errorf(
			err,
			`failure during running guntalina`,
		)
	}

	return nil
}

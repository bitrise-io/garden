package cli

import "github.com/codegangsta/cli"

const (
	// --- Standard app flags

	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// --- Other app flags

	// ZoneKey ...
	ZoneKey = "zone"
	// PlantKey ...
	PlantKey = "plant"
)

var (
	commands = []cli.Command{
		{
			Name:   "grow",
			Usage:  "Grow your plants!",
			Action: grow,
		},
		{
			Name:   "reap",
			Usage:  "Use your plants!",
			Action: reap,
		},
		{
			Name:   "view",
			Usage:  "View your plants!",
			Action: view,
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
		cli.StringFlag{
			Name:  ZoneKey,
			Value: "",
			Usage: "Zone filter: work only on plants which belong to this zone",
		},
		cli.StringFlag{
			Name:  PlantKey,
			Value: "",
			Usage: "Plant ID filter: work only on the specified plant",
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}

package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/codegangsta/cli"
)

var (
	// WorkWithZone ...
	WorkWithZone string
	// WorkWithPlantID ...
	WorkWithPlantID string
)

func before(c *cli.Context) error {
	// Log level
	if logLevel, err := log.ParseLevel(c.String(LogLevelKey)); err != nil {
		log.Fatal("Failed to parse log level:", err)
	} else {
		log.SetLevel(logLevel)
	}

	WorkWithPlantID = c.String(PlantKey)
	if WorkWithPlantID == "" {
		WorkWithZone = c.String(ZoneKey)
		if WorkWithZone != "" {
			log.Infoln(" =>", colorstring.Blue("Working with Zone"), ":", WorkWithZone)
		}
	} else {
		log.Infoln(" =>", colorstring.Yellow("Working with Plant"), ":", WorkWithPlantID)
	}

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run the Envman CLI.
func Run() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "garden"
	app.Version = "0.9.1"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = appFlags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Finished with error:", err)
	}
}

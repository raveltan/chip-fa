package main

import (
	"log"
	"os"

	"github.com/raveltan/chip-fa/emulator"
	"github.com/urfave/cli/v2"
)

func main() {
	var scaling float64
	var cyclePerSecond int
	var romFile string
	var isDebugging bool
	var hdpiScale float64

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "prints the current version of the chip-fa",
	}
	app := &cli.App{
		Name:    "Chip-Fa",
		Usage:   "Chip8's emulator written in GO",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases:     []string{"r"},
				Name:        "rom",
				Usage:       "`PATH` to ROM file that will be run on the Chip8's emulator",
				Required:    true,
				Destination: &romFile,
			},
			&cli.Float64Flag{
				Aliases:     []string{"s"},
				Name:        "scale",
				Value:       1,
				Usage:       "Window size scaling",
				Destination: &scaling,
			},
			&cli.Float64Flag{
				Aliases:     []string{"x"},
				Name:        "hdi-scale",
				Value:       0,
				Usage:       "HDPI Pixel Scaling",
				Destination: &hdpiScale,
			},
			&cli.IntFlag{
				Aliases:     []string{"c"},
				Name:        "cycle",
				Value:       60,
				Usage:       "Cycle per second for the CPU emulation ",
				Destination: &cyclePerSecond,
			},
			&cli.BoolFlag{
				Aliases:     []string{"d"},
				Name:        "debug",
				Value:       false,
				Usage:       "Enable debugging",
				Destination: &isDebugging,
			},
		},
		Action: func(c *cli.Context) error {
			emulator.StartEmulation(romFile, hdpiScale, scaling, cyclePerSecond, isDebugging)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

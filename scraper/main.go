package main

import (
	"fmt"
	"github.com/gookit/goutil/dump"
	"github.com/jessevdk/go-flags"
	"github.com/mfojtik/euclid/scraper/scrape"
	"github.com/mfojtik/euclid/scraper/types"
	"os"
)

var opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	DisplayFile string `short:"f" long:"file"`

	AcondURL      string `long:"acond-url" default:"https://localhost:4443"`
	AcondUser     string `long:"acond-user" default:"acond"`
	AcondPassword string `long:"acond-password" default:"acond"`

	MaxKeepValues int `long:"max-values" default:"5"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	if len(opts.DisplayFile) == 0 {
		fmt.Println("ERROR: No display file specified")
		os.Exit(1)
	}

	scraper := scrape.New(opts.AcondURL, opts.AcondUser, opts.AcondPassword)
	val, err := scraper.Scrape()
	if err != nil {
		panic(err)
	}

	display, err := types.ReadDisplayFromFile(opts.DisplayFile)
	if os.IsNotExist(err) {
		fmt.Println("ERR: creating new file")
		display = &types.Display{
			Temperatures: []types.Temperature{
				{Name: "Outside", Values: []types.Value{}, MaxValue: 30.0, MinValue: -5.0},
				{Name: "Living Room", Values: []types.Value{}, MaxValue: 24.0, MinValue: 17.0},
				{Name: "Upstairs", Values: []types.Value{}, MaxValue: 24.0, MinValue: 17.0},
				{Name: "Cellar", Values: []types.Value{}, MaxValue: 17.0, MinValue: 5.0},
			},
		}
	}

	for i := range display.Temperatures {
		switch display.Temperatures[i].Name {
		case "Outside":
			display.Temperatures[i].Values = types.RecordValue(val.Outdoor, display.Temperatures[i].Values, opts.MaxKeepValues)
		case "Living Room":
			display.Temperatures[i].Values = types.RecordValue(val.DownstairsCurrent, display.Temperatures[i].Values, opts.MaxKeepValues)
		case "Upstairs":
			display.Temperatures[i].Values = types.RecordValue(val.UpstairsCurrent, display.Temperatures[i].Values, opts.MaxKeepValues)
		}
		display.Temperatures[i].Trend = types.SetTrend(display.Temperatures[i].Values)
	}

	if err := types.WriteDisplayToFile(opts.DisplayFile, *display); err != nil {
		panic(err)
	}
	dump.P(display)
}

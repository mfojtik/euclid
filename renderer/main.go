package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mfojtik/euclid/renderer/template"
	"github.com/mfojtik/euclid/scraper/types"
	"os"
)

var opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	DisplayFile string `short:"f" long:"file"`
	Output      string `short:"o" long:"output"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	if len(opts.DisplayFile) == 0 || len(opts.Output) == 0 {
		fmt.Println("ERROR: No display or output file specified")
		os.Exit(1)
	}

	display, err := types.ReadDisplayFromFile(opts.DisplayFile)
	if err != nil {
		panic(err)
	}

	out, err := template.RenderTemplate(*display)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(opts.Output, out, os.ModePerm); err != nil {
		panic(err)
	}
}

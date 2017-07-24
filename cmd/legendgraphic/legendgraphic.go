package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"io"

	"github.com/omniscale/legendgraphic"
)

func main() {
	log.SetFlags(0)
	flags := flag.NewFlagSet("legendgraphic", flag.ExitOnError)
	config := flags.String("config", "legend.json", "legend configuration")
	outFile := flags.String("out", "-", "output file (- for stdout)")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(*config)
	if err != nil {
		log.Fatal("opening config: ", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	legend := &legendgraphic.Legend{}
	if err := dec.Decode(legend); err != nil {
		log.Fatal(err)
	}

	err, missing := legendgraphic.FillVars(legend, flags.Args())
	if err != nil {
		log.Fatal(err)
	}
	if missing != nil {
		log.Printf("%s references variables not found in %s:", *config, flags.Args())
		for _, v := range missing {
			log.Println(" ", v)
		}
		os.Exit(1)
	}

	var out io.Writer
	if *outFile != "-" {
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatal("writing legend: ", err)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}
	if err := legendgraphic.RenderLegend(out, legend); err != nil {
		log.Fatal(err)
	}
}

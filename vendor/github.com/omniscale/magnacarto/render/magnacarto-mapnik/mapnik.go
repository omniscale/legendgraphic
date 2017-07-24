package main

import (
	"fmt"
	"log"

	"github.com/natefinch/pie"

	"path/filepath"

	"github.com/omniscale/go-mapnik"
	"github.com/omniscale/magnacarto/render"
)

func renderReq(mapfile string, mapReq render.Request) ([]byte, error) {
	style := filepath.Base(mapfile)
	style = style[:len(style)-len(filepath.Ext(style))] // wihout suffix

	m := mapnik.NewSized(mapReq.Width, mapReq.Height)
	defer m.Free()
	err := m.Load(mapfile)
	if err != nil {
		return nil, err
	}

	m.SetSRS(fmt.Sprintf("+init=epsg:%d", mapReq.EPSGCode))
	m.ZoomTo(mapReq.BBOX[0], mapReq.BBOX[1], mapReq.BBOX[2], mapReq.BBOX[3])

	renderOpts := mapnik.RenderOpts{}
	renderOpts.Format = mapReq.Format
	renderOpts.ScaleFactor = mapReq.ScaleFactor

	b, err := m.Render(renderOpts)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type api struct {
}

type Args struct {
	Mapfile string
	Req     render.Request
}

func (api) Render(args *Args, response *[]byte) error {
	tmp, err := renderReq(args.Mapfile, args.Req)
	*response = tmp
	return err
}

func (api) Is3(args struct{}, response *bool) error {
	if mapnik.Version.Major == 3 {
		*response = true
	} else {
		*response = false
	}
	return nil
}

func (api) RegisterFonts(fontDir string, _ *interface{}) error {
	mapnik.RegisterFonts(fontDir)
	return nil
}

func main() {
	p := pie.NewProvider()
	if err := p.RegisterName("Mapnik", api{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	p.Serve()
}

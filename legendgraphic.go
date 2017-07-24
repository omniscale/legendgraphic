package legendgraphic

import (
	"html/template"
	"io"
	"strings"

	"bytes"
	"strconv"

	"github.com/omniscale/magnacarto/mss"
	"github.com/pkg/errors"
)

type Legend struct {
	Title  string
	Groups []Group
}

type Group struct {
	Title  string
	Layers []Layer
}

type Layer struct {
	LineStyle
	PolygonStyle
	Title string
}

type LineStyle struct {
	LineWidth     float64 `json:"line-width"`
	LineColor     string  `json:"line-color"`
	LineDasharray string  `json:"line-dasharray"`
	OutlineWidth  float64 `json:"outline-width"`
	OutlineColor  string  `json:"outline-color"`
}

type PolygonStyle struct {
	FillColor string `json:"fill-color"`
}

func RenderLegend(w io.Writer, l *Legend) error {
	tmpl := template.New("legend.html")
	tmpl.Funcs(template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	tmpl = template.Must(tmpl.ParseGlob("template/*.html"))
	return tmpl.Execute(w, l)
}

func FillVars(l *Legend, mssFiles []string) (err error, missing []string) {
	style := mss.New()
	for _, mssFile := range mssFiles {
		if err := style.ParseFile(mssFile); err != nil {
			return errors.Wrapf(err, "parsing mss file %s", mssFile), nil
		}
	}
	if err := style.Evaluate(); err != nil {
		return errors.Wrap(err, "evaluating mss files"), nil
	}
	vars := style.Vars()

	for _, g := range l.Groups {
		for i, l := range g.Layers {
			fillVar := func(v string) string {
				if strings.HasPrefix(v, "@") {
					c, ok := vars.GetColor(v[1:len(v)])
					if ok {
						return c.HexString()
					}
					fl, ok := vars.GetFloatList(v[1:len(v)])
					if ok {
						return formatDashArray(fl)
					}
					missing = append(missing, v)
				}
				return v
			}
			g.Layers[i].FillColor = fillVar(l.FillColor)
			g.Layers[i].LineColor = fillVar(l.LineColor)
			g.Layers[i].OutlineColor = fillVar(l.OutlineColor)
			g.Layers[i].LineDasharray = fillVar(l.LineDasharray)
		}
	}

	return nil, missing
}

func formatDashArray(values []float64) string {
	buf := bytes.Buffer{}
	for i, f := range values {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(strconv.FormatFloat(f, 'f', 2, 64))
	}
	return buf.String()
}

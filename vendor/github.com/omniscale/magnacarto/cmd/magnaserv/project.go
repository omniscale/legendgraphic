package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/omniscale/magnacarto/mml"
)

type byLastChange []project

func (p byLastChange) Len() int           { return len(p) }
func (p byLastChange) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byLastChange) Less(i, j int) bool { return p[i].LastChange.Before(p[j].LastChange) }

type project struct {
	Name         string    `json:"name"`
	Base         string    `json:"base"`
	MML          string    `json:"mml"`
	MCP          string    `json:"mcp"`
	LastChange   time.Time `json:"last_change"`
	AvailableMSS []string  `json:"available_mss"`
	SelectedMSS  []string  `json:"selected_mss"`
}

// findProjects searches for .mml files in path and in all sub-directories of
// path, but not any deeper.
func findProjects(path string) ([]project, error) {
	projects := []project{}
	var mmls []string
	if files, err := filepath.Glob(filepath.Join(path, "*.mml")); err != nil {
		return nil, err
	} else {
		mmls = append(mmls, files...)
	}

	if files, err := filepath.Glob(filepath.Join(path, "*", "*.mml")); err != nil {
		return nil, err
	} else {
		mmls = append(mmls, files...)
	}

	for _, mmlFile := range mmls {
		projDir := filepath.Dir(mmlFile)
		projBase, _ := filepath.Rel(path, projDir)

		mssFiles, err := findMSS(projDir)
		if err != nil {
			return nil, err
		}

		r, err := os.Open(mmlFile)
		if err != nil {
			return nil, err
		}
		parsedMML, err := mml.Parse(r)
		r.Close()
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", mmlFile, err)
		}

		lastChange := lastModTime(append([]string{mmlFile}, mssFiles...)...)

		// remove base dir from mml/mss
		mmlFile = filepath.Base(mmlFile)
		for i := range mssFiles {
			mssFiles[i], _ = filepath.Rel(projDir, mssFiles[i])
		}

		name := parsedMML.Name

		projects = append(projects,
			project{
				Name:         name,
				Base:         projBase,
				LastChange:   lastChange,
				MML:          mmlFile,
				MCP:          strings.TrimSuffix(mmlFile, filepath.Ext(mmlFile)) + ".mcp",
				AvailableMSS: mssFiles,
				SelectedMSS:  parsedMML.Stylesheets,
			})
	}

	return projects, nil
}

func findMSS(base string) ([]string, error) {
	var mss []string

	err := filepath.Walk(base, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".mss") {
			mss = append(mss, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mss, nil
}

func lastModTime(files ...string) time.Time {
	mod := time.Time{}
	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			continue
		}
		if fi.ModTime().After(mod) {
			mod = fi.ModTime()
		}
	}
	return mod
}

package main

import "testing"

func TestFindProjects(t *testing.T) {
	projects, err := findProjects("../../regression/cases")
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) == 0 {
		t.Error(projects)
	}

	project := projects[0]
	if project.Name != "Magnacarto Test" ||
		project.MML != "test.mml" ||
		project.MCP != "test.mcp" ||
		project.Base != "010-linestrings-default" ||
		len(project.SelectedMSS) != 1 ||
		len(project.AvailableMSS) != 1 ||
		project.SelectedMSS[0] != "test.mss" ||
		project.AvailableMSS[0] != "test.mss" {
		t.Error("unexpected project", project)
	}
}

package pluginadapter

import (
	"github.com/EME130/lazymd/internal/nav"
	"github.com/EME130/lazymd/internal/pluginapi"
)

// NavAdapter wraps a nav.Navigator as a pluginapi.NavAPI.
type NavAdapter struct {
	Nav nav.Navigator
}

func (a *NavAdapter) ReadSection(headingPath string) (*pluginapi.SectionResult, error) {
	sec, err := a.Nav.ReadSection(headingPath)
	if err != nil {
		return nil, err
	}
	return &pluginapi.SectionResult{
		Title:     sec.Title,
		Level:     sec.Level,
		Content:   sec.Content,
		StartLine: sec.StartLine,
		EndLine:   sec.EndLine,
	}, nil
}

func (a *NavAdapter) ListTasks(section *string, done *bool) ([]pluginapi.TaskInfo, error) {
	status := nav.TaskAll
	if done != nil {
		if *done {
			status = nav.TaskDone
		} else {
			status = nav.TaskPending
		}
	}
	tasks, err := a.Nav.ListTasks(section, status)
	if err != nil {
		return nil, err
	}
	out := make([]pluginapi.TaskInfo, len(tasks))
	for i, t := range tasks {
		out[i] = pluginapi.TaskInfo{
			Line:       t.Line,
			Done:       t.Done,
			Text:       t.Text,
			Breadcrumb: t.Breadcrumb,
		}
	}
	return out, nil
}

func (a *NavAdapter) GetBreadcrumb(line int) (string, error) {
	return a.Nav.GetBreadcrumb(line)
}

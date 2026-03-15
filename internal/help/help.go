package help

import (
	"embed"
	"strings"
)

//go:embed content/*.md
var content embed.FS

var topicFiles = map[string]string{
	"overview": "content/overview.md",
	"keys":     "content/keys.md",
	"brain":    "content/brain.md",
	"commands": "content/commands.md",
	"mcp":      "content/mcp.md",
	"panels":   "content/panels.md",
}

// GetTopic returns the markdown content for a help topic.
func GetTopic(name string) (string, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "help" {
		name = "overview"
	}
	file, ok := topicFiles[name]
	if !ok {
		return "", false
	}
	data, err := content.ReadFile(file)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// Topics returns all available static topic names.
func Topics() []string {
	names := make([]string, 0, len(topicFiles))
	for name := range topicFiles {
		names = append(names, name)
	}
	return names
}

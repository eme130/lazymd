package help

import "testing"

func TestGetTopicOverview(t *testing.T) {
	content, ok := GetTopic("overview")
	if !ok {
		t.Fatal("expected overview topic to exist")
	}
	if len(content) == 0 {
		t.Fatal("expected overview content to be non-empty")
	}
}

func TestGetTopicAllStatic(t *testing.T) {
	topics := []string{"overview", "keys", "brain", "commands", "mcp", "panels"}
	for _, name := range topics {
		content, ok := GetTopic(name)
		if !ok {
			t.Errorf("topic %q not found", name)
		}
		if len(content) == 0 {
			t.Errorf("topic %q is empty", name)
		}
	}
}

func TestGetTopicUnknown(t *testing.T) {
	_, ok := GetTopic("nonexistent")
	if ok {
		t.Error("expected unknown topic to return false")
	}
}

func TestTopics(t *testing.T) {
	topics := Topics()
	if len(topics) < 6 {
		t.Errorf("expected at least 6 topics, got %d", len(topics))
	}
}

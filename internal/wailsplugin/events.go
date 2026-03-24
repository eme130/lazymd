package wailsplugin

import "github.com/EME130/lazymd/internal/pluginapi"

var eventNameMap = map[pluginapi.EventType]string{
	pluginapi.EventBufferChanged: "buffer:changed",
	pluginapi.EventFileOpened:    "file:opened",
	pluginapi.EventFileSaved:     "file:saved",
	pluginapi.EventFileClosed:    "file:closed",
	pluginapi.EventCursorMoved:   "cursor:moved",
	pluginapi.EventModeChanged:   "mode:changed",
	pluginapi.EventGraphUpdated:  "graph:updated",
}

func WailsEventName(t pluginapi.EventType) string {
	if name, ok := eventNameMap[t]; ok {
		return name
	}
	return string(t)
}

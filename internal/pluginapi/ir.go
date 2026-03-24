package pluginapi

// OpType identifies the type of operation flowing through the IR.
type OpType string

const (
	OpInsertText  OpType = "insert_text"
	OpDeleteRange OpType = "delete_range"
	OpOpenFile    OpType = "open_file"
	OpSaveFile    OpType = "save_file"
	OpSetCursor   OpType = "set_cursor"
	OpSetMode     OpType = "set_mode"
	OpRunCommand  OpType = "run_command"
)

// Operation is a mutation flowing from frontend → passes → core → backend.
type Operation struct {
	Type     OpType
	Data     map[string]any
	Canceled bool
	Origin   string // plugin that emitted it
}

// NewOperation creates an operation with the given type and origin.
func NewOperation(typ OpType, origin string) *Operation {
	return &Operation{
		Type:   typ,
		Data:   make(map[string]any),
		Origin: origin,
	}
}

// Cancel marks this operation as canceled. Passes use this to abort operations.
func (op *Operation) Cancel() {
	op.Canceled = true
}

// EventType identifies the type of event broadcast to plugins.
type EventType string

const (
	EventFileOpened    EventType = "file_opened"
	EventFileSaved     EventType = "file_saved"
	EventFileClosed    EventType = "file_closed"
	EventBufferChanged EventType = "buffer_changed"
	EventModeChanged   EventType = "mode_changed"
	EventCursorMoved   EventType = "cursor_moved"
	EventGraphUpdated  EventType = "graph_updated"
)

// Event is a notification broadcast from the core to all plugins.
type Event struct {
	Type EventType
	Data map[string]any
}

// NewEvent creates an event with the given type.
func NewEvent(typ EventType) *Event {
	return &Event{
		Type: typ,
		Data: make(map[string]any),
	}
}

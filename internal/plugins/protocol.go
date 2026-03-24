package plugins

import "encoding/json"

// JSON-RPC 2.0 message types for the external plugin protocol.

// rpcRequest is a JSON-RPC 2.0 request from host to plugin.
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// rpcResponse is a JSON-RPC 2.0 response from plugin to host.
type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

// rpcError is a JSON-RPC 2.0 error object.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// initializeParams is sent with the "initialize" request.
type initializeParams struct {
	HostVersion  string   `json:"host_version"`
	Capabilities []string `json:"capabilities"`
}

// initializeResult is returned by the plugin on initialize.
type initializeResult struct {
	OK bool `json:"ok"`
}

// commandParams is sent with the "command" request.
type commandParams struct {
	Name    string         `json:"name"`
	Args    string         `json:"args"`
	Context commandContext `json:"context"`
}

// commandContext is a snapshot of editor state sent with commands.
type commandContext struct {
	FilePath  string `json:"file_path"`
	CursorRow int    `json:"cursor_row"`
	CursorCol int    `json:"cursor_col"`
	Mode      string `json:"mode"`
	Content   string `json:"content"`
}

// commandResult is returned by the plugin after executing a command.
type commandResult struct {
	Status    string     `json:"status,omitempty"`
	Mutations []mutation `json:"mutations,omitempty"`
}

// mutation is an editor mutation requested by an external plugin.
type mutation struct {
	Type    string `json:"type"` // set_status, insert_text, delete_lines, set_cursor, cancel
	Text    string `json:"text,omitempty"`
	IsError bool   `json:"is_error,omitempty"`
	Row     int    `json:"row,omitempty"`
	Col     int    `json:"col,omitempty"`
	Start   int    `json:"start,omitempty"`
	End     int    `json:"end,omitempty"`
}

// eventParams is sent with the "event" notification.
type eventParams struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

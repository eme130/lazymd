package plugins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/EME130/lazymd/internal/pluginapi"
)

const (
	rpcTimeout = 2 * time.Second
)

// ExternalPlugin wraps a child process that communicates via JSON-RPC 2.0 over stdio.
// It implements pluginapi.BackendPlugin (external frontend plugins are future work).
type ExternalPlugin struct {
	manifest *Manifest
	cmd      *exec.Cmd
	encoder  *json.Encoder
	scanner  *bufio.Scanner
	mu       sync.Mutex
	nextID   int
	running  bool

	// subscribed events from manifest
	events map[string]bool
}

// NewExternalPlugin creates an ExternalPlugin from a manifest.
func NewExternalPlugin(m *Manifest, entryPath string) *ExternalPlugin {
	events := make(map[string]bool, len(m.Events))
	for _, e := range m.Events {
		events[e] = true
	}
	return &ExternalPlugin{
		manifest: m,
		events:   events,
		cmd:      exec.Command(entryPath),
	}
}

func (p *ExternalPlugin) Info() pluginapi.PluginInfo {
	return pluginapi.PluginInfo{
		Name:        p.manifest.Name,
		Version:     p.manifest.Version,
		Author:      p.manifest.Author,
		Description: p.manifest.Description,
	}
}

func (p *ExternalPlugin) Init(_ *pluginapi.BackendContext) error {
	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("plugin %s: stdin pipe: %w", p.manifest.Name, err)
	}
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("plugin %s: stdout pipe: %w", p.manifest.Name, err)
	}

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("plugin %s: start: %w", p.manifest.Name, err)
	}

	p.encoder = json.NewEncoder(stdin)
	p.scanner = bufio.NewScanner(stdout)
	p.scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB max line
	p.running = true

	// Send initialize
	params := initializeParams{
		HostVersion:  "0.1.0",
		Capabilities: p.manifest.Capabilities,
	}
	_, err = p.call("initialize", params)
	if err != nil {
		p.Shutdown()
		return fmt.Errorf("plugin %s: initialize: %w", p.manifest.Name, err)
	}

	return nil
}

func (p *ExternalPlugin) Shutdown() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return nil
	}
	p.running = false

	// Try graceful shutdown
	_ = p.sendRequest("shutdown", nil)

	done := make(chan error, 1)
	go func() { done <- p.cmd.Wait() }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		p.cmd.Process.Kill()
		<-done
	}
	return nil
}

func (p *ExternalPlugin) OnOperation(_ *pluginapi.BackendContext, _ *pluginapi.Operation) {
	// External plugins receive operations as events for simplicity
}

func (p *ExternalPlugin) OnEvent(_ *pluginapi.BackendContext, event *pluginapi.Event) {
	if !p.events[string(event.Type)] {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return
	}

	params := eventParams{
		Type: string(event.Type),
		Data: event.Data,
	}
	_ = p.sendNotification("event", params)
}

func (p *ExternalPlugin) Commands() []pluginapi.CommandDef {
	cmds := make([]pluginapi.CommandDef, len(p.manifest.Commands))
	for i, mc := range p.manifest.Commands {
		cmdName := mc.Name
		cmds[i] = pluginapi.CommandDef{
			Name:        mc.Name,
			Description: mc.Description,
			Handler: func(ctx *pluginapi.BackendContext, args string) error {
				return p.executeCommand(ctx, cmdName, args)
			},
		}
	}
	return cmds
}

func (p *ExternalPlugin) Capabilities() []pluginapi.Capability {
	caps := make([]pluginapi.Capability, len(p.manifest.Capabilities))
	for i, c := range p.manifest.Capabilities {
		caps[i] = pluginapi.Capability{Name: c}
	}
	return caps
}

func (p *ExternalPlugin) executeCommand(ctx *pluginapi.BackendContext, name, args string) error {
	params := commandParams{
		Name: name,
		Args: args,
	}
	if ctx != nil && ctx.Editor != nil {
		params.Context = commandContext{
			FilePath:  ctx.Editor.FilePath(),
			CursorRow: ctx.Editor.CursorRow(),
			CursorCol: ctx.Editor.CursorCol(),
			Mode:      ctx.Editor.Mode(),
			Content:   ctx.Editor.Content(),
		}
	}

	raw, err := p.call("command", params)
	if err != nil {
		return err
	}

	var result commandResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return fmt.Errorf("parse command result: %w", err)
	}

	// Apply mutations
	p.applyMutations(ctx, result.Mutations)

	if result.Status != "" && ctx != nil && ctx.Editor != nil {
		ctx.Editor.SetStatus(result.Status, false)
	}

	return nil
}

func (p *ExternalPlugin) applyMutations(ctx *pluginapi.BackendContext, mutations []mutation) {
	if ctx == nil || ctx.Editor == nil {
		return
	}
	for _, m := range mutations {
		switch m.Type {
		case "set_status":
			ctx.Editor.SetStatus(m.Text, m.IsError)
		case "insert_text":
			ctx.Editor.InsertAt(m.Row, m.Col, m.Text)
		case "delete_lines":
			ctx.Editor.DeleteLines(m.Start, m.End)
		case "set_cursor":
			ctx.Editor.SetCursorRow(m.Row)
			ctx.Editor.SetCursorCol(m.Col)
		}
	}
}

// call sends a JSON-RPC request and waits for a response.
func (p *ExternalPlugin) call(method string, params any) (json.RawMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return nil, fmt.Errorf("plugin not running")
	}

	p.nextID++
	id := p.nextID

	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  raw,
	}

	if err := p.encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	// Read response (with timeout)
	done := make(chan struct{})
	var resp rpcResponse
	var scanErr error

	go func() {
		defer close(done)
		if p.scanner.Scan() {
			scanErr = json.Unmarshal(p.scanner.Bytes(), &resp)
		} else {
			scanErr = fmt.Errorf("plugin closed stdout")
			if err := p.scanner.Err(); err != nil {
				scanErr = err
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(rpcTimeout):
		return nil, fmt.Errorf("plugin %s: timeout waiting for response", p.manifest.Name)
	}

	if scanErr != nil {
		return nil, scanErr
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("plugin error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}

// sendRequest sends a request without waiting for response.
func (p *ExternalPlugin) sendRequest(method string, params any) error {
	p.nextID++
	raw, _ := json.Marshal(params)
	req := rpcRequest{
		JSONRPC: "2.0",
		ID:      p.nextID,
		Method:  method,
		Params:  raw,
	}
	return p.encoder.Encode(req)
}

// sendNotification sends a JSON-RPC notification (no ID, no response expected).
func (p *ExternalPlugin) sendNotification(method string, params any) error {
	raw, _ := json.Marshal(params)
	req := rpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  raw,
	}
	return p.encoder.Encode(req)
}

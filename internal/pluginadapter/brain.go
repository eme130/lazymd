package pluginadapter

import (
	"github.com/EME130/lazymd/internal/brain"
	"github.com/EME130/lazymd/internal/pluginapi"
)

// BrainAdapter wraps a brain.Graph as a pluginapi.BrainAPI.
type BrainAdapter struct {
	Graph *brain.Graph
}

func (a *BrainAdapter) NodeCount() int {
	if a.Graph == nil {
		return 0
	}
	return a.Graph.NodeCount()
}

func (a *BrainAdapter) Nodes() []pluginapi.NodeInfo {
	if a.Graph == nil {
		return nil
	}
	nodes := a.Graph.Nodes
	out := make([]pluginapi.NodeInfo, len(nodes))
	for i, n := range nodes {
		out[i] = pluginapi.NodeInfo{
			Name:      n.Name,
			Path:      n.Path,
			LinkCount: len(n.OutLinks) + len(n.InLinks),
		}
	}
	return out
}

func (a *BrainAdapter) GetBacklinks(name string) []string {
	if a.Graph == nil {
		return nil
	}
	id, ok := a.Graph.Resolve(name)
	if !ok {
		return nil
	}
	backlinks := a.Graph.GetBacklinks(id)
	out := make([]string, len(backlinks))
	for i, bl := range backlinks {
		out[i] = a.Graph.Nodes[bl].Name
	}
	return out
}

func (a *BrainAdapter) GetNeighbors(name string) []string {
	if a.Graph == nil {
		return nil
	}
	id, ok := a.Graph.Resolve(name)
	if !ok {
		return nil
	}
	neighbors := a.Graph.GetNeighbors(id, 1)
	out := make([]string, 0, len(neighbors))
	for _, nid := range neighbors {
		if nid != id {
			out = append(out, a.Graph.Nodes[nid].Name)
		}
	}
	return out
}

func (a *BrainAdapter) FindPath(from, to string) ([]string, bool) {
	if a.Graph == nil {
		return nil, false
	}
	fromID, ok1 := a.Graph.Resolve(from)
	toID, ok2 := a.Graph.Resolve(to)
	if !ok1 || !ok2 {
		return nil, false
	}
	path, found := a.Graph.FindPath(fromID, toID)
	if !found {
		return nil, false
	}
	out := make([]string, len(path))
	for i, nid := range path {
		out[i] = a.Graph.Nodes[nid].Name
	}
	return out, true
}

func (a *BrainAdapter) GetOrphans() []string {
	if a.Graph == nil {
		return nil
	}
	orphans := a.Graph.GetOrphans()
	out := make([]string, len(orphans))
	for i, id := range orphans {
		out[i] = a.Graph.Nodes[id].Name
	}
	return out
}

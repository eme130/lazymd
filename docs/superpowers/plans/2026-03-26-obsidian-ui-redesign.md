# Obsidian-like UI Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the Wails desktop frontend from a 3-column split layout to an Obsidian-inspired single-panel layout with reading/editing mode toggle, tab bar, icon ribbon, and navigation history.

**Architecture:** Replace the current 3-column grid (file tree | editor | preview) with a 2-column layout: icon ribbon + optional file tree sidebar on the left, and a content area on the right that switches between reading mode (rendered markdown), editing mode (CodeMirror), and brain graph view. Tab bar and navigation bar sit above the content. State management uses Svelte stores for tabs, navigation history, and view mode.

**Tech Stack:** Svelte 4, TypeScript, CodeMirror 6, remark/rehype, D3.js, Wails v2, Go

**Spec:** `docs/superpowers/specs/2026-03-26-obsidian-ui-redesign.md`

---

## File Structure

### New Files (Frontend)
| File | Responsibility |
|------|---------------|
| `src/lib/stores.ts` | Svelte stores: tabs, viewMode, navHistory, sidebarOpen |
| `src/components/IconRibbon.svelte` | Vertical icon strip (file explorer, graph, settings) |
| `src/components/TabBar.svelte` | Horizontal open-file tabs with close buttons |
| `src/components/NavBar.svelte` | Back/forward arrows, file title, mode toggle |
| `src/components/ContentArea.svelte` | Orchestrates reading/editing/graph views |

### Modified Files (Frontend)
| File | Changes |
|------|---------|
| `src/App.svelte` | Rewrite layout to ribbon + sidebar + content grid |
| `src/components/FileTree.svelte` | Add selected-file highlight, use stores |
| `src/components/StatusBar.svelte` | Add word count, char count, backlink count |
| `src/components/Preview.svelte` | Accept content as prop (not self-fetching) |
| `src/components/Editor.svelte` | Accept content as prop, emit save events |

### Modified Files (Go Backend)
| File | Changes |
|------|---------|
| `internal/wailsplugin/app.go` | Add `GetWordCount`, `GetCharCount`, `GetBacklinkCount` methods |
| `internal/wailsplugin/app_test.go` | Tests for new methods |

---

## Task 1: Add Go backend methods (word count, char count, backlink count)

**Files:**
- Modify: `internal/wailsplugin/app.go:264-274`
- Modify: `internal/wailsplugin/app_test.go`

- [ ] **Step 1: Write failing tests for GetWordCount, GetCharCount, GetBacklinkCount**

Add to `internal/wailsplugin/app_test.go`:

```go
func TestAppGetWordCount(t *testing.T) {
	editor := newMockEditor("hello world foo")
	app := newTestApp(editor, nil, nil)
	if got := app.GetWordCount(); got != 3 {
		t.Errorf("expected 3 words, got %d", got)
	}
}

func TestAppGetWordCountEmpty(t *testing.T) {
	editor := newMockEditor("")
	app := newTestApp(editor, nil, nil)
	if got := app.GetWordCount(); got != 0 {
		t.Errorf("expected 0 words, got %d", got)
	}
}

func TestAppGetCharCount(t *testing.T) {
	editor := newMockEditor("hello world")
	app := newTestApp(editor, nil, nil)
	if got := app.GetCharCount(); got != 11 {
		t.Errorf("expected 11 chars, got %d", got)
	}
}

func TestAppGetBacklinkCount(t *testing.T) {
	app := newTestApp(newMockEditor(""), nil, &mockBrainAPI{
		backlinks: map[string][]string{"test": {"a.md", "b.md"}},
	})
	app.editor.(*mockEditorAPI).filePath = "test.md"
	if got := app.GetBacklinkCount(); got != 2 {
		t.Errorf("expected 2 backlinks, got %d", got)
	}
}

func TestAppGetBacklinkCountNoBrain(t *testing.T) {
	app := newTestApp(newMockEditor(""), nil, nil)
	if got := app.GetBacklinkCount(); got != 0 {
		t.Errorf("expected 0 backlinks, got %d", got)
	}
}
```

Also add mock brain to test file:

```go
type mockBrainAPI struct {
	backlinks map[string][]string
}

func (m *mockBrainAPI) NodeCount() int                          { return 0 }
func (m *mockBrainAPI) Nodes() []pluginapi.NodeInfo             { return nil }
func (m *mockBrainAPI) GetBacklinks(name string) []string       { return m.backlinks[name] }
func (m *mockBrainAPI) GetNeighbors(name string) []string       { return nil }
func (m *mockBrainAPI) FindPath(from, to string) ([]string, bool) { return nil, false }
func (m *mockBrainAPI) GetOrphans() []string                    { return nil }
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/wailsplugin/ -run "TestAppGet(Word|Char|Backlink)" -v`
Expected: FAIL — methods not defined

- [ ] **Step 3: Implement the three methods**

Add to `internal/wailsplugin/app.go` after the existing `GetStatus` method:

```go
func (a *App) GetWordCount() int {
	if a.editor == nil {
		return 0
	}
	return len(strings.Fields(a.editor.Content()))
}

func (a *App) GetCharCount() int {
	if a.editor == nil {
		return 0
	}
	return len([]rune(a.editor.Content()))
}

func (a *App) GetBacklinkCount() int {
	if a.brain == nil || a.editor == nil {
		return 0
	}
	filePath := a.editor.FilePath()
	name := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	if name == "" {
		return 0
	}
	return len(a.brain.GetBacklinks(name))
}
```

Note: `strings` and `path/filepath` are already imported in app.go.

- [ ] **Step 4: Run tests to verify they pass**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/wailsplugin/ -run "TestAppGet(Word|Char|Backlink)" -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/wailsplugin/ -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/wailsplugin/app.go internal/wailsplugin/app_test.go
git commit -m "feat(desktop): add word count, char count, backlink count methods"
```

---

## Task 2: Create Svelte stores for shared UI state

**Files:**
- Create: `cmd/lm-desktop/frontend/src/lib/stores.ts`

- [ ] **Step 1: Create the stores file**

```typescript
import { writable, derived } from 'svelte/store';

// --- Tab State ---
export interface Tab {
  path: string;
  name: string;
  modified: boolean;
}

export const tabs = writable<Tab[]>([]);
export const activeTabIndex = writable<number>(-1);

export const activeTab = derived(
  [tabs, activeTabIndex],
  ([$tabs, $idx]) => $idx >= 0 && $idx < $tabs.length ? $tabs[$idx] : null
);

export function openTab(path: string, name: string) {
  tabs.update(t => {
    const existing = t.findIndex(tab => tab.path === path);
    if (existing >= 0) {
      activeTabIndex.set(existing);
      return t;
    }
    const newTabs = [...t, { path, name, modified: false }];
    activeTabIndex.set(newTabs.length - 1);
    return newTabs;
  });
}

export function closeTab(index: number) {
  tabs.update(t => {
    const newTabs = t.filter((_, i) => i !== index);
    activeTabIndex.update(current => {
      if (newTabs.length === 0) return -1;
      if (index <= current) return Math.max(0, current - 1);
      return current;
    });
    return newTabs;
  });
}

// --- View Mode ---
export type ViewMode = 'reading' | 'editing' | 'graph';
export const viewMode = writable<ViewMode>('reading');

export function toggleEditMode() {
  viewMode.update(m => m === 'editing' ? 'reading' : 'editing');
}

// --- Navigation History ---
const navHistory = writable<string[]>([]);
const navIndex = writable<number>(-1);

export const canGoBack = derived(navIndex, $i => $i > 0);
export const canGoForward = derived(
  [navIndex, navHistory],
  ([$i, $h]) => $i < $h.length - 1
);

export function pushHistory(path: string) {
  navHistory.update(h => {
    let idx: number;
    navIndex.subscribe(v => idx = v)();
    // Trim forward history
    const trimmed = h.slice(0, idx + 1);
    trimmed.push(path);
    navIndex.set(trimmed.length - 1);
    return trimmed;
  });
}

export function goBack(): string | null {
  let result: string | null = null;
  navIndex.update(i => {
    if (i > 0) {
      const newIdx = i - 1;
      navHistory.subscribe(h => result = h[newIdx])();
      return newIdx;
    }
    return i;
  });
  return result;
}

export function goForward(): string | null {
  let result: string | null = null;
  navIndex.update(i => {
    let len = 0;
    navHistory.subscribe(h => len = h.length)();
    if (i < len - 1) {
      const newIdx = i + 1;
      navHistory.subscribe(h => result = h[newIdx])();
      return newIdx;
    }
    return i;
  });
  return result;
}

// --- Sidebar ---
export const sidebarOpen = writable<boolean>(true);
```

- [ ] **Step 2: Verify it compiles**

Run: `cd cmd/lm-desktop/frontend && npx svelte-check --threshold warning 2>&1 | head -20`
Expected: No errors in stores.ts

- [ ] **Step 3: Commit**

```bash
git add cmd/lm-desktop/frontend/src/lib/stores.ts
git commit -m "feat(desktop): add Svelte stores for tabs, nav history, view mode"
```

---

## Task 3: Create IconRibbon component

**Files:**
- Create: `cmd/lm-desktop/frontend/src/components/IconRibbon.svelte`

- [ ] **Step 1: Create the component**

```svelte
<script lang="ts">
  import { sidebarOpen, viewMode } from '../lib/stores';

  function toggleSidebar() {
    sidebarOpen.update(v => !v);
  }

  function showGraph() {
    viewMode.set('graph');
  }

  $: isGraphActive = $viewMode === 'graph';
</script>

<nav class="ribbon">
  <button
    class="ribbon-btn"
    class:active={$sidebarOpen}
    on:click={toggleSidebar}
    title="File explorer (Cmd+B)"
  >
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
      <polyline points="14 2 14 8 20 8"/>
      <line x1="16" y1="13" x2="8" y2="13"/>
      <line x1="16" y1="17" x2="8" y2="17"/>
      <polyline points="10 9 9 9 8 9"/>
    </svg>
  </button>

  <button
    class="ribbon-btn"
    class:active={isGraphActive}
    on:click={showGraph}
    title="Graph view"
  >
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="5" r="3"/>
      <circle cx="5" cy="19" r="3"/>
      <circle cx="19" cy="19" r="3"/>
      <line x1="12" y1="8" x2="5" y2="16"/>
      <line x1="12" y1="8" x2="19" y2="16"/>
    </svg>
  </button>

  <div class="spacer"></div>

  <button class="ribbon-btn" title="Settings" disabled>
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="3"/>
      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.32 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
    </svg>
  </button>
</nav>

<style>
  .ribbon {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 0;
    gap: 2px;
    background: var(--lm-ribbon-bg, #13141c);
    border-right: 1px solid var(--lm-border, #3b4261);
    width: 40px;
    min-width: 40px;
  }

  .ribbon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: var(--lm-fg, #c0caf5);
    opacity: 0.5;
    cursor: pointer;
    border-radius: 4px;
    position: relative;
  }

  .ribbon-btn:hover {
    opacity: 0.8;
    background: var(--lm-highlight, #292e42);
  }

  .ribbon-btn.active {
    opacity: 1;
  }

  .ribbon-btn.active::before {
    content: '';
    position: absolute;
    left: -8px;
    top: 6px;
    bottom: 6px;
    width: 3px;
    background: var(--lm-link, #7aa2f7);
    border-radius: 0 2px 2px 0;
  }

  .ribbon-btn:disabled {
    opacity: 0.25;
    cursor: default;
  }

  .spacer {
    flex: 1;
  }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/IconRibbon.svelte
git commit -m "feat(desktop): add IconRibbon component"
```

---

## Task 4: Create TabBar component

**Files:**
- Create: `cmd/lm-desktop/frontend/src/components/TabBar.svelte`

- [ ] **Step 1: Create the component**

```svelte
<script lang="ts">
  import { tabs, activeTabIndex, closeTab } from '../lib/stores';

  function selectTab(index: number) {
    activeTabIndex.set(index);
  }

  function onClose(e: MouseEvent, index: number) {
    e.stopPropagation();
    closeTab(index);
  }
</script>

<div class="tab-bar">
  {#each $tabs as tab, i}
    <button
      class="tab"
      class:active={i === $activeTabIndex}
      on:click={() => selectTab(i)}
    >
      <span class="tab-name">{tab.name}{tab.modified ? ' •' : ''}</span>
      <span class="tab-close" on:click={(e) => onClose(e, i)} on:keydown={() => {}}>×</span>
    </button>
  {/each}
</div>

<style>
  .tab-bar {
    display: flex;
    background: var(--lm-ribbon-bg, #13141c);
    border-bottom: 1px solid var(--lm-border, #3b4261);
    overflow-x: auto;
    min-height: 35px;
    -webkit-app-region: drag;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: transparent;
    border: none;
    border-right: 1px solid var(--lm-border, #3b4261);
    color: var(--lm-fg, #c0caf5);
    opacity: 0.5;
    cursor: pointer;
    font-size: 12px;
    white-space: nowrap;
    -webkit-app-region: no-drag;
  }

  .tab:hover {
    opacity: 0.8;
  }

  .tab.active {
    opacity: 1;
    background: var(--lm-bg, #1a1b26);
  }

  .tab-close {
    font-size: 14px;
    line-height: 1;
    opacity: 0;
    border-radius: 3px;
    padding: 0 2px;
  }

  .tab:hover .tab-close {
    opacity: 0.6;
  }

  .tab-close:hover {
    opacity: 1 !important;
    background: var(--lm-highlight, #292e42);
  }

  .tab-name {
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/TabBar.svelte
git commit -m "feat(desktop): add TabBar component"
```

---

## Task 5: Create NavBar component

**Files:**
- Create: `cmd/lm-desktop/frontend/src/components/NavBar.svelte`

- [ ] **Step 1: Create the component**

```svelte
<script lang="ts">
  import { activeTab, viewMode, canGoBack, canGoForward, goBack, goForward, toggleEditMode } from '../lib/stores';
  import { OpenFile } from '../../wailsjs/go/wailsplugin/App';

  function handleBack() {
    const path = goBack();
    if (path) OpenFile(path);
  }

  function handleForward() {
    const path = goForward();
    if (path) OpenFile(path);
  }

  $: isEditing = $viewMode === 'editing';
  $: isGraph = $viewMode === 'graph';
</script>

<div class="navbar">
  <div class="nav-left">
    <button class="nav-btn" disabled={!$canGoBack} on:click={handleBack} title="Back (Cmd+[)">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6"/>
      </svg>
    </button>
    <button class="nav-btn" disabled={!$canGoForward} on:click={handleForward} title="Forward (Cmd+])">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6"/>
      </svg>
    </button>
  </div>

  <div class="nav-title">
    {#if isGraph}
      Graph View
    {:else if $activeTab}
      {$activeTab.name}
    {:else}
      No file open
    {/if}
  </div>

  <div class="nav-right">
    {#if !isGraph}
      <button
        class="nav-btn"
        on:click={toggleEditMode}
        title={isEditing ? 'Reading mode (Cmd+E)' : 'Editing mode (Cmd+E)'}
      >
        {#if isEditing}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
        {:else}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
        {/if}
      </button>
    {/if}
  </div>
</div>

<style>
  .navbar {
    display: flex;
    align-items: center;
    padding: 4px 12px;
    border-bottom: 1px solid var(--lm-border, #3b4261);
    min-height: 32px;
    gap: 8px;
  }

  .nav-left, .nav-right {
    display: flex;
    gap: 2px;
    min-width: 60px;
  }

  .nav-right {
    justify-content: flex-end;
  }

  .nav-title {
    flex: 1;
    text-align: center;
    font-size: 13px;
    opacity: 0.7;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .nav-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border: none;
    background: transparent;
    color: var(--lm-fg, #c0caf5);
    cursor: pointer;
    border-radius: 4px;
    opacity: 0.6;
  }

  .nav-btn:hover:not(:disabled) {
    opacity: 1;
    background: var(--lm-highlight, #292e42);
  }

  .nav-btn:disabled {
    opacity: 0.2;
    cursor: default;
  }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/NavBar.svelte
git commit -m "feat(desktop): add NavBar component with back/forward and mode toggle"
```

---

## Task 6: Create ContentArea component

**Files:**
- Create: `cmd/lm-desktop/frontend/src/components/ContentArea.svelte`

- [ ] **Step 1: Create the component**

This component orchestrates which view is shown based on `viewMode` store.

```svelte
<script lang="ts">
  import { viewMode } from '../lib/stores';
  import Editor from './Editor.svelte';
  import Preview from './Preview.svelte';
  import BrainGraph from './BrainGraph.svelte';
</script>

<div class="content-area">
  {#if $viewMode === 'graph'}
    <BrainGraph />
  {:else}
    <div class="view" class:hidden={$viewMode !== 'reading'}>
      <Preview />
    </div>
    <div class="view" class:hidden={$viewMode !== 'editing'}>
      <Editor />
    </div>
  {/if}
</div>

<style>
  .content-area {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .view { flex: 1; overflow: hidden; }
  .view.hidden { display: none; }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/ContentArea.svelte
git commit -m "feat(desktop): add ContentArea component for view mode switching"
```

---

## Task 7: Update StatusBar with word count, char count, backlinks

**Files:**
- Modify: `cmd/lm-desktop/frontend/src/components/StatusBar.svelte`

- [ ] **Step 1: Read current StatusBar.svelte** (already known from exploration)

- [ ] **Step 2: Rewrite StatusBar**

Replace the entire component with:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { GetWordCount, GetCharCount, GetBacklinkCount } from '../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged, onFileOpened, onFileSaved } from '../lib/events';
  import { viewMode } from '../lib/stores';

  let wordCount = 0;
  let charCount = 0;
  let backlinkCount = 0;

  async function refresh() {
    [wordCount, charCount, backlinkCount] = await Promise.all([
      GetWordCount(),
      GetCharCount(),
      GetBacklinkCount(),
    ]);
  }

  onMount(() => {
    refresh();
    onBufferChanged(() => refresh());
    onFileOpened(() => refresh());
    onFileSaved(() => refresh());
  });
</script>

<div class="status">
  <span class="left">
    <span class="backlinks">{backlinkCount} backlinks</span>
  </span>
  <span class="right">
    {#if $viewMode === 'editing'}
      <span class="indicator">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
      </span>
    {/if}
    <span>{wordCount} words</span>
    <span>{charCount} characters</span>
  </span>
</div>

<style>
  .status {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
    opacity: 0.6;
  }
  .left, .right { display: flex; gap: 16px; align-items: center; }
  .indicator { display: flex; align-items: center; }
</style>
```

- [ ] **Step 3: Regenerate Wails bindings**

Run: `cd cmd/lm-desktop && export PATH="/usr/local/go/bin:$PATH" && wails generate module`

This generates the TypeScript bindings for the new Go methods (`GetWordCount`, `GetCharCount`, `GetBacklinkCount`).

- [ ] **Step 4: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/StatusBar.svelte
git commit -m "feat(desktop): update StatusBar with word count, char count, backlinks"
```

---

## Task 8: Update FileTree to use stores

**Files:**
- Modify: `cmd/lm-desktop/frontend/src/components/FileTree.svelte`

- [ ] **Step 1: Update FileTree to highlight selected file and use tab store**

Replace the `open` function and add selected highlighting:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { ListFiles, OpenFile } from '../../wailsjs/go/wailsplugin/App';
  import { onFileOpened, onFileSaved } from '../lib/events';
  import { openTab, activeTab, pushHistory, viewMode } from '../lib/stores';

  interface FileEntry {
    name: string;
    path: string;
    isDir: boolean;
    children?: FileEntry[];
  }

  let files: FileEntry[] = [];
  let expanded: Set<string> = new Set();

  async function refresh() {
    files = await ListFiles('.');
  }

  function toggle(path: string) {
    if (expanded.has(path)) {
      expanded.delete(path);
    } else {
      expanded.add(path);
    }
    expanded = expanded;
  }

  async function open(path: string, name: string) {
    openTab(path, name);
    pushHistory(path);
    viewMode.set('reading');
    await OpenFile(path);
  }

  $: selectedPath = $activeTab?.path ?? '';

  onMount(() => {
    refresh();
    onFileOpened(() => refresh());
    onFileSaved(() => refresh());
  });
</script>

<div class="filetree-content">
  {#each files as entry}
    {#if entry.isDir}
      <div class="dir" on:click={() => toggle(entry.path)} on:keydown={() => {}}>
        <span class="icon">{expanded.has(entry.path) ? '▼' : '▶'}</span>
        {entry.name}
      </div>
      {#if expanded.has(entry.path) && entry.children}
        <div class="children">
          {#each entry.children as child}
            {#if child.isDir}
              <div class="dir">{child.name}</div>
            {:else}
              <div
                class="file"
                class:selected={child.path === selectedPath}
                on:click={() => open(child.path, child.name)}
                on:keydown={() => {}}
              >
                {child.name}
              </div>
            {/if}
          {/each}
        </div>
      {/if}
    {:else}
      <div
        class="file"
        class:selected={entry.path === selectedPath}
        on:click={() => open(entry.path, entry.name)}
        on:keydown={() => {}}
      >
        {entry.name}
      </div>
    {/if}
  {/each}
</div>

<style>
  .filetree-content { font-size: 13px; user-select: none; }
  .dir { cursor: pointer; padding: 2px 4px; font-weight: bold; }
  .dir:hover, .file:hover { background: var(--lm-highlight, #292e42); border-radius: 3px; }
  .file { cursor: pointer; padding: 2px 4px 2px 16px; }
  .file.selected { background: var(--lm-highlight, #292e42); border-radius: 3px; }
  .children { padding-left: 12px; }
  .icon { font-size: 10px; margin-right: 4px; }
</style>
```

- [ ] **Step 2: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/FileTree.svelte
git commit -m "feat(desktop): update FileTree to use stores with selected highlight"
```

---

## Task 9: Rewrite App.svelte with new layout

**Files:**
- Modify: `cmd/lm-desktop/frontend/src/App.svelte`

- [ ] **Step 1: Rewrite App.svelte**

Replace the entire file:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { NeedsSetup, OpenFile } from '../wailsjs/go/wailsplugin/App';
  import SetupWizard from './components/SetupWizard.svelte';
  import IconRibbon from './components/IconRibbon.svelte';
  import FileTree from './components/FileTree.svelte';
  import TabBar from './components/TabBar.svelte';
  import NavBar from './components/NavBar.svelte';
  import ContentArea from './components/ContentArea.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import CommandBar from './components/CommandBar.svelte';
  import { applyTheme } from './lib/theme';
  import { sidebarOpen, viewMode, toggleEditMode, goBack, goForward, closeTab, activeTabIndex, tabs } from './lib/stores';

  let needsSetup: boolean | null = null;
  let commandBar: CommandBar;

  function handleKeydown(e: KeyboardEvent) {
    if (needsSetup !== false) return;
    const mod = e.metaKey || e.ctrlKey;

    if (mod && e.key === 'e') {
      e.preventDefault();
      toggleEditMode();
      return;
    }
    if (mod && e.key === 'b') {
      e.preventDefault();
      sidebarOpen.update(v => !v);
      return;
    }
    if (mod && e.key === 'p') {
      e.preventDefault();
      commandBar.open();
      return;
    }
    if (mod && e.key === 'w') {
      e.preventDefault();
      let idx: number;
      activeTabIndex.subscribe(v => idx = v)();
      if (idx >= 0) closeTab(idx);
      return;
    }
    if (mod && e.key === '[') {
      e.preventDefault();
      const path = goBack();
      if (path) OpenFile(path);
      return;
    }
    if (mod && e.key === ']') {
      e.preventDefault();
      const path = goForward();
      if (path) OpenFile(path);
      return;
    }
  }

  function onSetupComplete() {
    needsSetup = false;
    applyTheme();
  }

  let showSidebar: boolean;
  $: showSidebar = $sidebarOpen;

  onMount(async () => {
    needsSetup = await NeedsSetup();
    if (!needsSetup) {
      applyTheme();
    }
  });
</script>

<svelte:window on:keydown={handleKeydown} />

{#if needsSetup === null}
  <!-- Loading -->
{:else if needsSetup}
  <SetupWizard on:complete={onSetupComplete} />
{:else}
  <CommandBar bind:this={commandBar} />

  <div class="app" class:hide-sidebar={!showSidebar}>
    <div class="ribbon">
      <IconRibbon />
    </div>

    {#if showSidebar}
      <aside class="sidebar">
        <FileTree />
      </aside>
    {/if}

    <div class="main">
      <TabBar />
      <NavBar />
      <ContentArea />
    </div>

    <footer class="statusbar">
      <StatusBar />
    </footer>
  </div>
{/if}

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background: var(--lm-bg, #1a1b26);
    color: var(--lm-fg, #c0caf5);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }

  .app {
    display: grid;
    grid-template-columns: 40px 250px 1fr;
    grid-template-rows: 1fr auto;
    height: 100vh;
    overflow: hidden;
  }

  .app.hide-sidebar {
    grid-template-columns: 40px 1fr;
  }

  .ribbon {
    grid-row: 1 / -1;
  }

  .sidebar {
    grid-row: 1;
    border-right: 1px solid var(--lm-border, #3b4261);
    overflow-y: auto;
    padding: 8px;
  }

  .main {
    grid-row: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .statusbar {
    grid-column: 2 / -1;
    grid-row: 2;
    padding: 4px 12px;
    background: var(--lm-statusbar-bg, #16161e);
    border-top: 1px solid var(--lm-border, #3b4261);
    font-size: 12px;
  }

  .app.hide-sidebar .statusbar {
    grid-column: 2 / -1;
  }
</style>
```

- [ ] **Step 2: Build and verify**

Run: `cd cmd/lm-desktop && export PATH="/usr/local/go/bin:$PATH" && wails build 2>&1 | tail -5`
Expected: Build succeeds

- [ ] **Step 3: Commit**

```bash
git add cmd/lm-desktop/frontend/src/App.svelte
git commit -m "feat(desktop): rewrite App.svelte with Obsidian-like layout"
```

---

## Task 10: Integration test — build and smoke test

**Files:** None (verification only)

- [ ] **Step 1: Run Go tests**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/wailsplugin/ -v`
Expected: All PASS

- [ ] **Step 2: Build the desktop app**

Run: `cd cmd/lm-desktop && export PATH="/usr/local/go/bin:$PATH" && wails build 2>&1 | tail -10`
Expected: Build succeeds

- [ ] **Step 3: Manual smoke test**

Launch the app and verify:
1. Icon ribbon appears on the left with file/graph/settings icons
2. File tree sidebar shows files
3. Clicking a file opens it in reading mode (rendered markdown)
4. Tab appears in tab bar
5. Cmd+E toggles to editing mode (CodeMirror)
6. Cmd+E toggles back to reading mode
7. Cmd+B hides/shows sidebar
8. Graph icon switches to brain graph view
9. Status bar shows backlinks, word count, character count
10. Back/forward navigation works

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "feat(desktop): complete Obsidian-like UI redesign"
```

# Obsidian-like UI Redesign

**Date:** 2026-03-26
**Status:** Approved
**Scope:** Wails desktop GUI (`cmd/lm-desktop/frontend/`)

## Goal

Redesign the Wails desktop frontend from a 3-column split layout (file tree | editor | preview) to an Obsidian-inspired single-panel layout with reading/editing mode toggle.

## Layout

```
┌──┬──────────┬─────────────────────────────────────────┐
│  │ Excali.. │  Tab1  │ Tab2 x │  Tab3  │      + │     │
│  │ BugList  ├────────┴────────┴────────┴────────┴─────│
│  │ ff       │  ← →     BugList              📖  ···   │
│  │ growth.. ├─────────────────────────────────────────│
│  │ Inspir.. │                                         │
│  │ Untitl.. │  # BugList                              │
│  │ V1 Fea.. │                                         │
│☐ │          │  ☐ md files are not opening              │
│🔗│          │  ☐                                       │
│⚙ │          │                                         │
│  │          │                                         │
├──┴──────────┴─────────────────────────────────────────┤
│ 0 backlinks   ✏ 7 words  37 characters                │
└───────────────────────────────────────────────────────┘
```

### 1. Icon Ribbon (leftmost, ~40px)

Vertical strip of icon buttons. Fixed width, full height.

| Icon | Action |
|------|--------|
| 📄 (file) | Show/hide file tree sidebar |
| 🔗 (graph) | Switch content area to brain graph view |
| ⚙ (settings) | Open settings (future) |

Active icon highlighted with left accent bar.

### 2. File Tree Sidebar (~250px)

- Same data source as current `FileTree.svelte`
- Collapsible via icon ribbon click or Cmd+B
- Selected file highlighted
- Vault name at bottom (like Obsidian's vault switcher area)

### 3. Content Area

#### Tab Bar
- Horizontal tabs for open files
- Each tab: file name + close (x) button
- Active tab highlighted
- "+" button to create new file (future)
- Scrollable if many tabs

#### Navigation Bar (below tabs)
- Left: back/forward arrows (navigation history)
- Center: current file title
- Right: reading/editing mode toggle icon, overflow menu (···)

#### Content Panel
Two modes, toggled with Cmd+E:

**Reading Mode (default):**
- Rendered markdown using current remark/rehype pipeline
- Clickable checkboxes (toggle tasks via `UpdateTask` binding)
- Clicking text area switches to editing mode at that position (optional enhancement)

**Editing Mode:**
- CodeMirror 6 editor (current `Editor.svelte`)
- Full height, same dark theme
- Cmd+S to save
- Cmd+E to return to reading mode

**Brain Graph View:**
- Activated from icon ribbon (not a tab)
- Replaces content panel entirely
- Current `BrainGraph.svelte` component
- Click node to navigate to that file

### 4. Status Bar (bottom, full width)

- Left: backlinks count
- Right: word count, character count
- Shows editing indicator (pencil icon) in editing mode

## State Management

### Tab State
```typescript
interface Tab {
  path: string;
  name: string;
  modified: boolean;
}

let openTabs: Tab[] = [];
let activeTabIndex: number = 0;
```

### Navigation History
```typescript
let history: string[] = [];   // file paths
let historyIndex: number = -1;
```

### View Mode
```typescript
type ViewMode = 'reading' | 'editing' | 'graph';
let viewMode: ViewMode = 'reading';
```

## Components (new/modified)

| Component | Status | Description |
|-----------|--------|-------------|
| `App.svelte` | **Rewrite** | New grid layout with ribbon + sidebar + content |
| `IconRibbon.svelte` | **New** | Vertical icon strip |
| `TabBar.svelte` | **New** | Horizontal file tabs |
| `NavBar.svelte` | **New** | Back/forward, title, mode toggle |
| `ContentArea.svelte` | **New** | Orchestrates reading/editing/graph modes |
| `FileTree.svelte` | **Modify** | Styling updates, selected file highlight |
| `StatusBar.svelte` | **Modify** | Backlinks, word count, character count |
| `Editor.svelte` | **Keep** | Reused as-is inside ContentArea |
| `Preview.svelte` | **Keep** | Reused as-is inside ContentArea (reading mode) |
| `BrainGraph.svelte` | **Keep** | Reused as-is inside ContentArea |
| `CommandBar.svelte` | **Keep** | Stays as overlay |
| `SetupWizard.svelte` | **Keep** | No changes |

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Cmd+E | Toggle reading/editing mode |
| Cmd+B | Toggle file tree sidebar |
| Cmd+P | Open command palette |
| Cmd+W | Close active tab |
| Ctrl+Tab | Next tab |
| Ctrl+Shift+Tab | Previous tab |
| Cmd+[ | Navigate back |
| Cmd+] | Navigate forward |

## Go Backend Changes

### New WailsApp Methods
- `GetWordCount() int` — word count for status bar
- `GetCharCount() int` — character count for status bar
- `GetBacklinkCount() int` — backlink count for current file

### Modified
- `GetStatus()` — extend to include word/char/backlink counts (or keep separate methods)

## Styling

- Dark theme: keep Tokyo Night base colors
- Icon ribbon: slightly darker than sidebar (`--lm-ribbon-bg`)
- Sidebar: current `--lm-bg`
- Content area: slightly lighter for contrast
- Tab bar: similar to Obsidian's muted tabs with active highlight
- Font: keep monospace for editor, use system sans-serif for UI chrome

## Tab Lifecycle

- **Create:** Opening a file (file tree click, wiki-link, command palette) creates a tab. If file is already open, focus existing tab.
- **Close:** Cmd+W closes active tab. No unsaved-change prompt (auto-save on switch).
- **Modified indicator:** Show a dot (•) after the file name in the tab.
- **No ephemeral/preview tabs** — every open is a persistent tab.

## Graph View Transitions

- Clicking graph icon in ribbon sets `viewMode = 'graph'`
- Cmd+E in graph view: no-op (only works in reading/editing)
- Clicking a file in file tree or graph node: exits graph view, sets `viewMode = 'reading'`
- Back/forward navigation: exits graph view if navigating to a file

## CodeMirror Lifecycle

Keep CodeMirror mounted but hidden (`display: none`) when in reading mode. Sync content on mode switch. This preserves cursor position and undo history.

## Navigation History

- Only explicit file opens push to history (file tree click, wiki-link, back/forward)
- Switching tabs does NOT push to history
- Graph view is not a history entry

## Out of Scope

- WYSIWYG editing
- Click-to-edit at cursor position
- Drag-and-drop tabs
- Split panes
- Plugin sidebar panels
- Settings panel UI
- "+" new file button (tab bar shows it but is non-functional)

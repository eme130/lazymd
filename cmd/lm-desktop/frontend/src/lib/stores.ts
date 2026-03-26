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

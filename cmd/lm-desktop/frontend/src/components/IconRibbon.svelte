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

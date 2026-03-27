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

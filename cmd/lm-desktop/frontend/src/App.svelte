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

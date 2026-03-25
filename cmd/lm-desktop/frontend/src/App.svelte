<script lang="ts">
  import { onMount } from 'svelte';
  import { NeedsSetup } from '../wailsjs/go/wailsplugin/App';
  import SetupWizard from './components/SetupWizard.svelte';
  import Editor from './components/Editor.svelte';
  import FileTree from './components/FileTree.svelte';
  import Preview from './components/Preview.svelte';
  import BrainGraph from './components/BrainGraph.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import CommandBar from './components/CommandBar.svelte';
  import { applyTheme } from './lib/theme';

  let needsSetup: boolean | null = null;
  let activeRightPanel: 'preview' | 'brain' = 'preview';
  let showFileTree = true;
  let commandBar: CommandBar;

  function handleKeydown(e: KeyboardEvent) {
    if (needsSetup !== false) return;
    if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.key === 'p') {
      e.preventDefault();
      activeRightPanel = activeRightPanel === 'preview' ? 'brain' : 'preview';
      return;
    }
    if ((e.metaKey || e.ctrlKey) && e.key === 'p') {
      e.preventDefault();
      commandBar.open();
      return;
    }
    if ((e.metaKey || e.ctrlKey) && e.key === 'b') {
      e.preventDefault();
      showFileTree = !showFileTree;
    }
  }

  function onSetupComplete() {
    needsSetup = false;
    applyTheme();
  }

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

  <div class="app" class:hide-filetree={!showFileTree}>
    {#if showFileTree}
      <aside class="filetree">
        <FileTree />
      </aside>
    {/if}

    <main class="editor">
      <Editor />
    </main>

    <aside class="right-panel">
      <div class="panel-tabs">
        <button class:active={activeRightPanel === 'preview'} on:click={() => activeRightPanel = 'preview'}>Preview</button>
        <button class:active={activeRightPanel === 'brain'} on:click={() => activeRightPanel = 'brain'}>Brain</button>
      </div>
      {#if activeRightPanel === 'preview'}
        <Preview />
      {:else}
        <BrainGraph />
      {/if}
    </aside>

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
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
  }

  .app {
    display: grid;
    grid-template-columns: 250px 1fr 350px;
    grid-template-rows: 1fr auto;
    height: 100vh;
    overflow: hidden;
  }

  .app.hide-filetree {
    grid-template-columns: 1fr 350px;
  }

  .filetree {
    grid-row: 1;
    border-right: 1px solid var(--lm-border, #3b4261);
    overflow-y: auto;
    padding: 8px;
  }

  .editor {
    grid-row: 1;
    overflow: hidden;
  }

  .right-panel {
    grid-row: 1;
    border-left: 1px solid var(--lm-border, #3b4261);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .panel-tabs {
    display: flex;
    border-bottom: 1px solid var(--lm-border, #3b4261);
  }

  .panel-tabs button {
    flex: 1;
    padding: 6px 12px;
    background: transparent;
    border: none;
    color: var(--lm-fg, #c0caf5);
    cursor: pointer;
    font-size: 12px;
    opacity: 0.6;
  }

  .panel-tabs button.active {
    opacity: 1;
    border-bottom: 2px solid var(--lm-link, #7aa2f7);
  }

  .statusbar {
    grid-column: 1 / -1;
    grid-row: 2;
    padding: 4px 12px;
    background: var(--lm-statusbar-bg, #16161e);
    border-top: 1px solid var(--lm-border, #3b4261);
    font-size: 12px;
  }
</style>

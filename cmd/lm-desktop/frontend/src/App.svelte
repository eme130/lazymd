<script lang="ts">
  let activeRightPanel: 'preview' | 'brain' = 'preview';
  let showFileTree = true;
</script>

<div class="app" class:hide-filetree={!showFileTree}>
  {#if showFileTree}
    <aside class="filetree">
      <p>File Tree (placeholder)</p>
    </aside>
  {/if}

  <main class="editor">
    <p>Editor (placeholder)</p>
  </main>

  <aside class="right-panel">
    {#if activeRightPanel === 'preview'}
      <p>Preview (placeholder)</p>
    {:else}
      <p>Brain Graph (placeholder)</p>
    {/if}
  </aside>

  <footer class="statusbar">
    <span>Status Bar (placeholder)</span>
  </footer>
</div>

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
    overflow-y: auto;
    padding: 8px;
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

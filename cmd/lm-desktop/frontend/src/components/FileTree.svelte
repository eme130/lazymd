<script lang="ts">
  import { onMount } from 'svelte';
  import { ListFiles, OpenFile } from '../../wailsjs/go/wailsplugin/App';
  import { onFileOpened, onFileSaved } from '../lib/events';

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

  async function open(path: string) {
    await OpenFile(path);
  }

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
              <div class="file" on:click={() => open(child.path)} on:keydown={() => {}}>
                {child.name}
              </div>
            {/if}
          {/each}
        </div>
      {/if}
    {:else}
      <div class="file" on:click={() => open(entry.path)} on:keydown={() => {}}>
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
  .children { padding-left: 12px; }
  .icon { font-size: 10px; margin-right: 4px; }
</style>

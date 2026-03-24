<script lang="ts">
  import { ListFiles, OpenFile } from '../../../wailsjs/go/wailsplugin/App';

  export let visible = false;

  interface CommandItem {
    label: string;
    action: () => void;
  }

  let query = '';
  let items: CommandItem[] = [];
  let filtered: CommandItem[] = [];
  let selected = 0;
  let input: HTMLInputElement;

  export async function open() {
    visible = true;
    query = '';
    selected = 0;
    const files = await ListFiles('.');
    items = flattenFiles(files).map((f) => ({
      label: f.path,
      action: () => { OpenFile(f.path); close(); },
    }));
    filtered = items;
    setTimeout(() => input?.focus(), 10);
  }

  function flattenFiles(entries: any[]): any[] {
    const result: any[] = [];
    for (const e of entries) {
      if (e.isDir && e.children) {
        result.push(...flattenFiles(e.children));
      } else if (!e.isDir) {
        result.push(e);
      }
    }
    return result;
  }

  function close() { visible = false; }

  function filter() {
    const q = query.toLowerCase();
    filtered = items.filter((i) => i.label.toLowerCase().includes(q));
    selected = 0;
  }

  function handleKey(e: KeyboardEvent) {
    if (e.key === 'Escape') { close(); return; }
    if (e.key === 'ArrowDown') { selected = Math.min(selected + 1, filtered.length - 1); e.preventDefault(); return; }
    if (e.key === 'ArrowUp') { selected = Math.max(selected - 1, 0); e.preventDefault(); return; }
    if (e.key === 'Enter' && filtered[selected]) { filtered[selected].action(); return; }
  }
</script>

{#if visible}
  <div class="overlay" on:click={close} on:keydown={() => {}}>
    <div class="command-bar" on:click|stopPropagation on:keydown={() => {}}>
      <input
        bind:this={input}
        bind:value={query}
        on:input={filter}
        on:keydown={handleKey}
        placeholder="Search files..."
      />
      <div class="results">
        {#each filtered.slice(0, 20) as item, i}
          <div
            class="result"
            class:selected={i === selected}
            on:click={item.action}
            on:keydown={() => {}}
          >
            {item.label}
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .overlay { position: fixed; inset: 0; background: rgba(0, 0, 0, 0.5); display: flex; justify-content: center; padding-top: 20vh; z-index: 100; }
  .command-bar { width: 500px; max-height: 400px; background: var(--lm-bg, #1a1b26); border: 1px solid var(--lm-border, #3b4261); border-radius: 8px; overflow: hidden; }
  input { width: 100%; padding: 12px 16px; background: transparent; border: none; border-bottom: 1px solid var(--lm-border, #3b4261); color: var(--lm-fg, #c0caf5); font-size: 14px; outline: none; box-sizing: border-box; }
  .results { overflow-y: auto; max-height: 320px; }
  .result { padding: 8px 16px; cursor: pointer; font-size: 13px; }
  .result:hover, .result.selected { background: var(--lm-highlight, #292e42); }
</style>

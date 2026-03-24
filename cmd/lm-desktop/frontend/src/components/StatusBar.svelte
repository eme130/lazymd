<script lang="ts">
  import { onMount } from 'svelte';
  import { GetStatus } from '../../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged, onCursorMoved, onFileOpened, onFileSaved } from '../lib/events';

  let status = { file: '', lines: 0, row: 0, col: 0, modified: false };

  async function refresh() {
    status = await GetStatus();
  }

  onMount(() => {
    refresh();
    onBufferChanged(() => refresh());
    onCursorMoved(() => refresh());
    onFileOpened(() => refresh());
    onFileSaved(() => refresh());
  });
</script>

<div class="status">
  <span class="file">{status.file || 'No file'}{status.modified ? ' *' : ''}</span>
  <span class="pos">Ln {status.row + 1}, Col {status.col + 1}</span>
  <span class="lines">{status.lines} lines</span>
</div>

<style>
  .status { display: flex; gap: 16px; align-items: center; }
  .file { flex: 1; }
  .pos, .lines { opacity: 0.7; }
</style>

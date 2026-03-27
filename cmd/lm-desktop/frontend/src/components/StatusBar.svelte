<script lang="ts">
  import { onMount } from 'svelte';
  import { GetWordCount, GetCharCount, GetBacklinkCount } from '../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged, onFileOpened, onFileSaved } from '../lib/events';
  import { viewMode } from '../lib/stores';

  let wordCount = 0;
  let charCount = 0;
  let backlinkCount = 0;

  async function refresh() {
    [wordCount, charCount, backlinkCount] = await Promise.all([
      GetWordCount(),
      GetCharCount(),
      GetBacklinkCount(),
    ]);
  }

  onMount(() => {
    refresh();
    onBufferChanged(() => refresh());
    onFileOpened(() => refresh());
    onFileSaved(() => refresh());
  });
</script>

<div class="status">
  <span class="left">
    <span class="backlinks">{backlinkCount} backlinks</span>
  </span>
  <span class="right">
    {#if $viewMode === 'editing'}
      <span class="indicator">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
      </span>
    {/if}
    <span>{wordCount} words</span>
    <span>{charCount} characters</span>
  </span>
</div>

<style>
  .status {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
    opacity: 0.6;
  }
  .left, .right { display: flex; gap: 16px; align-items: center; }
  .indicator { display: flex; align-items: center; }
</style>

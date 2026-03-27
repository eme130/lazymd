<script lang="ts">
  import { viewMode } from '../lib/stores';
  import Editor from './Editor.svelte';
  import Preview from './Preview.svelte';
  import BrainGraph from './BrainGraph.svelte';
</script>

<div class="content-area">
  {#if $viewMode === 'graph'}
    <BrainGraph />
  {:else}
    <div class="view" class:hidden={$viewMode !== 'reading'}>
      <Preview />
    </div>
    <div class="view" class:hidden={$viewMode !== 'editing'}>
      <Editor />
    </div>
  {/if}
</div>

<style>
  .content-area {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .view { flex: 1; overflow: hidden; }
  .view.hidden { display: none; }
</style>

<script lang="ts">
  import { tabs, activeTabIndex, closeTab } from '../lib/stores';

  function selectTab(index: number) {
    activeTabIndex.set(index);
  }

  function onClose(e: MouseEvent, index: number) {
    e.stopPropagation();
    closeTab(index);
  }
</script>

<div class="tab-bar">
  {#each $tabs as tab, i}
    <button
      class="tab"
      class:active={i === $activeTabIndex}
      on:click={() => selectTab(i)}
    >
      <span class="tab-name">{tab.name}{tab.modified ? ' •' : ''}</span>
      <span class="tab-close" on:click={(e) => onClose(e, i)} on:keydown={() => {}}>×</span>
    </button>
  {/each}
</div>

<style>
  .tab-bar {
    display: flex;
    background: var(--lm-ribbon-bg, #13141c);
    border-bottom: 1px solid var(--lm-border, #3b4261);
    overflow-x: auto;
    min-height: 35px;
    -webkit-app-region: drag;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: transparent;
    border: none;
    border-right: 1px solid var(--lm-border, #3b4261);
    color: var(--lm-fg, #c0caf5);
    opacity: 0.5;
    cursor: pointer;
    font-size: 12px;
    white-space: nowrap;
    -webkit-app-region: no-drag;
  }

  .tab:hover {
    opacity: 0.8;
  }

  .tab.active {
    opacity: 1;
    background: var(--lm-bg, #1a1b26);
  }

  .tab-close {
    font-size: 14px;
    line-height: 1;
    opacity: 0;
    border-radius: 3px;
    padding: 0 2px;
  }

  .tab:hover .tab-close {
    opacity: 0.6;
  }

  .tab-close:hover {
    opacity: 1 !important;
    background: var(--lm-highlight, #292e42);
  }

  .tab-name {
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>

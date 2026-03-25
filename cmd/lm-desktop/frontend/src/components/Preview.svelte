<script lang="ts">
  import { onMount } from 'svelte';
  import { remark } from 'remark';
  import remarkHtml from 'remark-html';
  import { GetContent } from '../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged } from '../lib/events';

  let html = '';
  let debounceTimer: ReturnType<typeof setTimeout>;

  async function render(content: string) {
    const result = await remark().use(remarkHtml).process(content);
    html = String(result);
  }

  function debouncedRender(content: string) {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => render(content), 150);
  }

  onMount(async () => {
    const content = await GetContent();
    await render(content);

    onBufferChanged((data) => {
      if (data?.content != null) {
        debouncedRender(data.content);
      }
    });
  });
</script>

<div class="preview-content">
  {@html html}
</div>

<style>
  .preview-content { padding: 16px; line-height: 1.6; overflow-y: auto; height: 100%; }
  .preview-content :global(h1) { font-size: 1.8em; margin: 0.5em 0; color: var(--lm-heading, #bb9af7); }
  .preview-content :global(h2) { font-size: 1.4em; margin: 0.5em 0; color: var(--lm-heading, #bb9af7); }
  .preview-content :global(h3) { font-size: 1.2em; margin: 0.5em 0; color: var(--lm-heading, #bb9af7); }
  .preview-content :global(code) { background: var(--lm-code-bg, #292e42); padding: 2px 6px; border-radius: 3px; font-size: 0.9em; }
  .preview-content :global(pre) { background: var(--lm-code-bg, #292e42); padding: 12px; border-radius: 6px; overflow-x: auto; }
  .preview-content :global(a) { color: var(--lm-link, #7aa2f7); }
  .preview-content :global(blockquote) { border-left: 3px solid var(--lm-border, #3b4261); padding-left: 12px; opacity: 0.8; }
</style>

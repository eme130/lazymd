<script lang="ts">
  import { onMount } from 'svelte';
  import { unified } from 'unified';
  import remarkParse from 'remark-parse';
  import remarkMath from 'remark-math';
  import remarkRehype from 'remark-rehype';
  import rehypeSanitize, { defaultSchema } from 'rehype-sanitize';
  import rehypeKatex from 'rehype-katex';
  import rehypeStringify from 'rehype-stringify';
  import 'katex/dist/katex.min.css';
  import { GetContent } from '../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged } from '../lib/events';

  let html = '';
  let debounceTimer: ReturnType<typeof setTimeout>;

  // Custom sanitize schema that preserves math class names on code elements.
  // remark-math produces <code class="math-inline"> and <code class="math-display">
  // which rehype-katex needs to find and render.
  const sanitizeSchema = {
    ...defaultSchema,
    attributes: {
      ...defaultSchema.attributes,
      code: [
        ...(defaultSchema.attributes?.code || []),
        ['className', 'math-inline', 'math-display']
      ]
    }
  };

  async function render(content: string) {
    const result = await unified()
      .use(remarkParse)
      .use(remarkMath)
      .use(remarkRehype)
      .use(rehypeSanitize, sanitizeSchema)
      .use(rehypeKatex)
      .use(rehypeStringify)
      .process(content);
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
  /* KaTeX math display */
  .preview-content :global(.katex-display) {
    margin: 1em 0;
    text-align: center;
    overflow-x: auto;
    overflow-y: hidden;
  }
  .preview-content :global(.katex) {
    font-size: 1.1em;
  }
</style>

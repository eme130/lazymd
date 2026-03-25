# Math Preview Panel Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add KaTeX math rendering to the Wails desktop preview panel so `$...$` and `$$...$$` display as typeset equations.

**Architecture:** Replace the current `remark` → `remark-html` pipeline in `Preview.svelte` with a `unified` → `remark-parse` → `remark-math` → `remark-rehype` → `rehype-sanitize` (custom schema) → `rehype-katex` → `rehype-stringify` pipeline. Sanitize runs before KaTeX so user HTML is cleaned but math nodes survive, then KaTeX renders trusted output. Import KaTeX CSS for equation styling. No Go changes.

**Tech Stack:** Svelte 4, unified, remark-math, rehype-katex, KaTeX

**Spec:** `docs/superpowers/specs/2026-03-25-math-preview-design.md`

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `cmd/lm-desktop/frontend/package.json` | Modify | Add/remove npm dependencies |
| `cmd/lm-desktop/frontend/src/components/Preview.svelte` | Modify | Replace remark-html pipeline with unified+rehype+katex pipeline |

---

### Task 1: Update NPM Dependencies

**Files:**
- Modify: `cmd/lm-desktop/frontend/package.json`

- [ ] **Step 1: Install new dependencies**

Run from the frontend directory:

```bash
cd /Users/eme130/projects/lazymd/cmd/lm-desktop/frontend && npm install unified remark-parse remark-math remark-rehype rehype-katex rehype-stringify katex
```

This adds:
- `unified` — core processor (was transitive via `remark`, now direct)
- `remark-parse` — markdown parser (was transitive via `remark`, now direct)
- `remark-math` — parses `$...$` and `$$...$$` into math AST nodes
- `remark-rehype` — bridges remark AST to rehype AST
- `rehype-katex` — renders math AST nodes as KaTeX HTML
- `rehype-stringify` — serializes rehype AST to HTML string
- `katex` — KaTeX library (peer dep of rehype-katex, provides CSS + fonts)

- [ ] **Step 2: Remove old dependencies**

```bash
cd /Users/eme130/projects/lazymd/cmd/lm-desktop/frontend && npm uninstall remark remark-html
```

This removes:
- `remark` — replaced by `unified` + `remark-parse`
- `remark-html` — replaced by `remark-rehype` + `rehype-stringify`

Note: `rehype-sanitize` is already in `package.json` and stays.

- [ ] **Step 3: Verify the final package.json dependencies**

Read `cmd/lm-desktop/frontend/package.json` and verify:
- `unified`, `remark-parse`, `remark-math`, `remark-rehype`, `rehype-katex`, `rehype-stringify`, `katex` are present
- `remark` and `remark-html` are gone
- `rehype-sanitize`, `codemirror`, `d3`, etc. are untouched

- [ ] **Step 4: Verify build still works**

```bash
cd /Users/eme130/projects/lazymd/cmd/lm-desktop/frontend && npm run build
```

Expected: Build succeeds (Preview.svelte still imports old modules, but tree-shaking means it'll warn or error — that's expected and fixed in Task 2).

Note: The build may fail because Preview.svelte still imports `remark` and `remark-html` which are now uninstalled. That's OK — Task 2 fixes the imports. If it fails, skip this verification and proceed to Task 2.

- [ ] **Step 5: Commit**

```bash
cd /Users/eme130/projects/lazymd && git add cmd/lm-desktop/frontend/package.json cmd/lm-desktop/frontend/package-lock.json
git commit -m "build(desktop): add KaTeX and unified pipeline dependencies, remove remark-html

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

### Task 2: Replace Preview Pipeline with KaTeX Rendering

**Files:**
- Modify: `cmd/lm-desktop/frontend/src/components/Preview.svelte`

- [ ] **Step 1: Replace the script section**

Replace the entire `<script>` block in `Preview.svelte` with:

```svelte
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
```

Key changes:
- `remark` + `remarkHtml` replaced with `unified` + 6 plugins
- KaTeX CSS imported via `import 'katex/dist/katex.min.css'`
- `rehypeSanitize` runs BEFORE `rehypeKatex` — sanitizes user HTML first, then KaTeX renders trusted math output
- Custom sanitize schema preserves `math-inline` and `math-display` class names on `code` elements so KaTeX can find them
- Everything else (debounce, events, onMount) identical

- [ ] **Step 2: Add KaTeX-specific CSS styles**

Add these CSS rules to the existing `<style>` block in `Preview.svelte`, after the existing rules:

```css
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
```

These ensure display math (`$$...$$`) is centered and inline math is slightly larger than body text.

- [ ] **Step 3: Verify build**

```bash
cd /Users/eme130/projects/lazymd/cmd/lm-desktop/frontend && npm run build
```

Expected: Build succeeds with no errors.

- [ ] **Step 4: Manual smoke test**

If Wails dev mode is available:

```bash
cd /Users/eme130/projects/lazymd && export PATH="/usr/local/go/bin:$PATH" && task desktop:dev
```

Create or open a markdown file with math content:

```markdown
# Test Math

Inline math: $x^2 + y^2 = z^2$

Display math:

$$\frac{a}{b} + \sum_{i=1}^{n} x_i$$

Code block (should NOT render as math):

```python
price = $5
```

Normal text after math.
```

Verify:
- `$x^2 + y^2 = z^2$` renders as inline typeset math
- `$$\frac{a}{b}...$$` renders as centered display math
- Code block shows literal `$5` (not math)
- Headings, links, code blocks render as before

If Wails dev mode is not available (missing Wails CLI), just verify the build succeeds — the pipeline is standard and well-tested upstream.

- [ ] **Step 5: Commit**

```bash
cd /Users/eme130/projects/lazymd && git add cmd/lm-desktop/frontend/src/components/Preview.svelte
git commit -m "feat(desktop): add KaTeX math rendering to preview panel

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

# Math Preview Panel Design

**Date**: 2026-03-25
**Status**: Approved
**Subsystem**: 2 of 4 (LaTeX Editor → Math Preview → Python REPL → Theorem Brain)

## Problem

LazyMD's desktop preview panel renders markdown to HTML but doesn't render `$...$` or `$$...$$` math zones. A mathematics PhD student needs to see typeset equations while writing.

## Design

### Approach: KaTeX in the Svelte Frontend

Replace the current `remark` → `remark-html` pipeline in `Preview.svelte` with a unified rehype pipeline that includes KaTeX math rendering.

**Current pipeline:**
```
remark() → remarkHtml → HTML string
```

**New pipeline:**
```
unified() → remarkParse → remarkMath → remarkRehype → rehypeSanitize (custom schema) → rehypeKatex → rehypeStringify → HTML string
```

Note: `rehype-sanitize` is already in `package.json` but unused. It's now wired into the pipeline BEFORE KaTeX — sanitizes user HTML first, then KaTeX renders trusted math output. A custom schema preserves `math-inline` and `math-display` class names on `code` elements so KaTeX can find and render them.

### NPM Dependencies

**Add:**
- `remark-math` — parses `$...$` and `$$...$$` into math AST nodes
- `remark-rehype` — bridges remark (markdown) AST to rehype (HTML) AST
- `rehype-katex` — renders math AST nodes as KaTeX HTML
- `rehype-stringify` — serializes rehype AST to HTML string
- `katex` — peer dependency of rehype-katex (provides CSS and fonts)
- `unified` — core processor (promoted from transitive dep via `remark` to direct dep)
- `remark-parse` — markdown parser plugin for unified (promoted from transitive to direct)

**Remove:**
- `remark` — replaced by `unified` + `remark-parse`
- `remark-html` — replaced by `remark-rehype` + `rehype-stringify`

**Retain (already in package.json):**
- `rehype-sanitize` — now wired into the pipeline (was unused)

### Preview.svelte Changes

**Script:** Replace the `render()` function:

```typescript
import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkMath from 'remark-math';
import remarkRehype from 'remark-rehype';
import rehypeSanitize, { defaultSchema } from 'rehype-sanitize';
import rehypeKatex from 'rehype-katex';
import rehypeStringify from 'rehype-stringify';
import 'katex/dist/katex.min.css';

async function render(content: string) {
  const result = await unified()
    .use(remarkParse)
    .use(remarkMath)
    .use(remarkRehype)
    .use(rehypeSanitize, { ...defaultSchema, attributes: { ...defaultSchema.attributes, code: [...(defaultSchema.attributes?.code || []), ['className', 'math-inline', 'math-display']] } })
    .use(rehypeKatex)
    .use(rehypeStringify)
    .process(content);
  html = String(result);
}
```

**CSS:** Import `katex/dist/katex.min.css` for equation styling (fonts, spacing, fraction layout). KaTeX CSS is self-contained — fonts are bundled or loaded from the katex npm package.

**Everything else unchanged:** Same `onMount`, same `debouncedRender`, same `onBufferChanged` event listener, same 150ms debounce.

### Rendering Behavior

- `$x^2$` → inline typeset math (flows with text)
- `$$\frac{a}{b}$$` → centered display math (own line)
- Non-math markdown → identical to current rendering
- Invalid LaTeX (e.g., `$\invalidcommand$`) → KaTeX shows red error text inline (default rehype-katex behavior, no custom error handling needed)
- Escaped `\$` → literal dollar sign (remark-math handles this)

### TUI Preview

No changes. The TUI preview uses Glamour which doesn't support math rendering. Math zones in the editor are already highlighted with the Math accent color (subsystem 1). The TUI preview passes `$...$` through as literal text.

## Modified Files

- `cmd/lm-desktop/frontend/src/components/Preview.svelte` — replace remark-html pipeline with unified+rehype+katex pipeline, import KaTeX CSS
- `cmd/lm-desktop/frontend/package.json` — add remark-math, remark-rehype, rehype-katex, rehype-stringify, katex, unified, remark-parse; remove remark, remark-html

## Key Test Cases

1. `$x^2$` in markdown → renders inline math (not literal `$x^2$`)
2. `$$\frac{a}{b}$$` on own line → renders centered display math
3. Mixed content: `The equation $E=mc^2$ is famous` → text with inline math
4. Code block containing `$` → not treated as math (remark-math respects code fences)
5. `\$5.00` → literal dollar sign, not math
6. Invalid LaTeX `$\badcommand$` → red error text (KaTeX default)
7. Buffer change → preview updates after 150ms debounce (same as current)
8. Raw HTML `<script>alert(1)</script>` with `$x^2$` → script sanitized, math rendered

## Not In Scope

- TUI math rendering (terminal can't display typeset equations)
- Math rendering in the Go side (go-latex/latex — reserved for future use)
- PDF export with math
- Custom KaTeX macros or configuration
- Syntax highlighting inside code blocks in the preview (separate concern)

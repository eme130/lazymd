# LaTeX Editor Design

**Date**: 2026-03-25
**Status**: Approved
**Subsystem**: 1 of 4 (LaTeX Editor → KaTeX Preview → Python REPL → Theorem Brain)

## Problem

LazyMD has no LaTeX support. A mathematics PhD student needs to write `.tex` files with syntax highlighting, and have `$...$` / `$$...$$` math zones visually distinct inside markdown files.

## User

Kleopetra — mathematics JRF. Writes theorems, proofs, lemmas. Her vault contains paper notes and her own theorem files, interlinked via `[[wiki-links]]`. LazyMD is her complete workspace (replaces Overleaf).

## Design

### A. LaTeX Language Definition

Add `latex` to `internal/highlight/languages.go`:

```go
{
    Name:    "latex",
    Aliases: []string{"tex"},
    Keywords: []string{
        "\\begin", "\\end", "\\section", "\\subsection", "\\subsubsection",
        "\\usepackage", "\\documentclass", "\\newcommand", "\\renewcommand",
        "\\def", "\\let", "\\input", "\\include",
        "\\bibliography", "\\bibliographystyle",
        "\\cite", "\\ref", "\\eqref", "\\label", "\\footnote",
        "\\textbf", "\\textit", "\\emph", "\\underline",
        "\\title", "\\author", "\\date", "\\maketitle",
        "\\tableofcontents", "\\appendix",
        "\\chapter", "\\paragraph", "\\subparagraph",
    },
    Types: []string{
        "theorem", "lemma", "proof", "definition", "corollary",
        "proposition", "remark", "example", "exercise",
        "equation", "align", "gather",
        "figure", "table", "tabular", "enumerate", "itemize",
        "document", "abstract", "cases", "matrix", "pmatrix",
        "bmatrix", "array",
    },
    Builtins: []string{
        "\\frac", "\\int", "\\sum", "\\prod", "\\lim",
        "\\infty", "\\partial", "\\nabla", "\\sqrt",
        "\\forall", "\\exists", "\\in", "\\notin",
        "\\subset", "\\subseteq", "\\supset", "\\supseteq",
        "\\cup", "\\cap", "\\setminus", "\\emptyset",
        "\\mathbb", "\\mathcal", "\\mathfrak", "\\mathrm",
        "\\alpha", "\\beta", "\\gamma", "\\delta", "\\epsilon",
        "\\theta", "\\lambda", "\\mu", "\\pi", "\\sigma",
        "\\omega", "\\phi", "\\psi", "\\rho", "\\tau",
        "\\Gamma", "\\Delta", "\\Theta", "\\Lambda", "\\Sigma",
        "\\Omega", "\\Phi", "\\Psi",
        "\\rightarrow", "\\leftarrow", "\\Rightarrow", "\\Leftarrow",
        "\\mapsto", "\\implies", "\\iff", "\\equiv",
        "\\leq", "\\geq", "\\neq", "\\approx", "\\sim",
        "\\cdot", "\\times", "\\otimes", "\\oplus",
    },
    LineComment: "%",
    // No StringDelimiters — LaTeX doesn't use quoted strings.
    // Brace grouping {…} is not string delimiting; braces tokenize as Normal.
}
```

The highlighter's keyword-based tokenizer already handles `LangDef` — this plugs in directly. Note: `*` variants like `align*` and `gather*` are not in the Types list because the identifier scanner stops at `*`. Environment names containing `*` appear inside `\begin{align*}` — the `align` part matches as a Type, and `*` tokenizes separately as Normal. This is acceptable.

**Backslash command matching:** The existing `builtin.go` tokenizer uses `isIdentStart()` which doesn't recognize `\`. Add a LaTeX-specific scan path in `tokenizeLine()`: when the current language is `latex` and `line[i] == '\\'`, scan forward to collect the full command name (`\frac`, `\int`, etc.), then match against Keywords/Builtins. This is a single `if` block before the normal identifier scan. Only `\command` is matched — brace arguments `{a}{b}` are tokenized normally as separate tokens.

### B. Math Zone Token Types

Extend `internal/markdown/syntax.go` with two new token types:

```go
const (
    // ... existing tokens ...
    MathInline  // $...$
    MathBlock   // $$...$$
)
```

**State machine:** Add `InMathBlock bool` to `LineContext` (alongside existing `InCodeBlock`).

**Tokenization rules:**
- `$$` at line start (with optional leading whitespace, nothing else on the line) opens/closes a math block. Tracked via `InMathBlock` in `LineContext`, same pattern as `InCodeBlock` for code fences. `$$` appearing mid-line (e.g., `text $$x$$ more`) is NOT a block delimiter — it is treated as two consecutive `$` characters, each potentially an inline math delimiter.
- `$` opens/closes inline math within a single line only. No state tracking needed — parsed in `tokenizeInline()` after code spans, before bold/italic.
- Inside math zones, the entire content is one `Span` with type `MathInline` or `MathBlock` (no sub-tokenization of the math content).
- `\$` is an escaped dollar sign, not a math delimiter. Only checked during math zone parsing (not generalized escape handling).
- **Precedence:** Code blocks > math blocks > inline code > inline math. If `InCodeBlock` is true, math zones are ignored. Inline: code spans (`` ` ``) are parsed first, then `$` math within remaining text.
- **Nesting:** `$` inside `$$...$$` is treated as literal (no nesting).
- **Unclosed zones:** Inline `$` with no closing `$` on the same line → treated as literal `$`. Unclosed `$$` block → `InMathBlock` persists until EOF (same as unclosed code fence).

**Color mapping:** Add `Math string` field to `ThemeColors` in `internal/themes/themes.go`. Each theme gets a distinct math accent color. Example: Tokyo Night uses `#bb9af7` (purple — distinct from headings and code keywords).

### C. File Type Detection

In the editor's file-open path, detect `.tex` extension:
- `.tex` files → use LaTeX highlighter for syntax coloring (keyword-based, via `highlight.FindLang("latex")`)
- `.md` / `.rndm` files → use markdown tokenizer (which now includes math zone awareness)

In `internal/editor/editor.go`:
- Add `FileType string` field to `EditorModel` (values: `"markdown"`, `"latex"`)
- In `OpenFile()`, after setting `e.File`, detect extension: `.tex` → `"latex"`, everything else → `"markdown"`
- The UI rendering path (`internal/ui/app.go`) checks `FileType` to decide which tokenizer to use for syntax coloring. Currently `renderEditor()` renders lines as plain text — this spec adds span-based coloring: for each line, get spans from the appropriate tokenizer, then apply lipgloss styles per span type.
- `.tex` files are exclusively LaTeX (no markdown parsing). `.md`/`.rndm` files use the markdown tokenizer with math zone support. No mixed mode.

### D. Website Use Case (deferred)

After all 4 subsystems ship, add a "For Mathematicians" section to the website showing:
- Writing theorems with LaTeX highlighting
- Math zones glowing in markdown notes
- Brain graph connecting `[[Theorem 4.2]]` → `[[Lemma 3.1]]`
- Python REPL verifying conjectures
- KaTeX preview rendering equations

This is documentation, not code — written last.

## Modified Files

- `internal/highlight/languages.go` — add LaTeX language definition
- `internal/highlight/builtin.go` — add LaTeX backslash command scan path in `tokenizeLine()`
- `internal/highlight/highlight_test.go` — tests for LaTeX tokenization (test cases 11–13)
- `internal/markdown/syntax.go` — add `MathInline`, `MathBlock` token types + tokenization in `tokenizeInline()` and `TokenizeLine()`
- `internal/markdown/syntax_test.go` — tests for math zone detection (test cases 1–10)
- `internal/themes/themes.go` — add `Math` color to `ThemeColors`, set in all themes
- `internal/editor/editor.go` — add `FileType` field, set on `OpenFile`
- `internal/ui/app.go` — update `renderEditor()` to apply syntax coloring via spans from the appropriate tokenizer based on `FileType`

## Key Test Cases

**Math zone tokenization (`syntax_test.go`):**
1. `$x^2$` → one `MathInline` span covering `$x^2$`
2. `text $a+b$ more` → Normal, MathInline, Normal
3. `$$` on its own line → opens/closes MathBlock (like code fence)
4. Multi-line: `$$` / `x = 1` / `$$` → MathBlock spans for content lines
5. `price is $5` (no closing `$`) → treated as literal, Normal token
6. `\$escaped` → Normal token (not math)
7. `` `code $not math$` `` → CodeInline only, no MathInline
8. Inside code block → no math tokenization
9. `$$x$$ and $y$` → `$$` mid-line is NOT a block delimiter, so parsed as inline: `$x$` is MathInline, ` and ` is Normal, `$y$` is MathInline (the leading/trailing `$` pair up as inline math delimiters)
10. Unclosed `$$` at EOF → MathBlock persists

**LaTeX highlighting (`highlight_test.go`):**
11. `\frac{a}{b}` → `\frac` as Builtin, `{`, `a`, `}`, `{`, `b`, `}` as Normal
12. `% comment` → Comment token
13. `\begin{theorem}` → `\begin` as Keyword, `theorem` as Type

## Not In Scope

- Math rendering/preview (subsystem 2 — KaTeX Preview Panel)
- Python execution (subsystem 3 — Python REPL)
- Theorem-aware brain links (subsystem 4)
- PDF export (future)
- `.tex` file compilation

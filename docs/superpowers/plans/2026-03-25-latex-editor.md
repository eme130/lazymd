# LaTeX Editor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add LaTeX syntax highlighting for `.tex` files and math zone detection (`$...$`, `$$...$$`) in markdown files.

**Architecture:** Three independent layers: (1) LaTeX language definition + backslash tokenizer for `.tex` files, (2) math zone token types in the markdown tokenizer for `.md`/`.rndm` files, (3) rendering integration in the TUI that applies syntax coloring based on file type. A new `Math` theme color and `FileType` editor field tie everything together.

**Tech Stack:** Go 1.24+, Bubble Tea v2, Lip Gloss

**Spec:** `docs/superpowers/specs/2026-03-25-latex-editor-design.md`

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `internal/themes/themes.go` | Modify | Add `Math` color field to `ThemeColors`, set in all 12 themes |
| `internal/markdown/syntax.go` | Modify | Add `MathInline`, `MathBlock` token types, `InMathBlock` to `LineContext`, math zone parsing in `TokenizeLine()` and `tokenizeInline()` |
| `internal/markdown/syntax_test.go` | Modify | Add 10 math zone test cases |
| `internal/highlight/languages.go` | Modify | Add LaTeX `LangDef` |
| `internal/highlight/builtin.go` | Modify | Add backslash command scan path for LaTeX |
| `internal/highlight/highlight_test.go` | Create | LaTeX tokenization tests (3 test cases) |
| `internal/editor/editor.go` | Modify | Add `FileType` field, set in `OpenFile()` |
| `internal/ui/app.go` | Modify | Update `renderEditor()` to apply syntax coloring from spans |

---

### Task 1: Add Math Color to Theme System

**Files:**
- Modify: `internal/themes/themes.go:10-59` (ThemeColors struct)
- Modify: `internal/themes/themes.go:69-232` (all 12 theme definitions)

- [ ] **Step 1: Add `Math` field to `ThemeColors` struct**

In `internal/themes/themes.go`, add `Math` after `CheckboxDone` in the Markdown syntax section:

```go
// Markdown syntax
H1, H2, H3, H4, H5, H6 Color
Bold                     Color
Italic                   Color
Code                     Color
CodeBg                   Color
CodeBlockBg              Color
Link                     Color
LinkURL                  Color
ListMarker               Color
Blockquote               Color
HR                       Color
Strikethrough            Color
Checkbox                 Color
CheckboxDone             Color
Math                     Color
```

- [ ] **Step 2: Set `Math` color in all 12 theme definitions**

Add `Math` to each theme. Choose a purple/violet accent that's distinct from headings and code keywords in each theme:

| Theme | Math Color | Rationale |
|-------|-----------|-----------|
| default | `#FF00FF` | Magenta (distinct from cyan H1) |
| dracula | `#BD93F9` | Purple (Dracula purple) |
| gruvbox | `#D3869B` | Pink/mauve |
| nord | `#B48EAD` | Aurora purple |
| solarized | `#6C71C4` | Violet |
| monokai | `#AE81FF` | Purple |
| catppuccin | `#CBA6F7` | Mauve |
| tokyo-night | `#BB9AF7` | Purple |
| one-dark | `#C678DD` | Purple |
| rose-pine | `#C4A7E7` | Iris |
| kanagawa | `#957FB8` | Spring violet |
| everforest | `#D699B6` | Purple |

For `defaultColors` (line 69), add after `CheckboxDone: "#00FF00"`:
```go
Math: "#FF00FF",
```

For each named theme, add `Math: "<color>"` after the `CheckboxDone` line. Example for dracula (line 97):
```go
CheckboxDone: "#50FA7B",
Math: "#BD93F9",
```

Repeat for all 12 themes with the colors from the table above.

- [ ] **Step 3: Run tests to verify no compilation errors**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go build ./...`
Expected: Clean build, no errors.

- [ ] **Step 4: Commit**

```bash
git add internal/themes/themes.go
git commit -m "feat(themes): add Math accent color to all 12 themes"
```

---

### Task 2: Add Math Zone Token Types and Tokenization

**Files:**
- Modify: `internal/markdown/syntax.go:8-33` (TokenType enum)
- Modify: `internal/markdown/syntax.go:42-45` (LineContext struct)
- Modify: `internal/markdown/syntax.go:99-150` (TokenizeLine function)
- Modify: `internal/markdown/syntax.go:152-263` (tokenizeInline function)
- Modify: `internal/markdown/syntax_test.go` (add test cases)

- [ ] **Step 1: Write failing tests for math zone tokenization**

Add to `internal/markdown/syntax_test.go`:

```go
func TestTokenizeMathInlineSimple(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine("$x^2$", ctx)
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d: %+v", len(spans), spans)
	}
	if spans[0].Token != MathInline {
		t.Errorf("expected MathInline, got %v", spans[0].Token)
	}
	if spans[0].Start != 0 || spans[0].End != 5 {
		t.Errorf("expected span [0,5), got [%d,%d)", spans[0].Start, spans[0].End)
	}
}

func TestTokenizeMathInlineWithText(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine("text $a+b$ more", ctx)
	if len(spans) != 3 {
		t.Fatalf("expected 3 spans, got %d: %+v", len(spans), spans)
	}
	if spans[0].Token != Normal {
		t.Errorf("span 0: expected Normal, got %v", spans[0].Token)
	}
	if spans[1].Token != MathInline {
		t.Errorf("span 1: expected MathInline, got %v", spans[1].Token)
	}
	if spans[2].Token != Normal {
		t.Errorf("span 2: expected Normal, got %v", spans[2].Token)
	}
}

func TestTokenizeMathBlockOpenClose(t *testing.T) {
	ctx := &LineContext{}

	// Opening $$
	spans := TokenizeLine("$$", ctx)
	if !ctx.InMathBlock {
		t.Fatal("expected InMathBlock after opening $$")
	}
	if len(spans) != 1 || spans[0].Token != MathBlock {
		t.Fatalf("expected 1 MathBlock span for opening $$, got %+v", spans)
	}

	// Content line
	spans = TokenizeLine("x = 1", ctx)
	if len(spans) != 1 || spans[0].Token != MathBlock {
		t.Fatalf("expected 1 MathBlock span for content, got %+v", spans)
	}

	// Closing $$
	spans = TokenizeLine("$$", ctx)
	if ctx.InMathBlock {
		t.Fatal("expected not InMathBlock after closing $$")
	}
	if len(spans) != 1 || spans[0].Token != MathBlock {
		t.Fatalf("expected 1 MathBlock span for closing $$, got %+v", spans)
	}
}

func TestTokenizeMathUnclosedInlineDollar(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine("price is $5", ctx)
	// No closing $ on same line → treated as literal
	for _, s := range spans {
		if s.Token == MathInline {
			t.Errorf("unclosed $ should not produce MathInline, got %+v", spans)
		}
	}
}

func TestTokenizeMathEscapedDollar(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine(`\$escaped`, ctx)
	for _, s := range spans {
		if s.Token == MathInline || s.Token == MathBlock {
			t.Errorf("escaped $ should not produce math token, got %+v", spans)
		}
	}
}

func TestTokenizeMathCodeSpanPrecedence(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine("`code $not math$`", ctx)
	for _, s := range spans {
		if s.Token == MathInline {
			t.Errorf("$ inside code span should not produce MathInline, got %+v", spans)
		}
	}
	// Should have CodeInline
	found := false
	for _, s := range spans {
		if s.Token == CodeInline {
			found = true
		}
	}
	if !found {
		t.Error("expected CodeInline span")
	}
}

func TestTokenizeMathInsideCodeBlock(t *testing.T) {
	ctx := &LineContext{InCodeBlock: true}
	spans := TokenizeLine("$x^2$", ctx)
	// Inside code block → no math tokenization
	for _, s := range spans {
		if s.Token == MathInline {
			t.Errorf("math inside code block should not produce MathInline, got %+v", spans)
		}
	}
	if len(spans) != 1 || spans[0].Token != CodeBlock {
		t.Errorf("expected CodeBlock, got %+v", spans)
	}
}

func TestTokenizeMathBlockPrecedenceOverCodeBlock(t *testing.T) {
	ctx := &LineContext{}
	// Code blocks take precedence: opening code fence, then $$ inside is code
	TokenizeLine("```", ctx) // opens code block
	if !ctx.InCodeBlock {
		t.Fatal("expected InCodeBlock")
	}
	spans := TokenizeLine("$$", ctx)
	if ctx.InMathBlock {
		t.Fatal("$$ inside code block should not open math block")
	}
	if len(spans) != 1 || spans[0].Token != CodeBlock {
		t.Errorf("expected CodeBlock, got %+v", spans)
	}
}

func TestTokenizeMathBlockMidLine(t *testing.T) {
	ctx := &LineContext{}
	// $$ mid-line is NOT a block delimiter — parsed as inline
	spans := TokenizeLine("$$x$$ and $y$", ctx)
	if ctx.InMathBlock {
		t.Fatal("mid-line $$ should not open math block")
	}
	// Should find MathInline spans for $x$ and $y$
	mathCount := 0
	for _, s := range spans {
		if s.Token == MathInline {
			mathCount++
		}
	}
	if mathCount != 2 {
		t.Errorf("expected 2 MathInline spans, got %d: %+v", mathCount, spans)
	}
}

func TestTokenizeMathBlockUnclosedEOF(t *testing.T) {
	ctx := &LineContext{}
	TokenizeLine("$$", ctx)
	if !ctx.InMathBlock {
		t.Fatal("expected InMathBlock after $$")
	}
	// Simulate EOF — InMathBlock persists
	TokenizeLine("some math content", ctx)
	if !ctx.InMathBlock {
		t.Fatal("InMathBlock should persist until closing $$")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/markdown/ -v -run "TestTokenizeMath"`
Expected: Compilation errors — `MathInline`, `MathBlock`, `InMathBlock` not defined.

- [ ] **Step 3: Add MathInline and MathBlock token types**

In `internal/markdown/syntax.go`, add after `WikiLink` (line 32):

```go
	WikiLink
	MathInline
	MathBlock
)
```

- [ ] **Step 4: Add InMathBlock to LineContext**

In `internal/markdown/syntax.go`, modify `LineContext` (line 43):

```go
type LineContext struct {
	InCodeBlock bool
	InMathBlock bool
}
```

- [ ] **Step 5: Add math block handling to TokenizeLine**

In `internal/markdown/syntax.go`, modify `TokenizeLine()`. Insert math block logic after the code fence check (line 103) and code block check (line 108), but before the empty-line check:

```go
func TokenizeLine(line string, ctx *LineContext) []Span {
	if IsCodeFence(line) {
		ctx.InCodeBlock = !ctx.InCodeBlock
		return []Span{{Start: 0, End: len(line), Token: CodeBlockMarker}}
	}

	if ctx.InCodeBlock {
		return []Span{{Start: 0, End: len(line), Token: CodeBlock}}
	}

	// Math block: $$ on its own line (optional leading whitespace)
	if isMathFence(line) {
		ctx.InMathBlock = !ctx.InMathBlock
		return []Span{{Start: 0, End: len(line), Token: MathBlock}}
	}

	if ctx.InMathBlock {
		return []Span{{Start: 0, End: len(line), Token: MathBlock}}
	}

	if len(line) == 0 {
		return nil
	}
	// ... rest unchanged
```

Add the `isMathFence` helper at the bottom of the file:

```go
// isMathFence returns true if the line is a $$ math block delimiter.
// Only matches $$ at line start (with optional leading whitespace, nothing else).
func isMathFence(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return trimmed == "$$"
}
```

- [ ] **Step 6: Add inline math parsing to tokenizeInline**

In `internal/markdown/syntax.go`, modify `tokenizeInline()`. Add `$` math detection after the inline code block (line 168), before bold+italic (line 171):

```go
		// Inline code
		if line[i] == '`' {
			if i > textStart {
				spans = append(spans, Span{Start: textStart, End: i, Token: Normal})
			}
			codeEnd := findInlineCode(line, i, end)
			spans = append(spans, Span{Start: i, End: codeEnd, Token: CodeInline})
			i = codeEnd
			textStart = i
			continue
		}

		// Inline math $...$
		if line[i] == '$' && !(i > 0 && line[i-1] == '\\') {
			mathEnd := findInlineMath(line, i+1, end)
			if mathEnd >= 0 {
				if i > textStart {
					spans = append(spans, Span{Start: textStart, End: i, Token: Normal})
				}
				spans = append(spans, Span{Start: i, End: mathEnd, Token: MathInline})
				i = mathEnd
				textStart = i
				continue
			}
		}

		// Bold + italic (*** or ___)
		// ... rest unchanged
```

Add the `findInlineMath` helper at the bottom of the file:

```go
// findInlineMath finds the closing $ for inline math starting after the opening $.
// Returns the position after the closing $ (like other find* functions), or -1 if not found.
// Rejects empty math spans (closing $ immediately after opening $, as in $$).
func findInlineMath(line string, start, end int) int {
	for i := start; i < end; i++ {
		if line[i] == '$' && i > start && !(line[i-1] == '\\') {
			return i + 1
		}
	}
	return -1
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/markdown/ -v -run "TestTokenizeMath"`
Expected: All 10 math tests PASS.

- [ ] **Step 8: Run all existing markdown tests to verify no regressions**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/markdown/ -v`
Expected: All tests PASS (existing + new).

- [ ] **Step 9: Commit**

```bash
git add internal/markdown/syntax.go internal/markdown/syntax_test.go
git commit -m "feat(markdown): add MathInline and MathBlock token types with math zone parsing"
```

---

### Task 3: Add LaTeX Language Definition and Backslash Tokenizer

**Files:**
- Modify: `internal/highlight/languages.go:35-467` (add LaTeX LangDef to allLanguages)
- Modify: `internal/highlight/builtin.go:60-174` (add backslash scan path)
- Create: `internal/highlight/highlight_test.go`

- [ ] **Step 1: Write failing tests for LaTeX highlighting**

Create `internal/highlight/highlight_test.go`:

```go
package highlight

import "testing"

func TestLatexBuiltinCommand(t *testing.T) {
	h := NewBuiltin()
	state := &State{}
	spans := h.Tokenize(`\frac{a}{b}`, "latex", state)
	if len(spans) == 0 {
		t.Fatal("expected spans, got none")
	}
	// \frac should be Builtin
	if spans[0].Kind != Builtin {
		t.Errorf("expected \\frac as Builtin, got %v", spans[0].Kind)
	}
	// Verify \frac span covers exactly \frac (5 chars)
	if spans[0].Start != 0 || spans[0].End != 5 {
		t.Errorf("expected \\frac span [0,5), got [%d,%d)", spans[0].Start, spans[0].End)
	}
}

func TestLatexComment(t *testing.T) {
	h := NewBuiltin()
	state := &State{}
	spans := h.Tokenize("% this is a comment", "latex", state)
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d: %+v", len(spans), spans)
	}
	if spans[0].Kind != Comment {
		t.Errorf("expected Comment, got %v", spans[0].Kind)
	}
}

func TestLatexKeywordAndType(t *testing.T) {
	h := NewBuiltin()
	state := &State{}
	// \begin{theorem} → \begin as Keyword, then {, theorem, }
	spans := h.Tokenize(`\begin{theorem}`, "latex", state)
	if len(spans) < 2 {
		t.Fatalf("expected at least 2 spans, got %d: %+v", len(spans), spans)
	}
	if spans[0].Kind != Keyword {
		t.Errorf("expected \\begin as Keyword, got %v", spans[0].Kind)
	}
	// Find "theorem" span — should be TypeName
	found := false
	for _, s := range spans {
		text := `\begin{theorem}`[s.Start:s.End]
		if text == "theorem" && s.Kind == TypeName {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'theorem' as TypeName, spans: %+v", spans)
	}
}

func TestLatexFindLang(t *testing.T) {
	lang := FindLang("latex")
	if lang == nil {
		t.Fatal("FindLang(\"latex\") returned nil")
	}
	if lang.Name != "latex" {
		t.Errorf("expected name 'latex', got %q", lang.Name)
	}

	// Test alias
	lang2 := FindLang("tex")
	if lang2 == nil {
		t.Fatal("FindLang(\"tex\") returned nil")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/highlight/ -v -run "TestLatex"`
Expected: FAIL — `FindLang("latex")` returns nil, backslash commands not tokenized.

- [ ] **Step 3: Add LaTeX LangDef to languages.go**

In `internal/highlight/languages.go`, add after the last entry in `allLanguages` (after the Java definition closing `}}`), before the final `}`:

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
	},
```

- [ ] **Step 4: Add backslash command scan path in builtin.go**

In `internal/highlight/builtin.go`, inside the `highlightLine()` main loop (line 60), add a LaTeX backslash scan **before** the identifier check (line 150). Insert between the annotation prefix block (ends ~line 147) and the identifiers block (line 150):

```go
		// LaTeX backslash commands: \frac, \begin, etc.
		if line[i] == '\\' && lang.Name == "latex" {
			cmdEnd := i + 1
			for cmdEnd < len(line) && isIdentStart(line[cmdEnd]) {
				cmdEnd++
			}
			if cmdEnd > i+1 {
				cmd := line[i:cmdEnd]
				if matchWord(lang.Keywords, cmd) {
					spans = append(spans, Span{Start: i, End: cmdEnd, Kind: Keyword})
				} else if matchWord(lang.Builtins, cmd) {
					spans = append(spans, Span{Start: i, End: cmdEnd, Kind: Builtin})
				} else {
					spans = append(spans, Span{Start: i, End: cmdEnd, Kind: Normal})
				}
				i = cmdEnd
				continue
			}
		}

		// Identifiers
		if isIdentStart(line[i]) {
```

- [ ] **Step 5: Run LaTeX tests to verify they pass**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/highlight/ -v -run "TestLatex"`
Expected: All 4 LaTeX tests PASS.

- [ ] **Step 6: Run all highlight tests to verify no regressions**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./internal/highlight/ -v`
Expected: All tests PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/highlight/languages.go internal/highlight/builtin.go internal/highlight/highlight_test.go
git commit -m "feat(highlight): add LaTeX language definition with backslash command tokenizer"
```

---

### Task 4: Add FileType to EditorModel

**Files:**
- Modify: `internal/editor/editor.go:66-84` (EditorModel struct)
- Modify: `internal/editor/editor.go:154-167` (OpenFile method)

- [ ] **Step 1: Add FileType field to EditorModel**

In `internal/editor/editor.go`, add `FileType` to `EditorModel` (after `File` on line 78):

```go
type EditorModel struct {
	Buf        *buffer.Buffer
	Row        int
	Col        int
	DesiredCol int
	ScrollRow  int
	ScrollCol  int
	ViewWidth  int
	ViewHeight int
	mode       Mode
	CmdBuf     string
	File       string
	FileType   string // "markdown" or "latex"
	Status     StatusMsg
	ShouldQuit bool
	Count      int
	PendingOp  rune
	CmdExec    CommandExecutor
	changed    bool
}
```

- [ ] **Step 2: Set FileType in OpenFile based on extension**

In `internal/editor/editor.go`, modify `OpenFile()`. Add file type detection after `e.File = path` (line 158). Add `"path/filepath"` to imports, then:

```go
func (e *EditorModel) OpenFile(path string) error {
	if err := e.Buf.LoadFile(path); err != nil {
		return err
	}
	e.File = path
	if filepath.Ext(path) == ".tex" {
		e.FileType = "latex"
	} else {
		e.FileType = "markdown"
	}
	e.Row = 0
	e.Col = 0
	e.ScrollRow = 0
	e.SetStatus("File opened", false)
	if e.CmdExec != nil {
		e.CmdExec.Broadcast("file_open", e)
	}
	return nil
}
```

Also set default in `LoadContent()`:

```go
func (e *EditorModel) LoadContent(name string, content string) {
	e.Buf.SetContent(content)
	e.File = name
	e.FileType = "markdown"
	e.Row = 0
	// ... rest unchanged
```

- [ ] **Step 3: Verify build succeeds**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Commit**

```bash
git add internal/editor/editor.go
git commit -m "feat(editor): add FileType field with .tex detection in OpenFile"
```

---

### Task 5: Wire Syntax Coloring into TUI Rendering

**Files:**
- Modify: `internal/ui/app.go:611-682` (renderEditor function)

This is the final integration task. The `renderEditor()` function currently renders lines as plain text. We add span-based coloring using the markdown tokenizer (for `.md`/`.rndm`) or the highlight tokenizer (for `.tex`).

- [ ] **Step 1: Add imports to app.go**

Add these imports to `internal/ui/app.go` (check existing imports first — some may already be present):

```go
"github.com/EME130/lazymd/internal/highlight"
"github.com/EME130/lazymd/internal/markdown"
"github.com/EME130/lazymd/internal/themes"
"github.com/charmbracelet/lipgloss"
```

- [ ] **Step 2: Add a colorLine helper method to AppModel**

Add this method somewhere in `internal/ui/app.go` (e.g., after `renderEditor()`):

```go
// colorLine applies syntax coloring to a line based on file type.
// For markdown files, uses the markdown tokenizer. For LaTeX, uses the highlight tokenizer.
// Returns the styled string.
func (m *AppModel) colorLine(line string, contentW int) string {
	if m.editor.FileType == "latex" {
		return m.colorLineLatex(line, contentW)
	}
	return m.colorLineMarkdown(line, contentW)
}

func (m *AppModel) colorLineMarkdown(line string, contentW int) string {
	spans := markdown.TokenizeLine(line, &m.mdLineCtx)
	if len(spans) == 0 {
		if len(line) > contentW {
			return line[:contentW]
		}
		return line
	}
	tc := themes.CurrentColors()
	var result string
	for _, s := range spans {
		text := line[s.Start:s.End]
		if len(text) == 0 {
			continue
		}
		color := mdTokenColor(s.Token, tc)
		if color != "" {
			result += lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
		} else {
			result += text
		}
	}
	return result
}

func (m *AppModel) colorLineLatex(line string, contentW int) string {
	spans := m.hlBuiltin.Tokenize(line, "latex", &m.hlState)
	if len(spans) == 0 {
		if len(line) > contentW {
			return line[:contentW]
		}
		return line
	}
	tc := themes.CurrentColors()
	var result string
	for _, s := range spans {
		text := line[s.Start:s.End]
		if len(text) == 0 {
			continue
		}
		color := tc.SyntaxColor(hlKindName(s.Kind))
		if color != "" {
			result += lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
		} else {
			result += text
		}
	}
	return result
}

func mdTokenColor(tok markdown.TokenType, tc *themes.ThemeColors) string {
	switch tok {
	case markdown.H1:
		return tc.H1
	case markdown.H2:
		return tc.H2
	case markdown.H3:
		return tc.H3
	case markdown.H4:
		return tc.H4
	case markdown.H5:
		return tc.H5
	case markdown.H6:
		return tc.H6
	case markdown.Bold:
		return tc.Bold
	case markdown.Italic:
		return tc.Italic
	case markdown.BoldItalic:
		return tc.Bold
	case markdown.CodeInline, markdown.CodeBlock, markdown.CodeBlockMarker:
		return tc.Code
	case markdown.LinkText:
		return tc.Link
	case markdown.LinkURL:
		return tc.LinkURL
	case markdown.ListBullet, markdown.ListNumber:
		return tc.ListMarker
	case markdown.Blockquote:
		return tc.Blockquote
	case markdown.HR:
		return tc.HR
	case markdown.Strikethrough:
		return tc.Strikethrough
	case markdown.TaskCheckbox:
		return tc.Checkbox
	case markdown.WikiLink:
		return tc.Link
	case markdown.MathInline, markdown.MathBlock:
		return tc.Math
	default:
		return ""
	}
}

func hlKindName(k highlight.TokenKind) string {
	switch k {
	case highlight.Keyword:
		return "keyword"
	case highlight.TypeName:
		return "type_name"
	case highlight.Builtin:
		return "builtin"
	case highlight.String:
		return "string"
	case highlight.Number:
		return "number"
	case highlight.Comment:
		return "comment"
	case highlight.Operator:
		return "operator"
	case highlight.Punctuation:
		return "punctuation"
	case highlight.Annotation:
		return "annotation"
	default:
		return ""
	}
}
```

- [ ] **Step 3: Add state fields to AppModel**

Find the `AppModel` struct in `app.go` and add these fields:

```go
	mdLineCtx markdown.LineContext
	hlBuiltin *highlight.BuiltinHighlighter
	hlState   highlight.State
```

Initialize `hlBuiltin` in the constructor (wherever `AppModel` is created — look for `NewApp` or similar):

```go
hlBuiltin: highlight.NewBuiltin(),
```

- [ ] **Step 4: Modify renderEditor to use colorLine**

In `renderEditor()`, replace the plain-text line rendering with colored rendering. The key change is at line 644 where `line := ed.Buf.Line(bufRow)` is used.

**Important:** Reset `mdLineCtx` and `hlState` at the start of each render pass (since we re-render from `ScrollRow` each time):

At the top of `renderEditor()`, after the `availH` check:

```go
	// Reset tokenizer state for this render pass
	m.mdLineCtx = markdown.LineContext{}
	m.hlState = highlight.State{}
	// Pre-process lines before ScrollRow to build up correct state
	for r := 0; r < ed.ScrollRow && r < ed.Buf.LineCount(); r++ {
		if m.editor.FileType == "latex" {
			m.hlBuiltin.Tokenize(ed.Buf.Line(r), "latex", &m.hlState)
		} else {
			markdown.TokenizeLine(ed.Buf.Line(r), &m.mdLineCtx)
		}
	}
```

Then, replace the line content rendering. Change this section (current lines ~644-668):

```go
		// Line content with syntax coloring
		rawLine := ed.Buf.Line(bufRow)
		coloredLine := m.colorLine(rawLine, contentW)

		// Apply horizontal scroll to the raw line for cursor calculation
		line := rawLine
		if ed.ScrollCol > 0 && ed.ScrollCol < len(line) {
			line = line[ed.ScrollCol:]
		} else if ed.ScrollCol >= len(line) {
			line = ""
		}

		if len(line) > contentW {
			line = line[:contentW]
		}

		// For display, apply scroll to colored output
		displayLine := coloredLine
		// TODO: horizontal scroll on colored output is complex (ANSI codes);
		// for now, fall back to plain text when scrolled
		if ed.ScrollCol > 0 {
			displayLine = line
		}

		// Render cursor on current line
		if bufRow == ed.Row {
			cursorCol := ed.Col - ed.ScrollCol
			if cursorCol >= 0 && cursorCol < len(line) {
				before := line[:cursorCol]
				cursor := m.styles.Cursor.Render(string(line[cursorCol]))
				after := ""
				if cursorCol+1 < len(line) {
					after = line[cursorCol+1:]
				}
				displayLine = before + cursor + after
			} else if cursorCol == len(line) {
				displayLine = line + m.styles.Cursor.Render(" ")
			}
		}
```

And update the pad calculation to use `displayLine` instead of `line`:

```go
		displayW := lipgloss.Width(displayLine)
		pad := contentW - displayW
		if pad < 0 {
			pad = 0
		}

		lines = append(lines, gutter+" "+displayLine+strings.Repeat(" ", pad))
```

- [ ] **Step 5: Verify build succeeds**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go build ./...`
Expected: Clean build.

- [ ] **Step 6: Run all tests**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go test ./...`
Expected: All tests PASS.

- [ ] **Step 7: Manual smoke test**

Run: `export PATH="/usr/local/go/bin:$PATH" && cd /Users/eme130/projects/lazymd && go run ./cmd/lm/ README.md`

Verify:
- Headings appear colored (not plain white)
- `$math$` zones appear in the Math accent color
- Code blocks still render correctly
- Cursor still works

If you have a `.tex` file, open it and verify:
- `\frac`, `\int` etc. appear as Builtin color
- `\begin`, `\end` appear as Keyword color
- `% comments` appear as Comment color
- `theorem`, `equation` etc. appear as Type color

- [ ] **Step 8: Commit**

```bash
git add internal/ui/app.go
git commit -m "feat(ui): add syntax coloring to editor for markdown and LaTeX files"
```

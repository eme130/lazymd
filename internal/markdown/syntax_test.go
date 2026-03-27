package markdown

import "testing"

func TestHeaderDetection(t *testing.T) {
	tests := []struct {
		input string
		level int
	}{
		{"# Hello", 1},
		{"### Test", 3},
		{"####### Too many", 0},
		{"Not a header", 0},
	}
	for _, tt := range tests {
		got := parseHeader(tt.input)
		if got != tt.level {
			t.Errorf("parseHeader(%q) = %d, want %d", tt.input, got, tt.level)
		}
	}
}

func TestHorizontalRule(t *testing.T) {
	trueTests := []string{"---", "***", "___", "- - -"}
	for _, s := range trueTests {
		if !IsHorizontalRule(s) {
			t.Errorf("IsHorizontalRule(%q) = false, want true", s)
		}
	}
	falseTests := []string{"--", "abc"}
	for _, s := range falseTests {
		if IsHorizontalRule(s) {
			t.Errorf("IsHorizontalRule(%q) = true, want false", s)
		}
	}
}

func TestCodeFence(t *testing.T) {
	trueTests := []string{"```", "```zig", "~~~"}
	for _, s := range trueTests {
		if !IsCodeFence(s) {
			t.Errorf("IsCodeFence(%q) = false, want true", s)
		}
	}
	if IsCodeFence("``") {
		t.Error("IsCodeFence(``) = true, want false")
	}
}

func TestParseCodeFence(t *testing.T) {
	tests := []struct {
		input    string
		isFence  bool
		language string
	}{
		{"```zig", true, "zig"},
		{"```python", true, "python"},
		{"```", true, ""},
		{"~~~js", true, "js"},
		{"not a fence", false, ""},
	}
	for _, tt := range tests {
		got := ParseCodeFence(tt.input)
		if got.IsFence != tt.isFence {
			t.Errorf("ParseCodeFence(%q).IsFence = %v, want %v", tt.input, got.IsFence, tt.isFence)
		}
		if got.Language != tt.language {
			t.Errorf("ParseCodeFence(%q).Language = %q, want %q", tt.input, got.Language, tt.language)
		}
	}
}

func TestTokenizeHeaderLine(t *testing.T) {
	ctx := &LineContext{}
	spans := TokenizeLine("## Hello world", ctx)

	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Token != H2 {
		t.Fatalf("expected H2, got %v", spans[0].Token)
	}
}

func TestTokenizeCodeBlock(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("```zig", ctx)
	if !ctx.InCodeBlock {
		t.Fatal("expected in_code_block after opening fence")
	}
	if spans[0].Token != CodeBlockMarker {
		t.Fatalf("expected CodeBlockMarker, got %v", spans[0].Token)
	}

	spans = TokenizeLine("const x = 5;", ctx)
	if spans[0].Token != CodeBlock {
		t.Fatalf("expected CodeBlock, got %v", spans[0].Token)
	}

	spans = TokenizeLine("```", ctx)
	if ctx.InCodeBlock {
		t.Fatal("expected not in_code_block after closing fence")
	}
}

func TestTokenizeInlineFormatting(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("hello **bold** world", ctx)
	if len(spans) != 3 {
		t.Fatalf("expected 3 spans, got %d", len(spans))
	}
	if spans[0].Token != Normal {
		t.Errorf("span 0: expected Normal, got %v", spans[0].Token)
	}
	if spans[1].Token != Bold {
		t.Errorf("span 1: expected Bold, got %v", spans[1].Token)
	}
	if spans[2].Token != Normal {
		t.Errorf("span 2: expected Normal, got %v", spans[2].Token)
	}
}

func TestTokenizeWikiLink(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("see [[my note]] here", ctx)
	if len(spans) != 3 {
		t.Fatalf("expected 3 spans, got %d", len(spans))
	}
	if spans[1].Token != WikiLink {
		t.Errorf("span 1: expected WikiLink, got %v", spans[1].Token)
	}
}

func TestTokenizeLink(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("click [here](https://example.com) now", ctx)
	if len(spans) != 4 {
		t.Fatalf("expected 4 spans, got %d", len(spans))
	}
	if spans[1].Token != LinkText {
		t.Errorf("span 1: expected LinkText, got %v", spans[1].Token)
	}
	if spans[2].Token != LinkURL {
		t.Errorf("span 2: expected LinkURL, got %v", spans[2].Token)
	}
}

func TestTokenizeList(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("- item one", ctx)
	if len(spans) < 1 {
		t.Fatal("expected at least 1 span")
	}
	if spans[0].Token != ListBullet {
		t.Errorf("span 0: expected ListBullet, got %v", spans[0].Token)
	}
}

func TestTokenizeNumberedList(t *testing.T) {
	ctx := &LineContext{}

	spans := TokenizeLine("1. first item", ctx)
	if len(spans) < 1 {
		t.Fatal("expected at least 1 span")
	}
	if spans[0].Token != ListNumber {
		t.Errorf("span 0: expected ListNumber, got %v", spans[0].Token)
	}
}

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
	spans := TokenizeLine("$$", ctx)
	if !ctx.InMathBlock {
		t.Fatal("expected InMathBlock after opening $$")
	}
	if len(spans) != 1 || spans[0].Token != MathBlock {
		t.Fatalf("expected 1 MathBlock span for opening $$, got %+v", spans)
	}
	spans = TokenizeLine("x = 1", ctx)
	if len(spans) != 1 || spans[0].Token != MathBlock {
		t.Fatalf("expected 1 MathBlock span for content, got %+v", spans)
	}
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
	TokenizeLine("```", ctx)
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
	spans := TokenizeLine("$$x$$ and $y$", ctx)
	if ctx.InMathBlock {
		t.Fatal("mid-line $$ should not open math block")
	}
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
	TokenizeLine("some math content", ctx)
	if !ctx.InMathBlock {
		t.Fatal("InMathBlock should persist until closing $$")
	}
}

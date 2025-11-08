package service

import "testing"

func TestStripRichTextTipTap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "heading content",
			input: `{"type":"doc","content":[{"type":"heading","attrs":{"level":2},"content":[{"type":"text","text":"12321"}]},{"type":"paragraph"}]}`,
			want:  "12321",
		},
		{
			name:  "hard break content",
			input: `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hello"},{"type":"hardBreak"},{"type":"text","text":"World"}]}]}`,
			want:  "Hello\nWorld",
		},
		{
			name:  "mention label fallback",
			input: `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"mention","attrs":{"label":"@admin"}}]}]}`,
			want:  "@admin",
		},
		{
			name: "multi docs",
			input: `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"A"}]}]}` + "\n" +
				`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"B"}]}]}`,
			want: "A\nB",
		},
		{
			name:  "empty doc",
			input: `{"type":"doc","content":[{"type":"paragraph"}]}`,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripRichText(tt.input); got != tt.want {
				t.Fatalf("stripRichText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripRichTextFallbackHTML(t *testing.T) {
	input := "<p>富文本</p>"
	want := "富文本"

	if got := stripRichText(input); got != want {
		t.Fatalf("stripRichText() = %q, want %q", got, want)
	}
}

func TestStripRichTextImageTipTap(t *testing.T) {
	attachmentBaseURLOverride = "https://seal.test"
	t.Cleanup(func() {
		attachmentBaseURLOverride = ""
	})
	input := `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"图像"},{"type":"image","attrs":{"src":"id:abc123"}}]}]}`
	want := "图像[CQ:image,file=image,url=https://seal.test/api/v1/attachment/abc123]"
	if got := stripRichText(input); got != want {
		t.Fatalf("stripRichText() = %q, want %q", got, want)
	}
}

func TestStripRichTextImageHTML(t *testing.T) {
	attachmentBaseURLOverride = "https://seal.test"
	t.Cleanup(func() {
		attachmentBaseURLOverride = ""
	})
	input := `<p>示例<img src="/api/v1/attachment/xyz" /></p>`
	want := "示例[CQ:image,file=image,url=https://seal.test/api/v1/attachment/xyz]"
	if got := stripRichText(input); got != want {
		t.Fatalf("stripRichText() = %q, want %q", got, want)
	}
}

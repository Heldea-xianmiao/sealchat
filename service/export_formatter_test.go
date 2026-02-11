package service

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeDomainToURLIPv6(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ipv6 with port",
			input: "[2001:db8::1]:3212",
			want:  "https://[2001:db8::1]:3212",
		},
		{
			name:  "ipv6 loopback without port",
			input: "::1",
			want:  "http://[::1]",
		},
		{
			name:  "ipv4 loopback",
			input: "127.0.0.1:8080",
			want:  "http://127.0.0.1:8080",
		},
		{
			name:  "ipv6 link-local",
			input: "fe80::1",
			want:  "http://[fe80::1]",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeDomainToURL(tt.input); got != tt.want {
				t.Fatalf("normalizeDomainToURL(%q) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestEnhancePlainContentForHTMLExportQuickFormat(t *testing.T) {
	input := "**粗体** *斜体* `代码` [链接](https://example.com)"
	result := enhancePlainContentForHTMLExport(input)

	expects := []string{
		"<strong>粗体</strong>",
		"<em>斜体</em>",
		"<code>代码</code>",
		`<a href="https://example.com" class="text-blue-500" target="_blank" rel="noopener noreferrer">链接</a>`,
	}

	for _, expected := range expects {
		if !strings.Contains(result, expected) {
			t.Fatalf("expect html contains %q, got %q", expected, result)
		}
	}
}

func TestEnhancePlainContentForHTMLExportInvalidLink(t *testing.T) {
	input := "[危险](javascript:alert(1))"
	result := enhancePlainContentForHTMLExport(input)
	if strings.Contains(result, "<a ") {
		t.Fatalf("invalid protocol should not become link, got %q", result)
	}
}

func TestBuildBBCodeTextLineFromQuickFormat(t *testing.T) {
	payload := &ExportPayload{WithoutTimestamp: true}
	msg := &ExportMessage{
		SenderName:  "测试",
		SenderColor: "#123abc",
		CreatedAt:   time.Unix(1700000000, 0),
		Content:     "**粗体** *斜体* `代码` [链接](https://example.com)",
	}

	line := buildBBCodeTextLine(payload, msg)
	expects := []string{"[b]粗体[/b]", "[i]斜体[/i]", "[code]代码[/code]", "[url=https://example.com]链接[/url]"}
	for _, expected := range expects {
		if !strings.Contains(line, expected) {
			t.Fatalf("expect bbcode contains %q, got %q", expected, line)
		}
	}
}

func TestEnhancePlainContentForHTMLExportDoesNotRenderCodeFence(t *testing.T) {
	input := "```\nconst a = 1\n```"
	result := enhancePlainContentForHTMLExport(input)
	if strings.Contains(result, "<pre><code>") {
		t.Fatalf("code fence should not become code block, got %q", result)
	}
	if !strings.Contains(result, "```") {
		t.Fatalf("code fence should remain literal text, got %q", result)
	}
}

func TestBuildBBCodeTextLineDoesNotRenderCodeFence(t *testing.T) {
	payload := &ExportPayload{WithoutTimestamp: true}
	msg := &ExportMessage{
		SenderName:  "测试",
		SenderColor: "#123abc",
		CreatedAt:   time.Unix(1700000000, 0),
		Content:     "```hello```",
	}

	line := buildBBCodeTextLine(payload, msg)
	if strings.Contains(line, "[code]") {
		t.Fatalf("code fence should not convert to [code], got %q", line)
	}
	if !strings.Contains(line, "```hello```") {
		t.Fatalf("code fence should remain literal text, got %q", line)
	}
}

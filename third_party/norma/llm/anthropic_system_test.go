package llm

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildAnthropicSystemCacheBoundary(t *testing.T) {
	sys := []string{"a", "b", "c"}

	cases := []struct {
		name       string
		boundary   int
		cache      bool
		wantBlocks int
		wantCached []bool // per-block: has cache_control
	}{
		{"no boundary (0) → single uncached", 0, true, 1, []bool{false}},
		{"negative → single uncached", -1, true, 1, []bool{false}},
		{"mid split → cached prefix + dynamic tail", 2, true, 2, []bool{true, false}},
		{"boundary == len → whole thing cached", 3, true, 1, []bool{true}},
		{"boundary > len → whole thing cached", 5, true, 1, []bool{true}},
		{"cache off → boundary ignored, all uncached", 2, false, 1, []bool{false}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blocks := buildAnthropicSystem(CompletionRequest{System: sys, DynamicBoundary: c.boundary}, c.cache)
			if len(blocks) != c.wantBlocks {
				t.Fatalf("blocks=%d want %d", len(blocks), c.wantBlocks)
			}
			for i, wantCached := range c.wantCached {
				gotCached := blocks[i].CacheControl != nil
				if gotCached != wantCached {
					t.Fatalf("block %d cached=%v want %v", i, gotCached, wantCached)
				}
			}
		})
	}
}

// TestBuildAnthropicToolCacheBreakpoint covers the PromptCache tool-list
// breakpoint: on → the LAST tool carries cache_control (whole tools block cached);
// off → no tool carries one. Builds a body through the provider so the wire JSON
// is exactly what a real request would send.
func TestBuildAnthropicToolCacheBreakpoint(t *testing.T) {
	tools := []ToolSchema{
		{Name: "tool_a", Description: "a", InputSchema: map[string]any{"type": "object"}},
		{Name: "tool_b", Description: "b", InputSchema: map[string]any{"type": "object"}},
	}
	req := CompletionRequest{System: []string{"sys"}, DynamicBoundary: 1, Tools: tools, MaxTokens: 64}

	on := anthropicProvider{cfg: Config{PromptCache: true, Model: "m"}}
	body, err := on.buildBody(req, false)
	if err != nil {
		t.Fatalf("buildBody: %v", err)
	}
	var wire struct {
		System []anthropicSystemBlock `json:"system"`
		Tools  []struct {
			Name         string         `json:"name"`
			CacheControl map[string]any `json:"cache_control"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(wire.Tools) != 2 {
		t.Fatalf("tools=%d want 2", len(wire.Tools))
	}
	if wire.Tools[0].CacheControl != nil {
		t.Fatalf("non-last tool must not carry cache_control: %v", wire.Tools[0].CacheControl)
	}
	if wire.Tools[1].CacheControl["type"] != "ephemeral" {
		t.Fatalf("last tool missing ephemeral cache_control: %v", wire.Tools[1].CacheControl)
	}
	if len(wire.System) != 1 || wire.System[0].CacheControl == nil {
		t.Fatalf("system breakpoint expected alongside tool breakpoint: %v", wire.System)
	}

	off := anthropicProvider{cfg: Config{PromptCache: false, Model: "m"}}
	body, err = off.buildBody(req, false)
	if err != nil {
		t.Fatalf("buildBody(off): %v", err)
	}
	// Fresh struct: json.Unmarshal reuses slice element structs, so unmarshaling
	// into the same wire value would leak the previous body's cache_control.
	var offWire struct {
		System []anthropicSystemBlock `json:"system"`
		Tools  []struct {
			Name         string         `json:"name"`
			CacheControl map[string]any `json:"cache_control"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(body, &offWire); err != nil {
		t.Fatalf("decode(off): %v", err)
	}
	if len(offWire.System) != 1 || offWire.System[0].CacheControl != nil {
		t.Fatalf("cache off: system breakpoint must be absent: %v", offWire.System)
	}
	for i, tw := range offWire.Tools {
		if tw.CacheControl != nil {
			t.Fatalf("cache off: tool %d must not carry cache_control: %v", i, tw.CacheControl)
		}
	}
	if strings.Contains(string(body), "cache_control") {
		t.Fatalf("cache off: wire body must not mention cache_control:\n%s", body)
	}
}

// TestBuildAnthropicMessageCacheBreakpoint covers the conversation-history
// breakpoint: a multi-turn request caches its whole prefix via a breakpoint on
// the last block of the last message; single-message (fresh one-shot) requests
// skip it so they never pay the cache-write premium on content no one reuses.
// The touched message must be copied — session state must stay cache-free.
func TestBuildAnthropicMessageCacheBreakpoint(t *testing.T) {
	tools := []ToolSchema{{Name: "t", Description: "d", InputSchema: map[string]any{"type": "object"}}}
	multi := []Message{
		UserText("first input"),
		{Role: RoleAssistant, Content: []ContentBlock{TextBlock("reply")}},
		{Role: RoleUser, Content: []ContentBlock{ToolResultText("toolu_1", "output", false)}}, // last = user, ends in tool_result
	}
	req := CompletionRequest{System: []string{"sys"}, DynamicBoundary: 1, Messages: multi, Tools: tools, MaxTokens: 64}

	on := anthropicProvider{cfg: Config{PromptCache: true, Model: "m"}}
	body, err := on.buildBody(req, false)
	if err != nil {
		t.Fatalf("buildBody: %v", err)
	}
	var wire struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content []ContentBlock  `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(wire.Messages) != 3 {
		t.Fatalf("messages=%d want 3", len(wire.Messages))
	}
	last := wire.Messages[2].Content[len(wire.Messages[2].Content)-1]
	if last.Type != BlockToolResult || last.CacheControl == nil || last.CacheControl["type"] != "ephemeral" {
		t.Fatalf("last block must carry ephemeral cache_control: %+v", last)
	}
	if wire.Messages[0].Content[0].CacheControl != nil || wire.Messages[1].Content[0].CacheControl != nil {
		t.Fatal("only the last block may carry the breakpoint")
	}
	// Session state untouched: the request-local copy never leaks back.
	for _, m := range multi {
		for _, b := range m.Content {
			if b.CacheControl != nil {
				t.Fatal("buildBody mutated the caller's message blocks")
			}
		}
	}

	// Single-message request → no message breakpoint (fresh one-shot session).
	single := anthropicProvider{cfg: Config{PromptCache: true, Model: "m"}}
	body, err = single.buildBody(CompletionRequest{System: []string{"sys"}, DynamicBoundary: 1, Messages: []Message{UserText("only")}, Tools: tools, MaxTokens: 64}, false)
	if err != nil {
		t.Fatalf("buildBody(single): %v", err)
	}
	if strings.Count(string(body), "cache_control") > 2 {
		t.Fatalf("single-message request must not add a message breakpoint:\n%s", body)
	}

	// Last message assistant → no message breakpoint.
	asstLast := anthropicProvider{cfg: Config{PromptCache: true, Model: "m"}}
	body, err = asstLast.buildBody(CompletionRequest{System: []string{"sys"}, DynamicBoundary: 1, Messages: []Message{UserText("q"), {Role: RoleAssistant, Content: []ContentBlock{TextBlock("a")}}}, Tools: tools, MaxTokens: 64}, false)
	if err != nil {
		t.Fatalf("buildBody(asst): %v", err)
	}
	if strings.Count(string(body), "cache_control") > 2 {
		t.Fatalf("assistant-final request must not add a message breakpoint:\n%s", body)
	}
}

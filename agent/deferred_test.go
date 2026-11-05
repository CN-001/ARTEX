package agent

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/norma/llm"
	actool "github.com/Autumn-27/norma/tool"
)

func TestDeferredSystem_NoGlobal(t *testing.T) {
	// No global MCP names → single-segment system, boundary >= len → the WHOLE
	// (task/session-stable) prompt is cacheable. Every follow-up turn/round then
	// reuses the cached prefix instead of re-paying input tokens.
	sys, boundary := deferredSystem("SYS", DeferredInfo{})
	if len(sys) != 1 || sys[0] != "SYS" || boundary != 1 {
		t.Fatalf("expected [SYS],1 — got %v,%d", sys, boundary)
	}
}

func TestDeferredSystemSplit(t *testing.T) {
	// Worker path: stable prefix + dynamic tail, boundary exactly at the seam.
	stable, dynamic := "ROLE", "INTENT"
	sys, boundary := deferredSystemSplit(stable, dynamic, DeferredInfo{})
	if len(sys) != 2 || sys[0] != "ROLE" || sys[1] != "INTENT" {
		t.Fatalf("expected [ROLE INTENT], got %v", sys)
	}
	if boundary != 1 {
		t.Fatalf("boundary=%d want 1 (only the stable prefix cached)", boundary)
	}
	// FindingGuidance joins the STABLE prefix; the deferred block joins it too —
	// both are stable per agent role, so they stay cacheable.
	def := DeferredInfo{
		FindingGuidance: "GUIDE",
		Deferred:       []string{"mcp__browser__navigate"},
		GlobalNames:    []string{"mcp__browser__navigate"},
	}
	sys, boundary = deferredSystemSplit(stable, dynamic, def)
	if boundary != 1 || len(sys) != 2 {
		t.Fatalf("split with def: sys=%v boundary=%d", sys, boundary)
	}
	if !strings.HasPrefix(sys[0], "ROLE") || !strings.Contains(sys[0], "GUIDE") ||
		!strings.Contains(sys[0], "<available-deferred-tools>") {
		t.Fatalf("stable prefix must carry role+guidance+deferred block:\n%s", sys[0])
	}
	if sys[1] != "INTENT" {
		t.Fatalf("dynamic tail changed: %q", sys[1])
	}
}

func TestDeferredSystem_WithGlobal(t *testing.T) {
	def := DeferredInfo{
		Deferred:    []string{"mcp__browser__navigate", "mcp__browser__click"},
		GlobalNames: []string{"mcp__browser__navigate", "mcp__browser__click"},
	}
	sys, boundary := deferredSystem("SYS", def)
	if len(sys) != 2 || sys[0] != "SYS" {
		t.Fatalf("expected [SYS, block], got %v", sys)
	}
	if !strings.Contains(sys[1], "<available-deferred-tools>") ||
		!strings.Contains(sys[1], "mcp__browser__navigate") {
		t.Fatalf("block missing names:\n%s", sys[1])
	}
	if boundary != len(sys) {
		t.Fatalf("boundary=%d want %d (whole prompt cached)", boundary, len(sys))
	}
}

func TestDeferredSystem_GatedNotInBlock(t *testing.T) {
	// A skill-gated server's tools are deferred but NOT in the global block.
	def := DeferredInfo{
		Deferred:    []string{"mcp__browser__navigate", "mcp__secret__do"},
		GlobalNames: []string{"mcp__browser__navigate"}, // secret gated → excluded
	}
	sys, _ := deferredSystem("SYS", def)
	if strings.Contains(sys[1], "mcp__secret__do") {
		t.Fatalf("gated tool must not appear in global block:\n%s", sys[1])
	}
	if !strings.Contains(sys[1], "mcp__browser__navigate") {
		t.Fatal("global tool should appear in block")
	}
}

func TestSeedUnlockFromHistory(t *testing.T) {
	skillCall := func(name string) llm.ContentBlock {
		return llm.ContentBlock{Type: llm.BlockToolUse, Name: "Skill", Input: json.RawMessage(`{"name":"` + name + `"}`)}
	}
	msgs := []llm.Message{
		{Role: llm.RoleAssistant, Content: []llm.ContentBlock{skillCall("browsing")}},
		{Role: llm.RoleAssistant, Content: []llm.ContentBlock{
			{Type: llm.BlockToolUse, Name: "Bash", Input: json.RawMessage(`{}`)}, // ignored
			skillCall("recon"),
		}},
	}
	var got []string
	seedUnlockFromHistory(msgs, func(name string) { got = append(got, name) })
	if strings.Join(got, ",") != "browsing,recon" {
		t.Fatalf("unlocked=%v want [browsing recon]", got)
	}
	seedUnlockFromHistory(msgs, nil) // nil → no-op, no panic
}

// TestUnlockGatingFlow mirrors what OnInvoke / seedUnlockFromHistory do: a gated
// tool starts locked and becomes callable only after its skill unlocks it.
func TestUnlockGatingFlow(t *testing.T) {
	serverTools := map[string][]string{"secret": {"mcp__secret__do"}}
	unlock := actool.NewUnlockSet("mcp__browser__navigate") // global only
	unlockSkill := func(name string) {
		if name == "unlock-secret" {
			unlock.Add(serverTools["secret"]...)
		}
	}
	if unlock.Has("mcp__secret__do") {
		t.Fatal("gated tool should start locked")
	}
	unlockSkill("unlock-secret")
	if !unlock.Has("mcp__secret__do") {
		t.Fatal("gated tool should be unlocked after skill load")
	}
}

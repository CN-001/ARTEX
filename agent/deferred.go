package agent

import (
	"encoding/json"

	"github.com/Autumn-27/norma/llm"
	actool "github.com/Autumn-27/norma/tool"
)

// deferredSystem builds an agent's system-prompt segments and cache boundary from
// its DeferredInfo. When globally-available MCP tools are present, their names + a
// "prefer core tools" instruction render into a <available-deferred-tools> block
// placed as the LAST system-prompt segment.
//
// The DynamicBoundary is ALWAYS set so the whole system prompt is cached
// (Anthropic cache_control on the static prefix; DeepSeek/OpenAI auto-cache hit
// the same stable prefix). The callers assemble it from session/task-stable parts
// (role template + constraints + work-dir spec + FindingGuidance), so caching it
// makes every follow-up turn/round reuse the prefix instead of paying
// input-token price again. Skill-gated MCP names are NOT in this block; they
// surface when their skill loads.
func deferredSystem(sysText string, def DeferredInfo) (system []string, boundary int) {
	sysText += def.FindingGuidance
	block := actool.RenderDeferredToolsBlock(def.GlobalNames)
	if block == "" {
		// Single segment: b >= len → whole prompt cached (SDK guard).
		return []string{sysText}, 1
	}
	system = []string{sysText, block}
	boundary = len(system) // b >= len → whole system prompt cached (SDK guard)
	return system, boundary
}

// deferredSystemSplit builds a TWO-segment system prompt with the cache boundary
// between them: the first segment (role body + FindingGuidance + deferred block)
// is session/task-stable and cached; the second carries per-run content that
// must not be cached because it changes per intent/session (the worker's intent
// text, target assets, run-dir spec and launch line). Boundary=1 marks exactly
// the stable prefix, so concurrent workers on the same task — and resume runs of
// the same intent — all reuse one cached prefix while the per-intent tail stays
// uncached. The deferred block joins the cached prefix (it is stable too).
func deferredSystemSplit(prefix, dynamic string, def DeferredInfo) (system []string, boundary int) {
	prefix += def.FindingGuidance
	if block := actool.RenderDeferredToolsBlock(def.GlobalNames); block != "" {
		prefix += block
	}
	return []string{prefix, dynamic}, 1
}

// seedUnlockFromHistory replays prior Skill() invocations in the conversation so
// their skill-gated MCPs are re-unlocked on a resumed session (design doc C2). The
// main agent builds a fresh session each turn; its in-memory unlock set would
// otherwise reset, leaving the model able to see a skill-revealed tool name yet
// unable to call it. No-op when unlockSkill is nil (no deferred tools).
func seedUnlockFromHistory(msgs []llm.Message, unlockSkill func(string)) {
	if unlockSkill == nil {
		return
	}
	for _, m := range msgs {
		for _, b := range m.ToolUses() {
			if b.Name != "Skill" {
				continue
			}
			var in struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(b.Input, &in) == nil && in.Name != "" {
				unlockSkill(in.Name)
			}
		}
	}
}

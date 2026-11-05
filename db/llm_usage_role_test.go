package db

import (
	"fmt"
	"testing"
	"time"
)

// UsageByWorker feeds the dashboard's per-role cache panel, so it must group the
// always-on llm_usage ledger by agent lane and order by spend. The ledger is
// shared with every other test in this database, so the assertions are DELTAS
// against the pre-insert state rather than absolute totals. Skips without a DB.
func TestUsageByWorkerGroupsByLane(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v) — skipping", err)
	}
	defer d.Close()
	// The app creates the metering table at startup; a fresh test database has
	// not run that yet.
	if err := d.EnsureLLMUsageTable(); err != nil {
		t.Fatalf("EnsureLLMUsageTable: %v", err)
	}

	byRole := func() map[string]RoleUsage {
		t.Helper()
		rows, err := d.UsageByWorker()
		if err != nil {
			t.Fatalf("UsageByWorker: %v", err)
		}
		// Ordering contract: heaviest input+output first (descending).
		prev := -1
		out := make(map[string]RoleUsage, len(rows))
		for _, r := range rows {
			spend := r.InputTokens + r.OutputTokens
			if prev >= 0 && spend > prev {
				t.Fatalf("rows not ordered by spend desc: %v then %v", prev, spend)
			}
			prev = spend
			out[r.Role] = r
		}
		return out
	}

	before := byRole()

	// Marker model keeps cleanup exact on a shared ledger.
	marker := fmt.Sprintf("role-usage-test-%d", time.Now().UnixNano())
	seeded := []LLMUsage{
		{Worker: "planner", Model: marker, InputTokens: 100, OutputTokens: 5, CacheRead: 80},
		{Worker: "planner", Model: marker, InputTokens: 100, OutputTokens: 5, CacheRead: 80},
		{Worker: "worker", Model: marker, InputTokens: 50, OutputTokens: 3, CacheRead: 10},
	}
	for i := range seeded {
		if err := d.InsertLLMUsage(&seeded[i]); err != nil {
			t.Fatalf("InsertLLMUsage: %v", err)
		}
	}
	t.Cleanup(func() { _, _ = d.Exec(`DELETE FROM llm_usage WHERE model = $1`, marker) })

	after := byRole()

	planner, plannerBefore := after["planner"], before["planner"]
	if got := planner.Calls - plannerBefore.Calls; got != 2 {
		t.Fatalf("planner calls delta = %d, want 2", got)
	}
	if got := planner.InputTokens - plannerBefore.InputTokens; got != 200 {
		t.Fatalf("planner input delta = %d, want 200", got)
	}
	if got := planner.CacheReadTokens - plannerBefore.CacheReadTokens; got != 160 {
		t.Fatalf("planner cache-read delta = %d, want 160", got)
	}
	if got := planner.OutputTokens - plannerBefore.OutputTokens; got != 10 {
		t.Fatalf("planner output delta = %d, want 10", got)
	}

	worker, workerBefore := after["worker"], before["worker"]
	if got := worker.Calls - workerBefore.Calls; got != 1 {
		t.Fatalf("worker calls delta = %d, want 1", got)
	}
	if got := worker.InputTokens - workerBefore.InputTokens; got != 50 {
		t.Fatalf("worker input delta = %d, want 50", got)
	}
	if got := worker.CacheReadTokens - workerBefore.CacheReadTokens; got != 10 {
		t.Fatalf("worker cache-read delta = %d, want 10", got)
	}
}

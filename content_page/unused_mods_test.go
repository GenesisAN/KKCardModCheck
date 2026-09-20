package content_page

import (
	"KKCardModCheck/config"
	"testing"
)

func TestUnusedScanRequiresCompleteCurrentResults(t *testing.T) {
	checkMu.Lock()
	oldPath := config.Instance.CardPath
	resetCheckStateLocked()
	config.Instance.CardPath = "cards"
	completedCardPath = "cards"
	checked = true
	progressParsed = 1
	usedMods = []string{"used"}
	checkMu.Unlock()
	t.Cleanup(func() { checkMu.Lock(); resetCheckStateLocked(); config.Instance.CardPath = oldPath; checkMu.Unlock() })
	ready, generation, used := unusedScanSnapshot()
	if !ready || len(used) != 1 {
		t.Fatal("valid scan rejected")
	}
	used[0] = "modified"
	if usedMods[0] != "used" {
		t.Fatal("snapshot aliases scan")
	}
	for _, scenario := range []string{"errors", "empty", "checking", "path", "unchecked"} {
		t.Run(scenario, func(t *testing.T) {
			checkMu.Lock()
			progressErrors = 0
			progressParsed = 1
			isChecking = false
			checked = true
			completedCardPath = "cards"
			switch scenario {
			case "errors":
				progressErrors = 1
			case "empty":
				progressParsed = 0
			case "checking":
				isChecking = true
			case "path":
				completedCardPath = "other"
			case "unchecked":
				checked = false
			}
			checkMu.Unlock()
			if ready, _, _ := unusedScanSnapshot(); ready {
				t.Fatal("incomplete/stale scan accepted")
			}
		})
	}
	checkMu.Lock()
	resetCheckStateLocked()
	checkMu.Unlock()
	if ready, next, _ := unusedScanSnapshot(); ready || next == generation {
		t.Fatal("reset did not invalidate preview")
	}
}

package content_page

import (
	"sort"
	"strings"

	"github.com/GenesisAN/illusionsCard/Base"
)

// summarizeCardMods returns the union of every referenced Zipmod and the
// subset that is absent from the local MOD database. KK scene cards expose the
// union of scene-level and embedded-character dependencies through the same
// CardInterface, so they can be summarized together with ordinary cards.
func summarizeCardMods(cards []Base.CardInterface, localGUIDs []string) (missing, used []string) {
	summary := newCardModSummary(localGUIDs)
	for _, card := range cards {
		summary.addCard(card)
	}
	return summary.lists()
}

type cardModSummary struct {
	local   map[string]struct{}
	used    map[string]struct{}
	missing map[string]struct{}
}

func newCardModSummary(localGUIDs []string) *cardModSummary {
	local := make(map[string]struct{}, len(localGUIDs))
	for _, guid := range localGUIDs {
		guid = strings.TrimSpace(guid)
		if guid != "" {
			local[guid] = struct{}{}
		}
	}
	return &cardModSummary{
		local:   local,
		used:    make(map[string]struct{}),
		missing: make(map[string]struct{}),
	}
}

func (summary *cardModSummary) addCard(card Base.CardInterface) {
	if card == nil {
		return
	}
	for _, guid := range card.GetZipmodsDependencies() {
		guid = strings.TrimSpace(guid)
		if guid == "" {
			continue
		}
		summary.used[guid] = struct{}{}
		if _, exists := summary.local[guid]; !exists {
			summary.missing[guid] = struct{}{}
		}
	}
}

func (summary *cardModSummary) lists() (missing, used []string) {
	return sortedModGUIDs(summary.missing), sortedModGUIDs(summary.used)
}

func sortedModGUIDs(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

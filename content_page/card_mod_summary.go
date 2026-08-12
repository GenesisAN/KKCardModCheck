package content_page

import (
	"sort"
	"strings"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
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
	dependencies := card.GetZipmodsDependencies()
	if scene, ok := card.(*KK.KKSceneCard); ok {
		dependencies = append(dependencies, sceneItemZipmodDependencies(scene)...)
	}
	for _, guid := range dependencies {
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

func sceneItemZipmodDependencies(scene *KK.KKSceneCard) []string {
	if scene == nil || scene.SceneExtended == nil {
		return nil
	}
	plugin := scene.SceneExtended[Base.UARExtID]
	if plugin == nil {
		plugin = scene.SceneExtended[Base.UARExtIDOld]
	}
	if plugin == nil {
		return nil
	}
	data, ok := plugin.Data.(map[string]interface{})
	if !ok {
		return nil
	}

	entries, ok := data["itemInfo"].([]interface{})
	if !ok {
		return nil
	}
	dependencies := make([]string, 0, len(entries))
	for _, entry := range entries {
		encoded, ok := entry.([]byte)
		if !ok {
			continue
		}
		var resolve Base.ResolveInfo
		if rest, err := resolve.UnmarshalMsg(encoded); err != nil || len(rest) != 0 {
			continue
		}
		if guid := strings.TrimSpace(resolve.GUID); guid != "" {
			dependencies = append(dependencies, guid)
		}
	}
	return dependencies
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

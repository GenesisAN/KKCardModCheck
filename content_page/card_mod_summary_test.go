package content_page

import (
	"reflect"
	"strings"
	"testing"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
)

type summaryTestCard struct {
	dependencies []string
}

func (card *summaryTestCard) CompareMissingZipMods([]string) []string { return nil }
func (card *summaryTestCard) GetZipmodsDependencies() []string        { return card.dependencies }
func (card *summaryTestCard) GetPath() string                         { return "test.png" }
func (card *summaryTestCard) GetDLLDependencies() []string            { return nil }

func TestSummarizeCardModsAggregatesAndSortsAllCards(t *testing.T) {
	cards := []Base.CardInterface{
		&summaryTestCard{dependencies: []string{"mod.c", "mod.a", "mod.a", ""}},
		&summaryTestCard{dependencies: []string{" mod.b ", "mod.c"}},
	}

	missing, used := summarizeCardMods(cards, []string{"mod.a"})
	if want := []string{"mod.a", "mod.b", "mod.c"}; !reflect.DeepEqual(used, want) {
		t.Fatalf("used MODs = %v, want %v", used, want)
	}
	if want := []string{"mod.b", "mod.c"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("missing MODs = %v, want %v", missing, want)
	}
}

func TestSummarizeCardModsSupportsSceneDependencyUnion(t *testing.T) {
	beachScene, err := (&Base.ResolveInfo{GUID: "BeachScene", Slot: 454271, LocalSlot: 100020938}).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	duplicateSceneMod, err := (&Base.ResolveInfo{GUID: "scene.mod"}).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	scene := &KK.KKSceneCard{Card: &Base.Card{ExtendedList: map[string]*Base.PluginDataEx{
		"scene-plugin": {
			RequiredZipmodGUIDs: []Base.ResolveInfo{{GUID: "scene.mod"}},
		},
		"embedded-character-plugin": {
			RequiredZipmodGUIDs: []Base.ResolveInfo{{GUID: "character.mod"}, {GUID: "shared.mod"}},
		},
	}}, SceneExtended: map[string]*Base.PluginData{
		Base.UARExtID: {
			Data: map[string]interface{}{
				"itemInfo": []interface{}{beachScene, duplicateSceneMod, "invalid"},
			},
		},
	}}

	missing, used := summarizeCardMods([]Base.CardInterface{scene}, []string{"shared.mod"})
	if want := []string{"BeachScene", "character.mod", "scene.mod", "shared.mod"}; !reflect.DeepEqual(used, want) {
		t.Fatalf("scene used MODs = %v, want %v", used, want)
	}
	if want := []string{"BeachScene", "character.mod", "scene.mod"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("scene missing MODs = %v, want %v", missing, want)
	}
	buildSceneTexts(scene)
	if !strings.Contains(sceneSummaryText, "BeachScene") || !strings.Contains(scenePluginsText, "BeachScene") {
		t.Fatalf("scene details omitted itemInfo dependency: summary=%q plugins=%q", sceneSummaryText, scenePluginsText)
	}
}

func TestSceneItemZipmodDependenciesSupportsLegacyUARID(t *testing.T) {
	encoded, err := (&Base.ResolveInfo{GUID: "legacy.scene.mod"}).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	scene := &KK.KKSceneCard{SceneExtended: map[string]*Base.PluginData{
		Base.UARExtIDOld: {Data: map[string]interface{}{"itemInfo": []interface{}{encoded}}},
	}}
	got := sceneItemZipmodDependencies(scene)
	if want := []string{"legacy.scene.mod"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("legacy scene dependencies = %v, want %v", got, want)
	}
}

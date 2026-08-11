package content_page

import (
	"reflect"
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
	scene := &KK.KKSceneCard{Card: &Base.Card{ExtendedList: map[string]*Base.PluginDataEx{
		"scene-plugin": {
			RequiredZipmodGUIDs: []Base.ResolveInfo{{GUID: "scene.mod"}},
		},
		"embedded-character-plugin": {
			RequiredZipmodGUIDs: []Base.ResolveInfo{{GUID: "character.mod"}, {GUID: "shared.mod"}},
		},
	}}}

	missing, used := summarizeCardMods([]Base.CardInterface{scene}, []string{"shared.mod"})
	if want := []string{"character.mod", "scene.mod", "shared.mod"}; !reflect.DeepEqual(used, want) {
		t.Fatalf("scene used MODs = %v, want %v", used, want)
	}
	if want := []string{"character.mod", "scene.mod"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("scene missing MODs = %v, want %v", missing, want)
	}
}

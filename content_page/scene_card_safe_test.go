package content_page

import (
	"os"
	"testing"

	ic "github.com/GenesisAN/illusionsCard"
	"github.com/GenesisAN/illusionsCard/KK"
)

func TestSceneReadyForExtractionRequiresParsedCurrentScene(t *testing.T) {
	originalPath := sceneCardPath
	originalParsedPath := parsedScenePath
	originalScene := parsedScene
	defer func() {
		sceneCardPath = originalPath
		parsedScenePath = originalParsedPath
		parsedScene = originalScene
	}()

	sceneCardPath = "scene.png"
	parsedScenePath = ""
	parsedScene = nil
	if sceneReadyForExtraction() {
		t.Fatal("scene must not be extractable before a successful parse")
	}
}

func TestSceneCardParseAndExtractSmoke(t *testing.T) {
	path := os.Getenv("KKCARDMODCHECK_SCENE_SMOKE")
	if path == "" {
		t.Skip("set KKCARDMODCHECK_SCENE_SMOKE to run the local scene-card smoke test")
	}

	card, err := ic.ReadCardFromPath(path)
	if err != nil {
		t.Fatalf("read scene card: %v", err)
	}
	scene, ok := card.(*KK.KKSceneCard)
	if !ok {
		t.Fatalf("expected KK scene card, got %T", card)
	}

	indexes := sceneExtractionIndexes(scene, true)
	for _, index := range indexes {
		if index < 0 || index >= len(scene.CharaCards) {
			t.Fatalf("character index %d outside [0,%d)", index, len(scene.CharaCards))
		}
	}
	if len(indexes) == 0 {
		t.Skip("scene parsed successfully but contains no extractable characters")
	}

	data, err := scene.ExtractCharaCard(indexes[0])
	if err != nil {
		t.Fatalf("extract first character in memory: %v", err)
	}
	if len(data) < 8 {
		t.Fatalf("extracted character data is unexpectedly short: %d bytes", len(data))
	}
}

package content_page

import g "github.com/AllenDang/giu"

// Keep toolbar controls consistent across classic pages.
const (
	classicButtonWidth        float32 = 144
	classicButtonHeight       float32 = 32
	classicInlineButtonWidth  float32 = 64
	classicInlineButtonHeight float32 = 28
)

func classicButton(label string) *g.ButtonWidget {
	return g.Button(label).Size(classicButtonWidth, classicButtonHeight)
}

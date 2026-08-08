package content_page

import (
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
)

type vectorIconV3 struct {
	id      string
	kind    string
	size    float32
	col     color.Color
	onClick func()
}

func iconV3(id, kind string, size float32, col color.Color) g.Widget {
	return &vectorIconV3{id: id, kind: kind, size: size, col: col}
}

func clickableIconV3(id, kind string, size float32, col color.Color, onClick func()) g.Widget {
	return &vectorIconV3{id: id, kind: kind, size: size, col: col, onClick: onClick}
}

func (w *vectorIconV3) Build() {
	pos := g.GetCursorScreenPos()
	g.InvisibleButton().ID(g.ID(w.id)).Size(w.size, w.size).OnClick(w.onClick).Build()
	c := g.GetCanvas()
	s := int(w.size)
	p := func(x, y int) image.Point { return pos.Add(image.Pt(x, y)) }
	left, top, right, bottom := 5, 5, s-5, s-5
	stroke := float32(2.2)

	switch w.kind {
	case "logo":
		c.AddRectFilled(p(3, 3), p(s-4, s-4), uiPrimary, 6, g.DrawFlagsRoundCornersAll)
		c.AddLine(p(9, s/2), p(14, s/2+4), color.White, 2.4)
		c.AddLine(p(14, s/2+4), p(22, 10), color.White, 2.4)
		c.AddCircleFilled(p(s-8, s-8), 6, uiSuccess)
		c.AddLine(p(s-11, s-8), p(s-8, s-5), color.White, 1.7)
		c.AddLine(p(s-8, s-5), p(s-4, s-11), color.White, 1.7)
	case "document":
		c.AddRect(p(left+3, top), p(right-3, bottom), w.col, 2, g.DrawFlagsRoundCornersAll, stroke)
		c.AddLine(p(left+8, top+9), p(right-8, top+9), w.col, stroke)
		c.AddLine(p(left+8, top+15), p(right-8, top+15), w.col, stroke)
		c.AddLine(p(left+8, top+21), p(right-11, top+21), w.col, stroke)
	case "scene":
		c.AddRect(p(left, top+2), p(right, bottom-2), w.col, 3, g.DrawFlagsRoundCornersAll, stroke)
		c.AddTriangle(p(left+5, bottom-7), p(s/2, top+11), p(right-5, bottom-7), w.col, stroke)
	case "mods":
		c.AddRect(p(left, top+7), p(s/2+3, bottom-5), w.col, 3, g.DrawFlagsRoundCornersAll, stroke)
		c.AddRect(p(s/2-2, top+2), p(right, bottom-10), w.col, 3, g.DrawFlagsRoundCornersAll, stroke)
		c.AddCircle(p(right-4, bottom-4), 6, w.col, 12, stroke)
	case "list":
		c.AddRect(p(left, top+3), p(right, bottom-3), w.col, 3, g.DrawFlagsRoundCornersAll, stroke)
		for y := top + 10; y <= bottom-7; y += 7 {
			c.AddCircleFilled(p(left+6, y), 1.5, w.col)
			c.AddLine(p(left+11, y), p(right-5, y), w.col, 1.8)
		}
	case "database":
		c.AddCircle(p(s/2, top+5), float32((right-left)/2), w.col, 18, stroke)
		c.AddLine(p(left, top+5), p(left, bottom-5), w.col, stroke)
		c.AddLine(p(right, top+5), p(right, bottom-5), w.col, stroke)
		c.PathClear()
		c.PathArcTo(p(s/2, bottom-5), float32((right-left)/2), 0, 3.15, 18)
		c.PathStroke(w.col, g.DrawFlagsNone, stroke)
		c.PathClear()
		c.PathArcTo(p(s/2, s/2), float32((right-left)/2), 0, 3.15, 18)
		c.PathStroke(w.col, g.DrawFlagsNone, stroke)
	case "settings":
		c.AddCircle(p(s/2, s/2), float32(s/3), w.col, 18, stroke)
		c.AddCircle(p(s/2, s/2), float32(s/8), w.col, 14, stroke)
		for _, d := range []image.Point{{0, -s/2 + 5}, {0, s/2 - 5}, {-s/2 + 5, 0}, {s/2 - 5, 0}} {
			c.AddLine(p(s/2, s/2), p(s/2+d.X, s/2+d.Y), w.col, stroke)
		}
	case "image":
		c.AddRect(p(left, top), p(right, bottom), w.col, 4, g.DrawFlagsRoundCornersAll, stroke)
		c.AddCircleFilled(p(right-7, top+7), 3, w.col)
		c.AddTriangle(p(left+4, bottom-5), p(s/2-2, top+13), p(right-4, bottom-5), w.col, stroke)
	case "warning":
		c.AddTriangle(p(s/2, top), p(right, bottom-2), p(left, bottom-2), w.col, stroke)
		c.AddLine(p(s/2, top+8), p(s/2, bottom-10), w.col, stroke)
		c.AddCircleFilled(p(s/2, bottom-6), 1.8, w.col)
	case "cube":
		c.AddQuad(p(s/2, top), p(right, top+7), p(s/2, top+14), p(left, top+7), w.col, stroke)
		c.AddLine(p(left, top+7), p(left, bottom-7), w.col, stroke)
		c.AddLine(p(right, top+7), p(right, bottom-7), w.col, stroke)
		c.AddLine(p(s/2, top+14), p(s/2, bottom), w.col, stroke)
		c.AddLine(p(left, bottom-7), p(s/2, bottom), w.col, stroke)
		c.AddLine(p(right, bottom-7), p(s/2, bottom), w.col, stroke)
	case "book":
		c.AddRect(p(left, top+3), p(s/2-1, bottom-2), w.col, 2, g.DrawFlagsRoundCornersAll, stroke)
		c.AddRect(p(s/2+1, top+3), p(right, bottom-2), w.col, 2, g.DrawFlagsRoundCornersAll, stroke)
		c.AddLine(p(s/2, top+5), p(s/2, bottom-2), w.col, stroke)
	default:
		c.AddCircle(p(s/2, s/2), float32(s/3), w.col, 18, stroke)
	}
}

func dashedRectV3(c *g.Canvas, min, max image.Point, col color.Color) {
	const dash = 8
	const gap = 6
	for x := min.X; x < max.X; x += dash + gap {
		x2 := x + dash
		if x2 > max.X {
			x2 = max.X
		}
		c.AddLine(image.Pt(x, min.Y), image.Pt(x2, min.Y), col, 1.5)
		c.AddLine(image.Pt(x, max.Y), image.Pt(x2, max.Y), col, 1.5)
	}
	for y := min.Y; y < max.Y; y += dash + gap {
		y2 := y + dash
		if y2 > max.Y {
			y2 = max.Y
		}
		c.AddLine(image.Pt(min.X, y), image.Pt(min.X, y2), col, 1.5)
		c.AddLine(image.Pt(max.X, y), image.Pt(max.X, y2), col, 1.5)
	}
}

func dropZoneV3(id string, height float32, widgets ...g.Widget) g.Widget {
	return g.Custom(func() {
		width, _ := g.GetAvailableRegion()
		pos := g.GetCursorScreenPos()
		g.Style().
			SetColor(g.StyleColorChildBg, color.RGBA{R: 0x10, G: 0x1B, B: 0x2C, A: 0xFF}).
			SetColor(g.StyleColorBorder, uiTransparent).
			SetStyle(g.StyleVarWindowPadding, 18, 18).
			SetStyleFloat(g.StyleVarChildRounding, 8).
			To(g.Child().ID(g.ID(id)).Size(width, height).Border(false).Flags(g.WindowFlagsNoScrollbar).Layout(widgets...)).Build()
		dashedRectV3(g.GetCanvas(), pos.Add(image.Pt(1, 1)), pos.Add(image.Pt(int(width)-2, int(height)-2)), color.RGBA{R: 0x72, G: 0x84, B: 0x9C, A: 0xFF})
	})
}

func placeholderCardV3(width, height float32) g.Widget {
	return g.Custom(func() {
		pos := g.GetCursorScreenPos()
		g.InvisibleButton().ID(g.ID("placeholder-card-v3")).Size(width, height).Build()
		c := g.GetCanvas()
		max := pos.Add(image.Pt(int(width), int(height)))
		c.AddRectFilled(pos, max, color.RGBA{R: 0x4A, G: 0x5D, B: 0x8D, A: 0xFF}, 9, g.DrawFlagsRoundCornersAll)
		for i := 0; i < 5; i++ {
			y1 := int(height) * i / 5
			y2 := int(height) * (i + 1) / 5
			shade := uint8(58 + i*8)
			c.AddRectFilled(pos.Add(image.Pt(0, y1)), pos.Add(image.Pt(int(width), y2)), color.RGBA{R: shade, G: shade + 12, B: shade + 38, A: 0xFF}, 0, g.DrawFlagsNone)
		}
		cx := int(width * .52)
		cy := int(height * .38)
		c.AddCircleFilled(pos.Add(image.Pt(cx, cy)), float32(width*.18), color.RGBA{R: 0x17, G: 0x22, B: 0x3C, A: 0xFF})
		c.AddTriangleFilled(pos.Add(image.Pt(int(width*.25), int(height*.92))), pos.Add(image.Pt(cx, int(height*.40))), pos.Add(image.Pt(int(width*.82), int(height*.92))), color.RGBA{R: 0x17, G: 0x22, B: 0x3C, A: 0xFF})
		c.AddRect(pos.Add(image.Pt(6, 6)), max.Sub(image.Pt(6, 6)), color.RGBA{R: 0xD5, G: 0xE1, B: 0xFF, A: 0xFF}, 7, g.DrawFlagsRoundCornersAll, 1.5)
	})
}

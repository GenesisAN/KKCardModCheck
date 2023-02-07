package Tools

import (
	"KKCardModCheck/IllusionsCard/Base"
	"bytes"
)

func get_png(file []byte) int {
	res1 := bytes.Index(file, Base.PngEndChunk)
	return res1
}

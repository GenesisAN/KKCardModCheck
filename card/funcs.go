package card

import "bytes"

var PngEndChunk = []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

var PngStartChunk = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D}

var PluginGUID = "com.bepis.bepinex.sideloader"

var UARExtID = "com.bepis.sideloader.universalautoresolver"

var UARExtIDOld = "EC.Core.Sideloader.UniversalAutoResolver"

func get_png(file []byte) int {
	res1 := bytes.Index(file, PngEndChunk)
	return res1
}

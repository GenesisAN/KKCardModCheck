package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"

	g "github.com/AllenDang/giu"
)

func configCardPathForScene() string {
	return config.Instance.CardPath
}

func classicSceneFilePicker() g.Widget {
	return util.FilePicker("场景卡文件", &sceneCardPath).ForFiles().Filters("*.png").LabelWidth(120).ButtonWidth(144)
}

func classicSceneOutputPicker() g.Widget {
	return util.FilePicker("输出目录", &sceneOutputDir).ForDirectories().LabelWidth(120).ButtonWidth(144)
}

func openSceneOutputDir() {
	util.OpenFile(sceneOutputDir)
}

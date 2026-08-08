package util

// BrowseFile 打开原生文件选择对话框，供自定义界面按钮复用。
func BrowseFile(title, defaultPath string, filters ...string) (string, error) {
	return showFileDialog(title, defaultPath, filters)
}

// BrowseDirectory 打开原生目录选择对话框，供自定义界面按钮复用。
func BrowseDirectory(title, defaultPath string) (string, error) {
	return showDirectoryDialog(title, defaultPath)
}

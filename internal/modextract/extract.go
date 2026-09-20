// Package modextract previews and moves unreferenced mod archives.
package modextract

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Relative string
	GUID     string
	info     os.FileInfo
}

type Plan struct {
	Source      string
	Destination string
	Entries     []Entry
	Skipped     []string
}

type ReadGUID func(string) (string, error)

func inside(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// canonical resolves even a destination whose final directories do not exist.
func canonical(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("请选择目录")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	parent := filepath.Dir(abs)
	if parent == abs {
		return "", err
	}
	parent, err = canonical(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(abs)), nil
}

func Preview(source, destination string, used []string, read ReadGUID) (*Plan, error) {
	source, err := canonical(source)
	if err != nil {
		return nil, err
	}
	destination, err = canonical(destination)
	if err != nil {
		return nil, err
	}
	if inside(source, destination) || inside(destination, source) {
		return nil, fmt.Errorf("MOD 目录与存放目录不能相同或互相包含")
	}
	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("MOD 路径必须是目录")
	}
	if info, err := os.Stat(destination); err == nil && !info.IsDir() {
		return nil, fmt.Errorf("存放路径必须是目录")
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	refs := make(map[string]bool, len(used))
	for _, guid := range used {
		refs[strings.TrimSpace(guid)] = true
	}
	plan := &Plan{Source: source, Destination: destination}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			plan.Skipped = append(plan.Skipped, path+": 跳过符号链接")
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".zipmod" && ext != ".zip" {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		guid, err := read(path)
		guid = strings.TrimSpace(guid)
		if err != nil || guid == "" {
			plan.Skipped = append(plan.Skipped, path+": 无法识别 GUID，保留原文件")
			return nil
		}
		if refs[guid] {
			return nil
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		plan.Entries = append(plan.Entries, Entry{Relative: rel, GUID: guid, info: info})
		return nil
	})
	return plan, err
}

// Move copies exclusively, syncs and closes the destination before removing the
// source. This also works across volumes and never overwrites an existing file.
func (p *Plan) Move(report func(string)) (moved, failed int) {
	for _, entry := range p.Entries {
		err := p.moveOne(entry)
		if err != nil {
			failed++
			report(fmt.Sprintf("失败 %s: %v", entry.Relative, err))
		} else {
			moved++
			report("已移动 " + entry.Relative)
		}
	}
	return
}

func (p *Plan) moveOne(entry Entry) error {
	source := filepath.Join(p.Source, entry.Relative)
	dest := filepath.Join(p.Destination, entry.Relative)
	resolvedSource, err := canonical(source)
	if err != nil {
		return err
	}
	resolvedDest, err := canonical(dest)
	if err != nil {
		return err
	}
	if resolvedSource != source || resolvedDest != dest || !inside(p.Source, source) || !inside(p.Destination, dest) {
		return fmt.Errorf("路径已变化，请重新预览")
	}
	current, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if !current.Mode().IsRegular() || !os.SameFile(entry.info, current) || current.Size() != entry.info.Size() || !current.ModTime().Equal(entry.info.ModTime()) {
		return fmt.Errorf("文件已变化，请重新预览")
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	if err = os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, current.Mode().Perm())
	if err != nil {
		return err
	}
	n, copyErr := io.Copy(out, in)
	if copyErr == nil {
		copyErr = out.Sync()
	}
	closeErr := out.Close()
	if copyErr == nil {
		copyErr = closeErr
	}
	if copyErr == nil && n != current.Size() {
		copyErr = fmt.Errorf("文件大小在复制时变化")
	}
	after, statErr := os.Lstat(source)
	if copyErr == nil && (statErr != nil || !os.SameFile(current, after) || after.Size() != current.Size() || !after.ModTime().Equal(current.ModTime())) {
		copyErr = fmt.Errorf("源文件在复制时变化")
	}
	if copyErr != nil {
		_ = os.Remove(dest)
		return copyErr
	}
	if err = in.Close(); err != nil {
		return err
	}
	if err = os.Remove(source); err != nil {
		return fmt.Errorf("已复制但无法移除源文件: %w", err)
	}
	return nil
}

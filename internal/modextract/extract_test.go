package modextract

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func put(t *testing.T, root, name, value string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(value), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func read(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

func TestPreviewAndMove(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "mods"), filepath.Join(root, "output")
	put(t, src, "used.zipmod", "used")
	put(t, src, "duplicate.ZIPMOD", "used")
	put(t, src, "a/spare.zipmod", "unused")
	put(t, src, "b/spare.zip", "another")
	put(t, src, "broken.zipmod", "")
	put(t, src, "notes.txt", "unused")
	plan, err := Preview(src, dst, []string{"used"}, read)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Entries) != 2 || len(plan.Skipped) != 1 {
		t.Fatalf("unexpected preview: %+v", plan)
	}
	moved, failed := plan.Move(func(string) {})
	if moved != 2 || failed != 0 {
		t.Fatalf("moved %d, failed %d", moved, failed)
	}
	for _, name := range []string{"a/spare.zipmod", "b/spare.zip"} {
		if _, err := os.Stat(filepath.Join(src, name)); !os.IsNotExist(err) {
			t.Fatal("source still exists", err)
		}
		if _, err := os.Stat(filepath.Join(dst, name)); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"used.zipmod", "duplicate.ZIPMOD", "broken.zipmod", "notes.txt"} {
		if _, err := os.Stat(filepath.Join(src, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestConflictAndChangedSourceArePreserved(t *testing.T) {
	for _, scenario := range []string{"conflict", "changed", "missing"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			src, dst := filepath.Join(root, "mods"), filepath.Join(root, "output")
			source := put(t, src, "a.zipmod", "spare")
			plan, err := Preview(src, dst, nil, read)
			if err != nil {
				t.Fatal(err)
			}
			switch scenario {
			case "conflict":
				put(t, dst, "a.zipmod", "existing")
			case "changed":
				put(t, src, "a.zipmod", "changed content")
			case "missing":
				if err := os.Remove(source); err != nil {
					t.Fatal(err)
				}
			}
			moved, failed := plan.Move(func(string) {})
			if moved != 0 || failed != 1 {
				t.Fatalf("moved %d failed %d", moved, failed)
			}
			if scenario != "missing" {
				if _, err := os.Stat(source); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "conflict" {
				b, _ := os.ReadFile(filepath.Join(dst, "a.zipmod"))
				if string(b) != "existing" {
					t.Fatal("overwritten")
				}
			}
		})
	}
}

func TestRejectOverlappingDirectories(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "mods")
	put(t, src, "a.zipmod", "a")
	destinations := []string{src, filepath.Join(src, "output"), root}
	if runtime.GOOS == "windows" {
		destinations = append(destinations, strings.ToUpper(src))
	}
	for _, dst := range destinations {
		if _, err := Preview(src, dst, nil, read); err == nil {
			t.Fatalf("accepted overlap %s", dst)
		}
	}
}

func TestSymlinkDestinationOutsideRoot(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "mods"), filepath.Join(root, "output")
	put(t, src, "nested/a.zipmod", "spare")
	plan, err := Preview(src, dst, nil, read)
	if err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(outside, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dst, "nested")); err != nil {
		t.Skip(err)
	}
	moved, failed := plan.Move(func(string) {})
	if moved != 0 || failed != 1 {
		t.Fatalf("moved %d failed %d", moved, failed)
	}
}

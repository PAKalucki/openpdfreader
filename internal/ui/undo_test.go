package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUndoManagerPushPop(t *testing.T) {
	m := newUndoManager(3)

	m.pushUndo("a")
	m.pushUndo("b")
	if !m.canUndo() {
		t.Fatal("expected canUndo=true")
	}

	got, ok := m.popUndo()
	if !ok || got != "b" {
		t.Fatalf("popUndo() = (%q,%v), want (b,true)", got, ok)
	}

	m.pushRedo("r1")
	if !m.canRedo() {
		t.Fatal("expected canRedo=true")
	}
	got, ok = m.popRedo()
	if !ok || got != "r1" {
		t.Fatalf("popRedo() = (%q,%v), want (r1,true)", got, ok)
	}
}

func TestUndoManagerClearsRedoOnNewUndo(t *testing.T) {
	m := newUndoManager(5)
	m.pushRedo("redo1")
	m.pushUndo("undo1")
	if m.canRedo() {
		t.Fatal("redo stack should be cleared when a new undo snapshot is pushed")
	}
}

func TestUndoManagerTrimRemovesOldFiles(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "1.pdf")
	f2 := filepath.Join(tmpDir, "2.pdf")
	f3 := filepath.Join(tmpDir, "3.pdf")

	if err := os.WriteFile(f1, []byte("1"), 0644); err != nil {
		t.Fatalf("WriteFile(f1): %v", err)
	}
	if err := os.WriteFile(f2, []byte("2"), 0644); err != nil {
		t.Fatalf("WriteFile(f2): %v", err)
	}
	if err := os.WriteFile(f3, []byte("3"), 0644); err != nil {
		t.Fatalf("WriteFile(f3): %v", err)
	}

	m := newUndoManager(2)
	m.pushUndo(f1)
	m.pushUndo(f2)
	m.pushUndo(f3)

	if _, err := os.Stat(f1); !os.IsNotExist(err) {
		t.Fatalf("expected oldest file to be removed, stat err=%v", err)
	}
}

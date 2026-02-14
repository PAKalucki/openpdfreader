package ui

import "os"

type undoManager struct {
	undo []string
	redo []string
	max  int
}

func newUndoManager(max int) *undoManager {
	if max <= 0 {
		max = 20
	}
	return &undoManager{max: max}
}

func (m *undoManager) canUndo() bool { return len(m.undo) > 0 }
func (m *undoManager) canRedo() bool { return len(m.redo) > 0 }

func (m *undoManager) pushUndo(path string) {
	if path == "" {
		return
	}
	m.undo = append(m.undo, path)
	m.clearRedo()

	for len(m.undo) > m.max {
		oldest := m.undo[0]
		m.undo = m.undo[1:]
		_ = os.Remove(oldest)
	}
}

func (m *undoManager) pushRedo(path string) {
	if path == "" {
		return
	}
	m.redo = append(m.redo, path)
}

func (m *undoManager) popUndo() (string, bool) {
	if len(m.undo) == 0 {
		return "", false
	}
	lastIdx := len(m.undo) - 1
	path := m.undo[lastIdx]
	m.undo = m.undo[:lastIdx]
	return path, true
}

func (m *undoManager) popRedo() (string, bool) {
	if len(m.redo) == 0 {
		return "", false
	}
	lastIdx := len(m.redo) - 1
	path := m.redo[lastIdx]
	m.redo = m.redo[:lastIdx]
	return path, true
}

func (m *undoManager) discardLastUndo() {
	path, ok := m.popUndo()
	if ok {
		_ = os.Remove(path)
	}
}

func (m *undoManager) clearAll() {
	for _, path := range m.undo {
		_ = os.Remove(path)
	}
	for _, path := range m.redo {
		_ = os.Remove(path)
	}
	m.undo = nil
	m.redo = nil
}

func (m *undoManager) clearRedo() {
	for _, path := range m.redo {
		_ = os.Remove(path)
	}
	m.redo = nil
}

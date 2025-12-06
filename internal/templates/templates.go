package templates

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"sync"
)

// Manager handles loading, caching, and rendering of HTML templates.
// It provides thread-safe template caching for improved performance.
type Manager struct {
	cache map[string]*template.Template
	dir   string
	mu    sync.RWMutex
}

// New creates a new template manager that loads templates from the specified directory.
func New(dir string) *Manager {
	return &Manager{
		cache: make(map[string]*template.Template),
		dir:   dir,
	}
}

// Get retrieves a parsed template from cache or loads it from disk.
// Templates are cached after first load for improved performance.
func (m *Manager) Get(name string) (*template.Template, error) {
	m.mu.RLock()
	tmpl, ok := m.cache[name]
	m.mu.RUnlock()

	if ok {
		return tmpl, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if tmpl, ok := m.cache[name]; ok {
		return tmpl, nil
	}

	path := filepath.Join(m.dir, name)
	tmpl, err := template.ParseFiles(path)

	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	m.cache[name] = tmpl
	return tmpl, nil
}

// Render executes a template with the provided data and writes to the writer.
func (m *Manager) Render(w io.Writer, name string, data any) error {
	tmpl, err := m.Get(name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// ClearCache removes all cached templates, forcing reload on next access.
// Useful during development when templates are being modified.
func (m *Manager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = make(map[string]*template.Template)
}

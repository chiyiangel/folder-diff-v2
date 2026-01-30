package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// KeyBinds handles keyboard shortcuts
type KeyBinds struct {
	app    *tview.Application
	layout *Layout
}

// NewKeyBinds creates a new key bindings handler
func NewKeyBinds(app *tview.Application, layout *Layout) *KeyBinds {
	return &KeyBinds{
		app:    app,
		layout: layout,
	}
}

// SetupGlobalKeys sets up global key bindings
func (k *KeyBinds) SetupGlobalKeys() {
	k.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle special keys
		switch event.Key() {
		case tcell.KeyEsc:
			k.app.Stop()
			return nil
		case tcell.KeyTab:
			k.layout.SwitchFocus()
			return nil
		case tcell.KeyLeft:
			k.layout.FocusLeft()
			return nil
		case tcell.KeyRight:
			k.layout.FocusRight()
			return nil
		case tcell.KeyCtrlC:
			k.app.Stop()
			return nil
		}

		// Handle rune keys
		switch event.Rune() {
		case 'q', 'Q':
			k.app.Stop()
			return nil
		case 'h', '?':
			k.layout.ShowHelp()
			return nil
		case 'd', 'D':
			k.jumpToNextDiff()
			return nil
		}

		return event
	})
}

// jumpToNextDiff jumps to the next difference in the current tree
func (k *KeyBinds) jumpToNextDiff() {
	// This is a placeholder for jumping to next diff
	// Implementation would traverse the tree and find next non-identical file
	// For now, we just let normal navigation handle it
}

// GetKeyHelp returns a formatted string of all key bindings
func GetKeyHelp() string {
	return `
╔══════════════════════════════════════════════════════════╗
║                    Keyboard Shortcuts                     ║
╠══════════════════════════════════════════════════════════╣
║  Navigation                                               ║
║    ↑ / k         Move selection up                        ║
║    ↓ / j         Move selection down                      ║
║    ← / →         Switch between panels                    ║
║    Tab           Toggle panel focus                       ║
║    Space         Expand/collapse folder                   ║
║    Enter         Expand/collapse folder                   ║
║    d             Jump to next difference                  ║
║                                                           ║
║  General                                                  ║
║    h / ?         Show this help                           ║
║    q / Esc       Quit application                         ║
║    Ctrl+C        Force quit                               ║
║                                                           ║
║  Status Legend                                            ║
║    ✓ (green)     Identical files                          ║
║    ~ (red)       Modified files                           ║
║    + (blue)      New files (in target only)               ║
║    - (gray)      Deleted files (in source only)           ║
╚══════════════════════════════════════════════════════════╝
`
}

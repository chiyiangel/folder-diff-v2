package tui

import (
	"folder-diff-v2/internal/compare"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App represents the TUI application
type App struct {
	app       *tview.Application
	layout    *Layout
	result    *compare.ComparisonResult
	sourceDir string
	targetDir string
}

// NewApp creates a new TUI application
func NewApp(result *compare.ComparisonResult, sourceDir, targetDir string) *App {
	return &App{
		app:       tview.NewApplication(),
		result:    result,
		sourceDir: sourceDir,
		targetDir: targetDir,
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	// Build tree structures from flat file lists
	sourceTree := BuildTree(a.result.SourceFiles, a.sourceDir)
	targetTree := BuildTree(a.result.TargetFiles, a.targetDir)

	// Create synchronized layout
	a.layout = NewLayout(a.app, sourceTree, targetTree, a.sourceDir, a.targetDir)

	// Set up global key bindings
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			a.app.Stop()
			return nil
		case tcell.KeyUp:
			a.layout.MoveUp()
			return nil
		case tcell.KeyDown:
			a.layout.MoveDown()
			return nil
		case tcell.KeyEnter:
			a.layout.ToggleExpand()
			return nil
		case tcell.KeyCtrlC:
			a.app.Stop()
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			a.app.Stop()
			return nil
		case 'h', '?':
			a.layout.ShowHelp()
			return nil
		case ' ':
			a.layout.ToggleExpand()
			return nil
		case 'd', 'D':
			a.layout.JumpToNextDiff()
			return nil
		case 'k':
			a.layout.MoveUp()
			return nil
		case 'j':
			a.layout.MoveDown()
			return nil
		}

		return event
	})

	// Set root and run
	a.app.SetRoot(a.layout.GetRoot(), true)
	a.app.EnableMouse(true)

	return a.app.Run()
}

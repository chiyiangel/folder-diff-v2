package tui

import (
	"fmt"

	"folder-diff-v2/internal/compare"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout manages the TUI layout
type Layout struct {
	app         *tview.Application
	root        *tview.Flex
	sourceTree  *TreeView
	targetTree  *TreeView
	statusBar   *tview.TextView
	focusLeft   bool
	helpModal   *tview.Modal
}

// NewLayout creates a new layout with source and target trees
func NewLayout(app *tview.Application, sourceRoot, targetRoot *compare.FileInfo, sourceDir, targetDir string) *Layout {
	l := &Layout{
		app:       app,
		focusLeft: true,
	}

	// Create tree views
	l.sourceTree = NewTreeView(sourceRoot, "Source: "+sourceDir)
	l.targetTree = NewTreeView(targetRoot, "Target: "+targetDir)

	// Create status bar
	l.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]↑↓[white] Navigate  [yellow]←→[white] Switch Panel  [yellow]Space[white] Expand/Collapse  [yellow]d[white] Next Diff  [yellow]h/?[white] Help  [yellow]q[white] Quit   |   [green]✓[white] Same  [red]~[white] Modified  [blue]+[white] New  [gray]-[white] Deleted")

	// Create title bar
	titleBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[::b]📁 Folder Diff - TUI Mode[::-]")
	titleBar.SetBackgroundColor(tcell.ColorDarkBlue)

	// Create main content with two panels
	content := tview.NewFlex().
		AddItem(l.sourceTree.GetPrimitive(), 0, 1, true).
		AddItem(l.targetTree.GetPrimitive(), 0, 1, false)

	// Create status bar container
	statusContainer := tview.NewFlex().
		AddItem(l.statusBar, 0, 1, false)
	statusContainer.SetBorder(false)

	// Assemble layout
	l.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 1, 0, false).
		AddItem(content, 0, 1, true).
		AddItem(statusContainer, 1, 0, false)

	// Create help modal
	l.helpModal = tview.NewModal().
		SetText(l.getHelpText()).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.app.SetRoot(l.root, true)
			if l.focusLeft {
				l.app.SetFocus(l.sourceTree.GetPrimitive())
			} else {
				l.app.SetFocus(l.targetTree.GetPrimitive())
			}
		})

	// Set initial focus
	app.SetFocus(l.sourceTree.GetPrimitive())

	return l
}

// GetRoot returns the root primitive
func (l *Layout) GetRoot() tview.Primitive {
	return l.root
}

// SwitchFocus switches focus between panels
func (l *Layout) SwitchFocus() {
	l.focusLeft = !l.focusLeft
	if l.focusLeft {
		l.app.SetFocus(l.sourceTree.GetPrimitive())
	} else {
		l.app.SetFocus(l.targetTree.GetPrimitive())
	}
}

// FocusLeft focuses the left panel
func (l *Layout) FocusLeft() {
	l.focusLeft = true
	l.app.SetFocus(l.sourceTree.GetPrimitive())
}

// FocusRight focuses the right panel
func (l *Layout) FocusRight() {
	l.focusLeft = false
	l.app.SetFocus(l.targetTree.GetPrimitive())
}

// ShowHelp displays the help modal
func (l *Layout) ShowHelp() {
	l.app.SetRoot(l.helpModal, true)
}

// getHelpText returns the help text
func (l *Layout) getHelpText() string {
	return fmt.Sprintf(`Folder Diff - Keyboard Shortcuts

Navigation:
  ↑/↓        Move up/down in tree
  ←/→        Switch between panels
  Tab        Toggle panel focus
  Space      Expand/collapse folder
  Enter      Expand/collapse folder
  d          Jump to next difference

Display:
  h / ?      Show this help
  q / Esc    Quit application

Legend:
  ✓ (green)  Identical files
  ~ (red)    Modified files
  + (blue)   New files (target only)
  - (gray)   Deleted files (source only)
`)
}

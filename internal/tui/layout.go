package tui

import (
	"fmt"
	"strings"

	"folder-diff-v2/internal/compare"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout manages the synchronized dual-pane TUI layout
type Layout struct {
	app           *tview.Application
	root          *tview.Flex
	sourceView    *tview.TextView
	targetView    *tview.TextView
	statusBar     *tview.TextView
	helpModal     *tview.Modal
	syncTree      *SyncNode
	flatNodes     []*SyncNode
	currentIndex  int
	sourceDir     string
	targetDir     string
}

// NewLayout creates a new synchronized layout
func NewLayout(app *tview.Application, sourceRoot, targetRoot *compare.FileInfo, sourceDir, targetDir string) *Layout {
	l := &Layout{
		app:          app,
		currentIndex: 0,
		sourceDir:    sourceDir,
		targetDir:    targetDir,
	}

	// Build synchronized tree
	l.syncTree = BuildSyncTree(sourceRoot, targetRoot)
	l.flatNodes = FlattenTree(l.syncTree)

	// Create text views for both panels
	l.sourceView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	l.sourceView.SetBorder(true).SetTitle(" Source: " + sourceDir + " ")

	l.targetView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	l.targetView.SetBorder(true).SetTitle(" Target: " + targetDir + " ")

	// Create status bar
	l.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]↑↓[white] Navigate  [yellow]Space[white] Expand/Collapse  [yellow]d[white] Next Diff  [yellow]h/?[white] Help  [yellow]q[white] Quit   |   [green]✓[white] Same  [red]~[white] Modified  [blue]+[white] New  [gray]-[white] Deleted")

	// Create title bar
	titleBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[::b]📁 Folder Diff - Synchronized View[::-]")
	titleBar.SetBackgroundColor(tcell.ColorDarkBlue)

	// Create main content with two panels
	content := tview.NewFlex().
		AddItem(l.sourceView, 0, 1, false).
		AddItem(l.targetView, 0, 1, false)

	// Assemble layout
	l.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 1, 0, false).
		AddItem(content, 0, 1, false).
		AddItem(l.statusBar, 1, 0, false)

	// Create help modal
	l.helpModal = tview.NewModal().
		SetText(l.getHelpText()).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.app.SetRoot(l.root, true)
		})

	// Initial render
	l.render()

	return l
}

// GetRoot returns the root primitive
func (l *Layout) GetRoot() tview.Primitive {
	return l.root
}

// render updates both views with current selection
func (l *Layout) render() {
	sourceText := ""
	targetText := ""

	for i, node := range l.flatNodes {
		level := l.getLevel(node)
		indent := strings.Repeat("  ", level)
		
		selected := i == l.currentIndex
		prefix := "  "
		if selected {
			prefix = "> "
		}

		// Render source side
		sourceText += l.renderNode(node, indent, prefix, true, selected) + "\n"
		
		// Render target side
		targetText += l.renderNode(node, indent, prefix, false, selected) + "\n"
	}

	l.sourceView.SetText(sourceText)
	l.targetView.SetText(targetText)

	// Scroll to current selection
	l.sourceView.ScrollTo(l.currentIndex, 0)
	l.targetView.ScrollTo(l.currentIndex, 0)
}

// renderNode renders a single node for source or target side
func (l *Layout) renderNode(node *SyncNode, indent, prefix string, isSource, selected bool) string {
	var icon, statusIcon, text string
	var color string

	// Choose file info
	var file *compare.FileInfo
	if isSource {
		file = node.SourceFile
	} else {
		file = node.TargetFile
	}

	// Determine icon
	if node.IsDir {
		if node.Expanded {
			icon = "📂"
		} else {
			icon = "📁"
		}
	} else {
		icon = "📄"
	}

	// Handle non-existent files
	if file == nil {
		text = "[Not exists]"
		if isSource && node.Status == compare.New {
			color = "gray"
			statusIcon = " +"
		} else if !isSource && node.Status == compare.Deleted {
			color = "gray"
			statusIcon = " -"
		} else {
			color = "gray"
			statusIcon = ""
		}
		
		if selected {
			return fmt.Sprintf("%s%s[black:white]%s %s%s[-:-]", prefix, indent, icon, text, statusIcon)
		}
		return fmt.Sprintf("%s%s[%s]%s %s%s[-]", prefix, indent, color, icon, text, statusIcon)
	}

	// Set status icon and color
	switch node.Status {
	case compare.Identical:
		statusIcon = " ✓"
		color = "green"
	case compare.Modified:
		statusIcon = " ~"
		color = "red"
	case compare.New:
		statusIcon = " +"
		color = "blue"
	case compare.Deleted:
		statusIcon = " -"
		color = "gray"
	}

	// Highlight if selected
	if selected {
		return fmt.Sprintf("%s%s[black:white]%s %s%s[-:-]", prefix, indent, icon, node.Name, statusIcon)
	}

	return fmt.Sprintf("%s%s[%s]%s %s%s[-]", prefix, indent, color, icon, node.Name, statusIcon)
}

// getLevel calculates the depth level of a node
func (l *Layout) getLevel(node *SyncNode) int {
	level := 0
	current := node.Parent
	for current != nil && current.RelPath != "." {
		level++
		current = current.Parent
	}
	return level
}

// MoveUp moves selection up
func (l *Layout) MoveUp() {
	if l.currentIndex > 0 {
		l.currentIndex--
		l.render()
	}
}

// MoveDown moves selection down
func (l *Layout) MoveDown() {
	if l.currentIndex < len(l.flatNodes)-1 {
		l.currentIndex++
		l.render()
	}
}

// ToggleExpand toggles expand/collapse for current directory
func (l *Layout) ToggleExpand() {
	if l.currentIndex < 0 || l.currentIndex >= len(l.flatNodes) {
		return
	}

	node := l.flatNodes[l.currentIndex]
	if node.IsDir {
		node.Expanded = !node.Expanded
		l.flatNodes = FlattenTree(l.syncTree)
		l.render()
	}
}

// JumpToNextDiff jumps to the next file with differences
func (l *Layout) JumpToNextDiff() {
	start := l.currentIndex + 1
	for i := start; i < len(l.flatNodes); i++ {
		if l.flatNodes[i].Status != compare.Identical {
			l.currentIndex = i
			l.render()
			return
		}
	}
	// Wrap around
	for i := 0; i < start; i++ {
		if l.flatNodes[i].Status != compare.Identical {
			l.currentIndex = i
			l.render()
			return
		}
	}
}

// ShowHelp displays the help modal
func (l *Layout) ShowHelp() {
	l.app.SetRoot(l.helpModal, true)
}

// getHelpText returns the help text
func (l *Layout) getHelpText() string {
	return `Folder Diff - Synchronized Navigation

Keyboard Shortcuts:

Navigation:
  ↑/↓        Move selection up/down (both panels)
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
  
Note: Both panels are synchronized - navigation 
affects both sides simultaneously.
`
}


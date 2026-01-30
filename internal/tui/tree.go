package tui

import (
	"path/filepath"
	"sort"
	"strings"

	"folder-diff-v2/internal/compare"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TreeView represents a file tree view
type TreeView struct {
	tree  *tview.TreeView
	root  *compare.FileInfo
	title string
}

// NewTreeView creates a new tree view
func NewTreeView(root *compare.FileInfo, title string) *TreeView {
	tv := &TreeView{
		tree:  tview.NewTreeView(),
		root:  root,
		title: title,
	}

	// Create root node
	rootNode := tv.createNode(root)
	tv.tree.SetRoot(rootNode)
	tv.tree.SetCurrentNode(rootNode)

	// Set up tree appearance
	tv.tree.SetBorder(true)
	tv.tree.SetTitle(" " + title + " ")
	tv.tree.SetTitleAlign(tview.AlignLeft)

	// Set selected function for expand/collapse
	tv.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		fileInfo := ref.(*compare.FileInfo)
		if fileInfo.IsDir {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tv
}

// createNode creates a tree node from FileInfo
func (tv *TreeView) createNode(info *compare.FileInfo) *tview.TreeNode {
	// Determine display text and color
	text, color := tv.formatNode(info)

	node := tview.NewTreeNode(text).
		SetReference(info).
		SetSelectable(true).
		SetColor(color)

	if info.IsDir {
		node.SetExpanded(true)
		// Add children sorted (dirs first, then files)
		children := tv.sortChildren(info.Children)
		for _, child := range children {
			childNode := tv.createNode(child)
			node.AddChild(childNode)
		}
	}

	return node
}

// formatNode returns the display text and color for a node
func (tv *TreeView) formatNode(info *compare.FileInfo) (string, tcell.Color) {
	var icon, statusIcon string
	var color tcell.Color

	// Set icon based on type
	if info.IsDir {
		icon = "📁 "
	} else {
		icon = "📄 "
	}

	// Set status icon and color
	switch info.Status {
	case compare.Identical:
		statusIcon = " ✓"
		color = tcell.ColorGreen
	case compare.Modified:
		statusIcon = " ~"
		color = tcell.ColorRed
	case compare.New:
		statusIcon = " +"
		color = tcell.ColorBlue
	case compare.Deleted:
		statusIcon = " -"
		color = tcell.ColorGray
	default:
		statusIcon = ""
		color = tcell.ColorWhite
	}

	name := info.Name
	if name == "" || name == "." {
		name = filepath.Base(info.Path)
	}

	return icon + name + statusIcon, color
}

// sortChildren sorts children (directories first, then alphabetically)
func (tv *TreeView) sortChildren(children []*compare.FileInfo) []*compare.FileInfo {
	sorted := make([]*compare.FileInfo, len(children))
	copy(sorted, children)

	sort.Slice(sorted, func(i, j int) bool {
		// Directories first
		if sorted[i].IsDir != sorted[j].IsDir {
			return sorted[i].IsDir
		}
		// Then alphabetically
		return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
	})

	return sorted
}

// GetPrimitive returns the tree primitive
func (tv *TreeView) GetPrimitive() tview.Primitive {
	return tv.tree
}

// BuildTree builds a tree structure from flat file list
func BuildTree(files []*compare.FileInfo, rootPath string) *compare.FileInfo {
	// Create root node
	root := &compare.FileInfo{
		Path:     rootPath,
		RelPath:  ".",
		Name:     filepath.Base(rootPath),
		IsDir:    true,
		Status:   compare.Identical,
		Children: []*compare.FileInfo{},
		Expanded: true,
	}

	// Map to track created directories
	dirMap := make(map[string]*compare.FileInfo)
	dirMap["."] = root

	// Sort files by path to ensure parents are created first
	sortedFiles := make([]*compare.FileInfo, len(files))
	copy(sortedFiles, files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].RelPath < sortedFiles[j].RelPath
	})

	for _, file := range sortedFiles {
		if file.RelPath == "." {
			continue
		}

		// Set name
		file.Name = filepath.Base(file.RelPath)

		// Find or create parent
		parentPath := filepath.Dir(file.RelPath)
		if parentPath == "" {
			parentPath = "."
		}

		parent, exists := dirMap[parentPath]
		if !exists {
			// Create intermediate directories
			parent = ensureParentDirs(dirMap, parentPath, root)
		}

		// Add to parent
		file.Parent = parent
		parent.Children = append(parent.Children, file)

		// If this is a directory, add to map
		if file.IsDir {
			dirMap[file.RelPath] = file
		}
	}

	return root
}

// ensureParentDirs creates parent directories as needed
func ensureParentDirs(dirMap map[string]*compare.FileInfo, path string, root *compare.FileInfo) *compare.FileInfo {
	if path == "." || path == "" {
		return root
	}

	if dir, exists := dirMap[path]; exists {
		return dir
	}

	// Create parent first
	parentPath := filepath.Dir(path)
	parent := ensureParentDirs(dirMap, parentPath, root)

	// Create this directory
	dir := &compare.FileInfo{
		RelPath:  path,
		Name:     filepath.Base(path),
		IsDir:    true,
		Status:   compare.Identical,
		Children: []*compare.FileInfo{},
		Parent:   parent,
		Expanded: true,
	}

	parent.Children = append(parent.Children, dir)
	dirMap[path] = dir

	return dir
}

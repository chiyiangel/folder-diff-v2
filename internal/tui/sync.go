package tui

import (
	"path/filepath"
	"sort"
	"strings"

	"folder-diff-v2/internal/compare"
)

// SyncNode represents a synchronized pair of nodes from source and target
type SyncNode struct {
	RelPath      string
	Name         string
	IsDir        bool
	SourceFile   *compare.FileInfo // nil if not exists in source
	TargetFile   *compare.FileInfo // nil if not exists in target
	Status       compare.FileStatus
	Children     []*SyncNode
	Parent       *SyncNode
	Expanded     bool
}

// BuildSyncTree creates a synchronized tree from source and target files
func BuildSyncTree(sourceRoot, targetRoot *compare.FileInfo) *SyncNode {
	root := &SyncNode{
		RelPath:    ".",
		Name:       ".",
		IsDir:      true,
		SourceFile: sourceRoot,
		TargetFile: targetRoot,
		Status:     compare.Identical,
		Children:   []*SyncNode{},
		Expanded:   true,
	}

	// Collect all unique paths from both trees
	pathMap := make(map[string]*SyncNode)
	pathMap["."] = root

	// Process source files
	if sourceRoot != nil {
		collectPaths(sourceRoot, pathMap, true, root)
	}

	// Process target files
	if targetRoot != nil {
		collectPaths(targetRoot, pathMap, false, root)
	}

	// Build parent-child relationships
	buildHierarchy(pathMap, root)

	return root
}

// collectPaths collects all paths from a tree
func collectPaths(node *compare.FileInfo, pathMap map[string]*SyncNode, isSource bool, root *SyncNode) {
	if node == nil {
		return
	}

	for _, child := range node.Children {
		if child.RelPath == "." {
			continue
		}

		syncNode, exists := pathMap[child.RelPath]
		if !exists {
			syncNode = &SyncNode{
				RelPath:  child.RelPath,
				Name:     child.Name,
				IsDir:    child.IsDir,
				Status:   child.Status,
				Children: []*SyncNode{},
				Expanded: true,
			}
			pathMap[child.RelPath] = syncNode
		}

		// Assign to source or target
		if isSource {
			syncNode.SourceFile = child
		} else {
			syncNode.TargetFile = child
		}

		// Update status based on both sides
		syncNode.Status = determineStatus(syncNode.SourceFile, syncNode.TargetFile)
		syncNode.IsDir = (syncNode.SourceFile != nil && syncNode.SourceFile.IsDir) ||
			(syncNode.TargetFile != nil && syncNode.TargetFile.IsDir)

		// Recursively process children
		collectPaths(child, pathMap, isSource, root)
	}
}

// buildHierarchy establishes parent-child relationships
func buildHierarchy(pathMap map[string]*SyncNode, root *SyncNode) {
	// Get all paths sorted
	var paths []string
	for path := range pathMap {
		if path != "." {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)

	// Assign children to parents
	for _, path := range paths {
		node := pathMap[path]
		parentPath := filepath.Dir(path)
		if parentPath == "" {
			parentPath = "."
		}

		parent, exists := pathMap[parentPath]
		if !exists {
			parent = root
		}

		node.Parent = parent
		if !contains(parent.Children, node) {
			parent.Children = append(parent.Children, node)
		}
	}

	// Sort children in each node
	sortChildren(root)
}

// sortChildren recursively sorts children (dirs first, then alphabetically)
func sortChildren(node *SyncNode) {
	sort.Slice(node.Children, func(i, j int) bool {
		// Directories first
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		// Then alphabetically
		return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortChildren(child)
	}
}

// determineStatus determines the status based on source and target files
func determineStatus(source, target *compare.FileInfo) compare.FileStatus {
	if source == nil && target != nil {
		return compare.New
	}
	if source != nil && target == nil {
		return compare.Deleted
	}
	if source != nil && target != nil {
		// Both exist - check their status
		if source.Status == compare.Modified || target.Status == compare.Modified {
			return compare.Modified
		}
		return compare.Identical
	}
	return compare.Identical
}

// contains checks if a node is in the children slice
func contains(children []*SyncNode, node *SyncNode) bool {
	for _, child := range children {
		if child.RelPath == node.RelPath {
			return true
		}
	}
	return false
}

// FlattenTree converts tree to flat list for display
func FlattenTree(root *SyncNode) []*SyncNode {
	var result []*SyncNode
	flattenRecursive(root, &result, 0)
	return result
}

// flattenRecursive recursively flattens the tree
func flattenRecursive(node *SyncNode, result *[]*SyncNode, level int) {
	if node == nil {
		return
	}

	// Skip root "."
	if node.RelPath != "." {
		*result = append(*result, node)
	}

	// Only process children if expanded
	if node.Expanded {
		for _, child := range node.Children {
			flattenRecursive(child, result, level+1)
		}
	}
}

// BuildTree builds a tree structure from flat file list (for compatibility)
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
			parent = ensureParentDirsForFile(dirMap, parentPath, root)
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

// ensureParentDirsForFile creates parent directories as needed
func ensureParentDirsForFile(dirMap map[string]*compare.FileInfo, path string, root *compare.FileInfo) *compare.FileInfo {
	if path == "." || path == "" {
		return root
	}

	if dir, exists := dirMap[path]; exists {
		return dir
	}

	// Create parent first
	parentPath := filepath.Dir(path)
	parent := ensureParentDirsForFile(dirMap, parentPath, root)

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

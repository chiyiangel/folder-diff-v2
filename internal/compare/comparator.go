package compare

type Comparator struct {
	mode ComparisonMode
}

func NewComparator(mode ComparisonMode) *Comparator {
	return &Comparator{mode: mode}
}

func (c *Comparator) Compare(source, target []*FileInfo) *ComparisonResult {
	result := &ComparisonResult{
		SourceFiles: source,
		TargetFiles: target,
		Mode:        c.mode,
	}

	sourceMap := make(map[string]*FileInfo)
	targetMap := make(map[string]*FileInfo)

	for _, file := range source {
		sourceMap[file.RelPath] = file
	}

	for _, file := range target {
		targetMap[file.RelPath] = file

		if sourceFile, exists := sourceMap[file.RelPath]; exists {
			if file.IsDir {
				file.Status = Identical
				sourceFile.Status = Identical
				continue
			}

			if c.mode == HashMode {
				if file.Hash == sourceFile.Hash {
					file.Status = Identical
					sourceFile.Status = Identical
				} else {
					file.Status = Modified
					sourceFile.Status = Modified
				}
			} else {
				file.Status = Identical
				sourceFile.Status = Identical
			}
		} else {
			file.Status = New
		}
	}

	// Mark deleted files
	for _, file := range source {
		if _, exists := targetMap[file.RelPath]; !exists {
			file.Status = Deleted
		}
	}

	return result
}

package compare

// ComparisonMode defines how files should be compared
type ComparisonMode string

const (
	HashMode     ComparisonMode = "hash"
	FilenameMode ComparisonMode = "filename"
)

// FileStatus represents the comparison status of a file
type FileStatus string

const (
	Identical FileStatus = "identical"
	Modified  FileStatus = "modified"
	New       FileStatus = "new"
	Deleted   FileStatus = "deleted"
)

// FileInfo represents a file or directory in the comparison
type FileInfo struct {
	Path     string
	RelPath  string
	Hash     string
	Status   FileStatus
	IsDir    bool
	Children []*FileInfo
}

// ComparisonResult contains the complete comparison results
type ComparisonResult struct {
	SourceRoot     string
	TargetRoot     string
	SourceFiles    []*FileInfo
	TargetFiles    []*FileInfo
	Mode           ComparisonMode
	ExcludePattern []string
}

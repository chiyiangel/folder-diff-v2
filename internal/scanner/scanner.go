package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"folder-diff-v2/internal/compare"
)

type Scanner struct {
	excludePatterns []string
}

func NewScanner(excludePatterns []string) *Scanner {
	return &Scanner{
		excludePatterns: excludePatterns,
	}
}

func (s *Scanner) shouldExclude(path string) bool {
	for _, pattern := range s.excludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

func (s *Scanner) calculateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *Scanner) ScanDirectory(root string) ([]*compare.FileInfo, error) {
	var files []*compare.FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if s.shouldExclude(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		fileInfo := &compare.FileInfo{
			Path:    path,
			RelPath: relPath,
			IsDir:   info.IsDir(),
		}

		if !info.IsDir() {
			hash, err := s.calculateHash(path)
			if err != nil {
				return err
			}
			fileInfo.Hash = hash
		}

		files = append(files, fileInfo)
		return nil
	})

	return files, err
}

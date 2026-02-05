package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FilesystemScanner implements the Scanner interface for local filesystem scanning
type FilesystemScanner struct {
	paths      []string
	parser     FileParser
	options    ScanOptions
	extensions map[string]bool
}

// NewFilesystemScanner creates a new scanner with the given paths and parser
func NewFilesystemScanner(paths []string, parser FileParser, opts ...ScanOptions) *FilesystemScanner {
	options := DefaultScanOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// Build extension map from parser
	extensions := make(map[string]bool)
	for _, ext := range parser.GetExtensions() {
		extensions[strings.ToLower(ext)] = true
	}

	return &FilesystemScanner{
		paths:      paths,
		parser:     parser,
		options:    options,
		extensions: extensions,
	}
}

// Scan scans all configured paths and returns discovered media files
func (s *FilesystemScanner) Scan(ctx context.Context) ([]ScanResult, error) {
	results, _, err := s.ScanWithSummary(ctx)
	return results, err
}

// ScanWithSummary scans and returns both results and statistics
func (s *FilesystemScanner) ScanWithSummary(ctx context.Context) ([]ScanResult, *ScanSummary, error) {
	// Check context before starting
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	var allResults []ScanResult
	summary := &ScanSummary{}

	for _, path := range s.paths {
		// Check context between paths
		select {
		case <-ctx.Done():
			return allResults, summary, ctx.Err()
		default:
		}

		results, err := s.ScanPath(ctx, path)
		if err != nil {
			// For context cancellation, return the error immediately
			if ctx.Err() != nil {
				return allResults, summary, ctx.Err()
			}
			summary.Errors = append(summary.Errors, fmt.Errorf("path %s: %w", path, err))
			continue
		}

		for _, result := range results {
			summary.TotalFiles++
			if result.IsMedia {
				summary.MediaFiles++
				if result.ParsedTitle != "" {
					summary.ParsedFiles++
				}
			} else {
				summary.SkippedFiles++
			}
			if result.Error != nil {
				summary.FailedParses++
			}
		}

		allResults = append(allResults, results...)
	}

	return allResults, summary, nil
}

// ScanPath scans a single path
func (s *FilesystemScanner) ScanPath(ctx context.Context, path string) ([]ScanResult, error) {
	// Check context before starting
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var results []ScanResult

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, walkErr error) error {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if walkErr != nil {
			// For the root directory, return the error (e.g., path doesn't exist)
			if filePath == path {
				return walkErr
			}
			// For other errors (permission issues, etc.), skip the file/dir
			return nil
		}

		// Handle directories
		if d.IsDir() {
			// Skip hidden directories unless configured otherwise
			if !s.options.IncludeHidden && strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}

			// Check exclude patterns
			for _, pattern := range s.options.ExcludePatterns {
				if matched, _ := filepath.Match(pattern, d.Name()); matched {
					return filepath.SkipDir
				}
			}

			// Check max depth
			if s.options.MaxDepth > 0 {
				relPath, _ := filepath.Rel(path, filePath)
				// Calculate depth: "." is 0, "a" is 1, "a/b" is 2
				var depth int
				if relPath != "." {
					depth = strings.Count(relPath, string(filepath.Separator)) + 1
				}
				if depth > s.options.MaxDepth {
					return filepath.SkipDir
				}
			}

			return nil
		}

		// Skip hidden files unless configured otherwise
		if !s.options.IncludeHidden && strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		// Check if file has a supported extension
		ext := strings.ToLower(filepath.Ext(filePath))
		if !s.extensions[ext] {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't stat
		}

		// Handle symlinks if configured
		if d.Type()&fs.ModeSymlink != 0 {
			if !s.options.FollowSymlinks {
				return nil
			}
			// Resolve symlink
			resolved, err := filepath.EvalSymlinks(filePath)
			if err != nil {
				return nil
			}
			resolvedInfo, err := os.Stat(resolved)
			if err != nil {
				return nil
			}
			info = resolvedInfo
		}

		// Parse filename using the configured parser
		fileName := d.Name()
		title, metadata := s.parser.Parse(fileName)

		result := ScanResult{
			FilePath:    filePath,
			FileName:    fileName,
			FileSize:    info.Size(),
			ParsedTitle: title,
			Metadata:    metadata,
			IsMedia:     true,
		}

		results = append(results, result)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetPaths returns the configured scan paths
func (s *FilesystemScanner) GetPaths() []string {
	return s.paths
}

// GetParser returns the configured file parser
func (s *FilesystemScanner) GetParser() FileParser {
	return s.parser
}

// GetOptions returns the scan options
func (s *FilesystemScanner) GetOptions() ScanOptions {
	return s.options
}

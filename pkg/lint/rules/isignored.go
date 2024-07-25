package rules

import (
    "path/filepath"
    "strings"
<<<<<<< HEAD
    "fmt"
)

type LintResult struct {
    Messages []string
}

func IsIgnored(errorMessage string, patterns map[string][]string) bool {
    for path, pathPatterns := range patterns {
        cleanedPath := filepath.Clean(path)
        if strings.Contains(errorMessage, cleanedPath) {
            for _, pattern := range pathPatterns {
                if strings.Contains(errorMessage, pattern) {
                    fmt.Printf("Ignoring error related to path: %s with pattern: %s\n", path, pattern)
                    return true
                }
=======
)

func IsIgnored(path string, patterns []string) bool {
    for _, pattern := range patterns {
        cleanedPath := filepath.Clean(path)
        cleanedPattern := filepath.Clean(pattern)
        if match, err := filepath.Match(cleanedPattern, cleanedPath); err == nil && match {
            return true
        }
        if strings.HasSuffix(cleanedPattern, "/") || strings.HasSuffix(cleanedPattern, "\\") {
            patternDir := strings.TrimRight(cleanedPattern, "/\\")
            if strings.HasPrefix(cleanedPath, patternDir) {
                return true
>>>>>>> ac283a55 (add the isignored file and fixed the ignore rules a bit)
            }
        }
    }
    return false
}

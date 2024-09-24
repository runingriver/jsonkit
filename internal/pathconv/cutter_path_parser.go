package pathconv

import "strings"

func PathToList(path string) (preEnd []string, end string) {
	paths := strings.Split(path, ".")
	return paths[:len(paths)-1], paths[len(paths)-1]
}

func ValidCutPath(paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, path := range paths {
		if path == "" {
			return false
		}
		if strings.HasPrefix(path, "*.") {
			return false
		}

	}
	return true
}

// ValidIncludeCutPath 不能以*号结尾
func ValidIncludeCutPath(paths []string) bool {
	for _, path := range paths {
		if strings.HasSuffix(path, ".*") || strings.HasPrefix(path, "*.") || strings.Contains(path, ".*.") {
			return false
		}
	}
	return true
}

func ValidExcludeCutPath(paths []string) bool {
	for _, path := range paths {
		if strings.HasSuffix(path, ".*") || strings.HasPrefix(path, "*.") || strings.Contains(path, ".*.") {
			return false
		}
	}
	return true
}
